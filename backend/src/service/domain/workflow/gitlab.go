package workflow_service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"backend/src/entities"
)

func gitlabWebhooksEventsToActions() map[string]string {
	return map[string]string{
		"Push Hook":          "New push",
		"Merge Request Hook": "Merge request update",
		"Issue Hook":         "Issue update",
		"Note Hook":          "Comment update",
		"Tag Push Hook":      "New tag push",
		"Wiki Page Hook":     "Wiki page update",
		"Release Hook":       "Release update",
		"Feature Flag Hook":  "Feature flag update",
		"Pipeline Hook":      "Pipeline update",
		"Job Hook":           "Job update",
		"Deployment Hook":    "Deployment update",
		"Emoji Hook":         "Emoji update",
	}
}

func gitlabWebhookEventsToState() map[string]string {
	return map[string]string{
		"push_events":           "true",
		"issues_events":         "true",
		"merge_requests_events": "true",
		"tag_push_events":       "true",
		"note_events":           "true",
		"job_events":            "true",
		"pipeline_events":       "true",
		"wiki_page_events":      "true",
		"deployment_events":     "true",
		"releases_events":       "true",
		"feature_flag_events":   "true",
		"emoji_events":          "true",
	}
}

func (self *WorkflowService) checkGitlabWebhooksWorkflow(workflow entities.Workflow, webhookProjectId int) error {
	webhookProjectIdString := strconv.Itoa(webhookProjectId)
	userProjectIdString, err := getWorkflowStringActionParam(workflow, "project")
	if err != nil {
		return err
	}

	if userProjectIdString == webhookProjectIdString {
		self.checkReactions(workflow)
	}
	return nil
}

func (self *WorkflowService) checkGitlabWebhooksWorkflows(headers http.Header, webhookJsonDataBytes []byte, eventsToActions map[string]string, serviceId string) error {
	var webhookResponse entities.GitlabWebhookTriggeredResponse
	if len(webhookJsonDataBytes) > 0 {
		err := json.Unmarshal(webhookJsonDataBytes, &webhookResponse)
		if err != nil {
			return err
		}
	}
	if webhookResponse.Project.Id == 0 {
		return fmt.Errorf("Incorrect project id")
	}

	eventName := headers["X-Gitlab-Event"][0]
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
		self.checkGitlabWebhooksWorkflow(workflow, webhookResponse.Project.Id)
	}
	return nil
}

func (self *WorkflowService) executeGitlabRequest(method, url, accessToken string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", bearerType+accessToken)
	req.Header.Set(contentType, "application/x-www-form-urlencoded")

	res, err := self.ServiceService.ExecuteRequest(req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (self *WorkflowService) createNewWorkflowGitlabWebhook(projectId, accessToken string) error {
	createWebhookUrl := "https://gitlab.com/api/v4/projects/" + projectId + "/hooks"
	webhookUrl := os.Getenv("NGROK_APP_URL") + "/webhooks/Gitlab"
	params := url.Values{}

	params.Set("url", webhookUrl)
	for webhookEvent, webhookState := range gitlabWebhookEventsToState() {
		params.Set(webhookEvent, webhookState)
	}

	res, err := self.executeGitlabRequest("POST", createWebhookUrl, accessToken, bytes.NewBufferString(params.Encode()))
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return nil
}

func (self *WorkflowService) isGitlabProjectWebhookPresent(projectId, accessToken string) (bool, error) {
	webhookUrl := os.Getenv("NGROK_APP_URL") + "/webhooks/Gitlab"
	getWebhooksUrl := "https://gitlab.com/api/v4/projects/" + projectId + "/hooks"

	res, err := self.executeGitlabRequest("GET", getWebhooksUrl, accessToken, nil)
	if err != nil {
		return false, err
	}
	defer res.Body.Close()

	webhooksJsonDataBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return false, err
	}

	webhooksJsonData, err := unmarshalJsonToMap(webhooksJsonDataBytes)
	if err != nil {
		return false, err
	}

	if len(webhooksJsonData) == 0 {
		return false, nil
	}

	for _, webhookJsonData := range webhooksJsonData {
		url, urlExists := webhookJsonData["url"]
		if !urlExists {
			return false, fmt.Errorf(errorMissingField)
		}

		urlString, urlIsString := url.(string)
		if !urlIsString {
			return false, fmt.Errorf(errorMissingField)
		}

		if urlString == webhookUrl {
			return true, nil
		}
	}

	return false, nil
}

func (self *WorkflowService) checkNewWorkflowGitlabWebhook(workflow entities.Workflow, accessToken string) error {
	projectId, err := getWorkflowStringActionParam(workflow, "project")
	if err != nil {
		return err
	}

	isWebhookPresent, err := self.isGitlabProjectWebhookPresent(projectId, accessToken)
	if err != nil {
		return err
	}

	if !isWebhookPresent {
		self.createNewWorkflowGitlabWebhook(projectId, accessToken)
		err := self.WorkflowRepository.UpdateWorkflow(workflow.Id, workflow)
		if err != nil {
			return err
		}
	}
	return nil
}

func (self *WorkflowService) checkNewWorkflowsGitlabWebhook(action entities.Action) error {
	workflows, err := self.WorkflowRepository.FindWorkflowsByActionId(action.Id)
	if err != nil {
		return err
	}

	for _, workflow := range workflows {
		accessToken, err := self.getAccessToken("Gitlab", workflow)
		if err != nil {
			continue
		}
		self.checkNewWorkflowGitlabWebhook(workflow, accessToken)
	}
	return nil
}

func (self *WorkflowService) CheckNewGitlabWorkflows() error {
	service, err := self.ServiceService.FindServiceByName("Gitlab")
	if err != nil {
		return err
	}

	actions, err := self.ActionRepository.FindActionsByServiceId(service.Id)
	if err != nil {
		return err
	}

	for _, action := range actions {
		err := self.checkNewWorkflowsGitlabWebhook(action)
		if err != nil {
			return err
		}
	}
	return nil
}
