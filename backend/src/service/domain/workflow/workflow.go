package workflow_service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"backend/src/entities"
	"backend/src/service"
	"backend/src/storage"
)

type WorkflowService struct {
	WorkflowRepository storage.WorkflowRepository
	UserRepository     storage.UserRepository
	ActionRepository   storage.ActionRepository
	ReactionRepository storage.ReactionRepository
	ServiceService     service.ServiceService
	UserServiceService service.UserServiceService
}

const bearerType = "Bearer "
const contentType = "Content-Type"

const errorUpdatingWorkflow = "Could not update workflow"
const errorUpdatingToken = "Could not update token"
const errorRetrievingReaction = "Error finding reaction"
const errorMissingField = "Missing required field"
const errorMarshaling = "Could not marshal JSON"

func NewWorkflowService(WorkflowRepository storage.WorkflowRepository, UserRepository storage.UserRepository,
	ActionRepository storage.ActionRepository, ReactionRepository storage.ReactionRepository,
	ServiceService service.ServiceService, UserServiceService service.UserServiceService) *WorkflowService {
	return &WorkflowService{
		WorkflowRepository: WorkflowRepository,
		UserRepository:     UserRepository,
		ActionRepository:   ActionRepository,
		ReactionRepository: ReactionRepository,
		ServiceService:     ServiceService,
		UserServiceService: UserServiceService,
	}
}

func (self *WorkflowService) CreateWorkflow(userEmail, userConnectionType string, newWorkflow entities.NewWorkflow) error {
	userFound, errFindingUser := self.UserRepository.FindUserByEmail(userEmail, userConnectionType)
	if errFindingUser != nil {
		return errFindingUser
	}
	errCreationWorkflow := self.WorkflowRepository.CreateWorkflow(newWorkflow.Name,
		userFound.Id, newWorkflow.ActionId, newWorkflow.ReactionId,
		newWorkflow.ActionParam, newWorkflow.ReactionParam, newWorkflow.ActionData)
	if errCreationWorkflow != nil {
		return errCreationWorkflow
	}
	return nil
}

func (self *WorkflowService) GetUserWorkflows(email, connectionType string) ([]entities.Workflow, error) {
	userFound, errUserFound := self.UserRepository.FindUserByEmail(email, connectionType)
	if errUserFound != nil {
		return nil, errUserFound
	}

	retrievedWorkflow, err := self.WorkflowRepository.FindWorkflowsByOwnerId(userFound.Id)
	if err != nil {
		return nil, err
	}
	return retrievedWorkflow, nil
}

func (self *WorkflowService) UpdateWorkflow(workflowId string, workflow entities.UpdatedWorkflow) error {
	updatedWorkflow, err := self.WorkflowRepository.FindWorkflowById(workflowId)
	if err != nil {
		return err
	}

	if workflow.Name != nil {
		updatedWorkflow.Name = *workflow.Name
	}
	if workflow.ActionId != nil {
		updatedWorkflow.ActionId = *workflow.ActionId
	}
	if workflow.ReactionId != nil {
		updatedWorkflow.ReactionId = *workflow.ReactionId
	}
	if workflow.IsActivated != nil {
		updatedWorkflow.IsActivated = *workflow.IsActivated
	}
	if workflow.ActionParam != nil {
		updatedWorkflow.ActionParam = *workflow.ActionParam
	}
	if workflow.ReactionParam != nil {
		updatedWorkflow.ReactionParam = *workflow.ReactionParam
	}

	err = self.WorkflowRepository.UpdateWorkflow(workflowId, updatedWorkflow)
	if err != nil {
		return err
	}
	return nil
}

func (self *WorkflowService) DeleteWorkflow(email, connectionType, workflowId string) error {
	user, err := self.UserRepository.FindUserByEmail(email, connectionType)
	if err != nil {
		return err
	}

	err = self.WorkflowRepository.DeleteWorkflow(workflowId, user.Id)
	if err != nil {
		return err
	}
	return nil
}

func (self *WorkflowService) getAccessToken(serviceName string, workflow entities.Workflow) (string, error) {
	foundUser, errUser := self.UserRepository.FindUserById(workflow.OwnerId)
	if errUser != nil {
		return "", fmt.Errorf("Error finding user")
	}
	accessToken, errRefreshing := self.UserServiceService.CallApiAndRefresh(foundUser.Email, foundUser.ConnectionType, serviceName)
	if errRefreshing != nil {
		return "", errRefreshing
	}
	return accessToken, nil
}

func (self *WorkflowService) refreshTokenForService(serviceName, reactionFoundName string, reactionsPossible []string, workflow entities.Workflow) (string, error) {
	var accessToken string
	var err error

	for i := range reactionsPossible {
		if reactionFoundName == reactionsPossible[i] {
			accessToken, err = self.getAccessToken(serviceName, workflow)
			if err != nil {
				return accessToken, fmt.Errorf(errorUpdatingToken)
			}
		}
	}
	return accessToken, nil
}

func concatanateSlashToPath(path, name string) string {
	if path[0] != '/' {
		path = "/" + path
	}
	if path[len(path)-1] != '/' {
		path = path + "/"
	}

	return path + name
}

func unmarshalJsonToMap(jsonDataBytes []byte) ([]map[string]interface{}, error) {
	var jsonData []map[string]interface{}
	if len(jsonDataBytes) > 0 {
		err := json.Unmarshal(jsonDataBytes, &jsonData)
		if err != nil {
			return jsonData, err
		}
	}
	return jsonData, nil
}

func getWorkflowStringActionParam(workflow entities.Workflow, paramKey string) (string, error) {
	param, paramExists := workflow.ActionParam[paramKey]
	if !paramExists {
		return "", fmt.Errorf(errorMissingField)
	}

	paramString, paramIsString := param.(string)
	if !paramIsString {
		return "", fmt.Errorf(errorMissingField)
	}
	return paramString, nil
}

func (self *WorkflowService) CheckWebhooksWorkflows(serviceName string, request *http.Request) error {
	webhookJsonDataBytes, err := io.ReadAll(request.Body)
	if err != nil {
		return err
	}

	service, err := self.ServiceService.FindServiceByName(serviceName)
	if err != nil {
		return err
	}

	switch service.Name {
	case "Gitlab":
		return self.checkGitlabWebhooksWorkflows(request.Header, webhookJsonDataBytes, gitlabWebhooksEventsToActions(), service.Id)
	case "Github":
		return self.checkGithubWebhooksWorkflows(request.Header, webhookJsonDataBytes, githubWebhooksEventsToActions(), service.Id)
	}
	return nil
}

func (self *WorkflowService) checkReactions(workflow entities.Workflow) {
	self.checkSpotifyReactions(workflow)
	self.checkDiscordReactions(workflow)
	self.checkLinkedinReactions(workflow)
	self.checkAsanaReactions(workflow)
	self.checkSMSReactions(workflow)
	self.checkSendEmailReactions(workflow)
	self.checkDropboxReactions(workflow)
	self.checkRedditReactions(workflow)
}
