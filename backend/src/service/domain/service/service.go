package service_service

import (
	"backend/src/entities"
	"backend/src/storage"
)

type ServiceService struct {
	ServiceRepository  storage.ServiceRepository
	UserRepository     storage.UserRepository
	TriggerRepository  storage.TriggerRepository
	WorkflowRepository storage.WorkflowRepository
}

func NewServiceService(ServiceRepository storage.ServiceRepository, UserRepository storage.UserRepository,
	TriggerRepository storage.TriggerRepository, WorkflowRepository storage.WorkflowRepository) *ServiceService {
	return &ServiceService{
		ServiceRepository:  ServiceRepository,
		UserRepository:     UserRepository,
		TriggerRepository:  TriggerRepository,
		WorkflowRepository: WorkflowRepository,
	}
}

func (self *ServiceService) FindServiceByName(name string) (entities.Service, error) {
	return self.ServiceRepository.FindServiceByName(name)
}

func (self *ServiceService) FindTriggerByIdService(idService string) ([]entities.Trigger, error) {
	return self.TriggerRepository.FindTriggerByIdService(idService)
}

func (self *ServiceService) FindWorkflowByIdTrigger(idTrigger string) ([]entities.Workflow, error) {
	return self.WorkflowRepository.FindWorkflowByIdTrigger(idTrigger)
}
