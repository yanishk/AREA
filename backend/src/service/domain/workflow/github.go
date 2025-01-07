package workflow_service

import (
	"encoding/json"
	"fmt"
	"net/http"

	"backend/src/entities"
)

const githubBaseUrl = "https://api.github.com/"

func (self *WorkflowService) executeGithubRequest(workflow entities.Workflow, accessToken, endpoint string) (*http.Response, error) {
	repository, repositoryExists := workflow.ActionParam["repository"]
	if !repositoryExists {
		return nil, fmt.Errorf(errorMissingField)
	}

	url := githubBaseUrl + "repos/" + repository.(string) + endpoint
	return self.ServiceService.ExecuteApiRequest(url, "GET", bearerType, accessToken, nil)
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

	resp, err := self.executeGithubRequest(workflow, accessToken, "/pulls")
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

	resp, err := self.executeGithubRequest(workflow, accessToken, "/branches")
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

	resp, err := self.executeGithubRequest(workflow, accessToken, "/commits")
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
