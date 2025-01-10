package workflow_service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"backend/src/entities"
)

func githubWebhooksEventsToActions() map[string]string {
	return map[string]string{
		"watch":        "New star",
		"public":       "Visibility update",
		"milestone":    "Milestone update",
		"release":      "Release update",
		"gollum":       "Wiki update",
		"workflow_job": "Workflow job update",
		"workflow_run": "Workflow run update",
		"fork":         "Fork update",
	}
}

func githubWebhookEvents() []string {
	return []string{
		"watch",
		"public",
		"milestone",
		"release",
		"gollum",
		"workflow_job",
		"workflow_run",
		"fork",
	}
}

const githubBaseUrl = "https://api.github.com/"
const githubRepositoryEndpoint = "repos/"

func (self *WorkflowService) executeGithubRequest(workflow entities.Workflow, method, accessToken, endpoint string) (*http.Response, error) {
	repository, repositoryExists := workflow.ActionParam["repository"]
	if !repositoryExists {
		return nil, fmt.Errorf(errorMissingField)
	}

	url := githubBaseUrl + githubRepositoryEndpoint + repository.(string) + endpoint
	return self.ServiceService.ExecuteApiRequest(url, method, bearerType, accessToken, nil)
}

func (self *WorkflowService) checkActionDataLen(workflow entities.Workflow, lenActualTurn float64) error {
	lenLastTurn, lenLastTurnExists := workflow.ActionData["len"]
	if !lenLastTurnExists {
		workflow.ActionData = make(map[string]interface{})
		workflow.ActionData["len"] = lenActualTurn
		err := self.WorkflowRepository.UpdateWorkflow(workflow.Id, workflow)
		if err != nil {
			return err
		}
		return nil
	}

	if lenLastTurn.(float64) < lenActualTurn {
		workflow.ActionData["len"] = lenActualTurn
		err := self.WorkflowRepository.UpdateWorkflow(workflow.Id, workflow)
		if err != nil {
			return err
		}
		self.checkReactions(workflow)
	}
	return nil
}

func (self *WorkflowService) checkGithubNewRepositoryAction(workflow entities.Workflow, accessToken string) error {
	repositories, err := self.ServiceService.RequestGithubUserRepositories(accessToken)
	if err != nil {
		return err
	}
	return self.checkActionDataLen(workflow, float64(len(repositories)))
}

func (self *WorkflowService) checkGithubNewIssueAction(workflow entities.Workflow, accessToken string) error {
	var issues []entities.GithubIssue

	resp, err := self.ServiceService.ExecuteApiRequest(githubBaseUrl+"issues", "GET", bearerType, accessToken, nil)
	if err != nil {
		return err
	}

	err = json.NewDecoder(resp.Body).Decode(&issues)
	if err != nil {
		return err
	}

	return self.checkActionDataLen(workflow, float64(len(issues)))
}

func (self *WorkflowService) checkGithubNewPullRequestAction(workflow entities.Workflow, accessToken string) error {
	var pullRequests []entities.GithubRepository

	resp, err := self.executeGithubRequest(workflow, "GET", accessToken, "/pulls")
	if err != nil {
		return err
	}

	err = json.NewDecoder(resp.Body).Decode(&pullRequests)
	if err != nil {
		return err
	}

	return self.checkActionDataLen(workflow, float64(len(pullRequests)))
}

func (self *WorkflowService) checkGithubNewBranchAction(workflow entities.Workflow, accessToken string) error {
	var branches []entities.GithubBranch

	resp, err := self.executeGithubRequest(workflow, "GET", accessToken, "/branches")
	if err != nil {
		return err
	}

	err = json.NewDecoder(resp.Body).Decode(&branches)
	if err != nil {
		return err
	}

	return self.checkActionDataLen(workflow, float64(len(branches)))
}

func (self *WorkflowService) checkGithubNewCommitAction(workflow entities.Workflow, accessToken string) error {
	var commits []entities.GithubCommit

	resp, err := self.executeGithubRequest(workflow, "GET", accessToken, "/commits")
	if err != nil {
		return err
	}

	err = json.NewDecoder(resp.Body).Decode(&commits)
	if err != nil {
		return err
	}

	return self.checkActionDataLen(workflow, float64(len(commits)))
}

func (self *WorkflowService) checkWorkflowsWithGithubActions(action entities.Action) error {
	workflows, err := self.WorkflowRepository.FindWorkflowsByActionId(action.Id)
	if err != nil {
		return err
	}

	for _, workflow := range workflows {
		if !workflow.IsActivated {
			continue
		}
		accessToken, err := self.getAccessToken("Github", workflow)
		if err != nil {
			continue
		}
		switch action.Name {
		case "New repository":
			self.checkGithubNewRepositoryAction(workflow, accessToken)
		case "New issue assignated":
			self.checkGithubNewIssueAction(workflow, accessToken)
		case "New pull request":
			self.checkGithubNewPullRequestAction(workflow, accessToken)
		case "New branch":
			self.checkGithubNewBranchAction(workflow, accessToken)
		case "New push":
			self.checkGithubNewCommitAction(workflow, accessToken)
		}
	}
	return nil
}

func (self *WorkflowService) CheckGithubActions() error {
	service, err := self.ServiceService.FindServiceByName("Github")
	if err != nil {
		return err
	}

	actions, err := self.ActionRepository.FindActionsByServiceId(service.Id)
	if err != nil {
		return err
	}

	for _, action := range actions {
		err := self.checkWorkflowsWithGithubActions(action)
		if err != nil {
			return err
		}
	}
	return nil
}

func (self *WorkflowService) checkGithubWebhooksWorkflow(workflow entities.Workflow, webhookRepository string) error {
	userRepository, err := getWorkflowStringActionParam(workflow, "repository")
	if err != nil {
		return err
	}

	if userRepository == webhookRepository {
		self.checkReactions(workflow)
	}
	return nil
}

func (self *WorkflowService) checkGithubWebhooksWorkflows(headers http.Header, webhookJsonDataBytes []byte, eventsToActions map[string]string, serviceId string) error {
	var webhookResponse entities.GithubWebhookTriggeredResponse
	if len(webhookJsonDataBytes) > 0 {
		err := json.Unmarshal(webhookJsonDataBytes, &webhookResponse)
		if err != nil {
			return err
		}
	}
	if webhookResponse.Repository.Name == "" {
		return fmt.Errorf("Incorrect repository name")
	}

	eventName := headers["X-Github-Event"][0]
	if eventName == "ping" {
		return nil
	}

	actionName, actionNameExists := eventsToActions[eventName]
	if !actionNameExists {
		return fmt.Errorf("No action mapped with this event")
	}

	action, err := self.ActionRepository.FindActionByNameAndServiceId(actionName, serviceId)
	if err != nil {
		return err
	}

	workflows, err := self.WorkflowRepository.FindWorkflowsByActionId(action.Id)
	if err != nil {
		return err
	}

	for _, workflow := range workflows {
		if !workflow.IsActivated {
			continue
		}
		self.checkGithubWebhooksWorkflow(workflow, webhookResponse.Repository.Name)
	}
	return nil
}

func (self *WorkflowService) createNewWorkflowGithubWebhook(repository, accessToken string) error {
	var webhook entities.GithubWebhookResponse
	createWebhookUrl := githubBaseUrl + "repos/" + repository + "/hooks"
	webhookUrl := os.Getenv("NGROK_APP_URL") + "/webhooks/Github"

	webhook.Name = "web"
	webhook.Active = true
	webhook.Events = githubWebhookEvents()
	webhook.Config.Url = webhookUrl
	webhook.Config.ContentType = "json"
	jsonBody, err := json.Marshal(webhook)
	if err != nil {
		return fmt.Errorf(errorMarshaling)
	}

	res, err := self.ServiceService.ExecuteApiRequest(createWebhookUrl, "POST", bearerType, accessToken, bytes.NewBuffer([]byte(jsonBody)))
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return nil
}

func (self *WorkflowService) isGithubRepositoryWebhookPresent(repository, accessToken string) (bool, error) {
	var webhooks []entities.GithubWebhookResponse
	getWebhooksUrl := githubBaseUrl + githubRepositoryEndpoint + repository + "/hooks"
	webhookUrl := os.Getenv("NGROK_APP_URL") + "/webhooks/Github"

	res, err := self.ServiceService.ExecuteApiRequest(getWebhooksUrl, "GET", bearerType, accessToken, nil)
	if err != nil {
		return false, err
	}
	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(&webhooks)
	if err != nil {
		return false, err
	}

	for _, webhook := range webhooks {
		if webhook.Config.Url == webhookUrl {
			return true, nil
		}
	}
	return false, nil
}

func (self *WorkflowService) checkNewWorkflowGithubWebhook(workflow entities.Workflow, accessToken string) error {
	repository, err := getWorkflowStringActionParam(workflow, "repository")
	if err != nil {
		return err
	}

	isWebhookPresent, err := self.isGithubRepositoryWebhookPresent(repository, accessToken)
	if err != nil {
		return err
	}

	if !isWebhookPresent {
		self.createNewWorkflowGithubWebhook(repository, accessToken)
		err := self.WorkflowRepository.UpdateWorkflow(workflow.Id, workflow)
		if err != nil {
			return err
		}
	}
	return nil
}

func (self *WorkflowService) checkNewWorkflowsGithubWebhook(action entities.Action) error {
	workflows, err := self.WorkflowRepository.FindWorkflowsByActionId(action.Id)
	if err != nil {
		return err
	}

	for _, workflow := range workflows {
		accessToken, err := self.getAccessToken("Github", workflow)
		if err != nil {
			continue
		}
		self.checkNewWorkflowGithubWebhook(workflow, accessToken)
	}
	return nil
}

func (self *WorkflowService) CheckNewGithubWorkflows() error {
	service, err := self.ServiceService.FindServiceByName("Github")
	if err != nil {
		return err
	}

	actions, err := self.ActionRepository.FindActionsByServiceId(service.Id)
	if err != nil {
		return err
	}

	for _, action := range actions {
		err := self.checkNewWorkflowsGithubWebhook(action)
		if err != nil {
			return err
		}
	}
	return nil
}
