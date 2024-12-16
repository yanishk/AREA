package workflow_service

import (
	"backend/src/entities"
	"backend/src/storage"
)

type WorkflowService struct {
	WorkflowRepository storage.WorkflowRepository
	UserRepository     storage.UserRepository
}

func NewWorkflowService(WorkflowRepository storage.WorkflowRepository, UserRepository storage.UserRepository) *WorkflowService {
	return &WorkflowService{
		WorkflowRepository: WorkflowRepository,
		UserRepository:     UserRepository,
	}
}

func (self *WorkflowService) CreateWorkflow(userEmail string, newWorkflow entities.NewWorkflow) error {
	userFound, errFindingUser := self.UserRepository.FindUserByEmail(userEmail)
	if errFindingUser != nil {
		return errFindingUser
	}
	errCreationWorkflow := self.WorkflowRepository.CreateWorkflow(newWorkflow.NameWorkflow,
		userFound.Id, newWorkflow.ActionId, newWorkflow.ReactionId,
		newWorkflow.ActionParam, newWorkflow.ReactionParam)
	if errCreationWorkflow != nil {
		return errCreationWorkflow
	}
	return nil
}

func (self *WorkflowService) FindWorkflowByUserId(userId string) ([]entities.Workflow, error) {
	return self.WorkflowRepository.FindWorkflowByUserId(userId)
}

func (self *WorkflowService) ActivateDeactivateWorkflow(workflowId string, isActivated bool) error {
	return self.WorkflowRepository.ActivateDeactivateWorkflow(workflowId, isActivated)
}

func (self *WorkflowService) FindWorkflowById(workflowId string) (entities.Workflow, error) {
	return self.WorkflowRepository.FindWorkflowById(workflowId)
}

func (self *WorkflowService) GetUserWorkflows(email string) ([]entities.Workflow, error) {
	userFound, errUserFound := self.UserRepository.FindUserByEmail(email)
	if errUserFound != nil {
		return nil, errUserFound
	}

	retrievedWorkflow, err := self.WorkflowRepository.FindWorkflowByUserId(userFound.Id)
	if err != nil {
		return nil, err
	}
	return retrievedWorkflow, nil
}
