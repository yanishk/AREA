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

func (self *WorkflowService) ActivateDeactivateWorkflow(workflowActivation entities.ChangeActivationWorkflow) error {
	foundWorkflow, errFetchWorkflow := self.WorkflowRepository.FindWorkflowById(workflowActivation.WorkflowId)
	if errFetchWorkflow != nil {
		return errFetchWorkflow
	}

	errUpdate := self.WorkflowRepository.ActivateDeactivateWorkflow(workflowActivation.WorkflowId, foundWorkflow.IsActivated)
	if errUpdate != nil {
		return errUpdate
	}
	return nil
}
