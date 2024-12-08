package workflow_service

import "backend/src/storage"

type WorkflowService struct {
	WorkflowRepository storage.WorkflowRepository
}

func NewWorkflowService(WorkflowRepository storage.WorkflowRepository) *WorkflowService {
	return &WorkflowService{WorkflowRepository: WorkflowRepository}
}

func (self *WorkflowService) CreateWorkflow(name string, ownerId string, actionId string, reactionId string,
	actionParam map[string]interface{}, reactionParam map[string]interface{}) error {
	return self.WorkflowRepository.CreateWorkflow(name, ownerId, actionId, reactionId, actionParam, reactionParam)
}
