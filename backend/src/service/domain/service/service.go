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
	ReactionRepository storage.ReactionRepository
}

func NewServiceService(ServiceRepository storage.ServiceRepository, UserRepository storage.UserRepository,
	TriggerRepository storage.TriggerRepository, WorkflowRepository storage.WorkflowRepository,
	ReactionRepository storage.ReactionRepository) *ServiceService {
	return &ServiceService{
		ServiceRepository:  ServiceRepository,
		UserRepository:     UserRepository,
		TriggerRepository:  TriggerRepository,
		WorkflowRepository: WorkflowRepository,
		ReactionRepository: ReactionRepository,
	}
}

func (self *ServiceService) FindServiceByName(name string) (entities.Service, error) {
	return self.ServiceRepository.FindServiceByName(name)
}

func (self *ServiceService) FindTriggerByServiceId(serviceId string) ([]entities.Trigger, error) {
	return self.TriggerRepository.FindTriggerByServiceId(serviceId)
}

func (self *ServiceService) FindWorkflowByTriggerId(triggerId string) ([]entities.Workflow, error) {
	return self.WorkflowRepository.FindWorkflowByTriggerId(triggerId)
}

func (self *ServiceService) FindReactionByID(id string) (entities.Action, error) {
	return self.ReactionRepository.FindReactionByID(id)
}

func (self *ServiceService) RetrieveActionsServices() ([]entities.ServiceUsefulInformations, error) {
	return self.ServiceRepository.FindActionsServices()
}

func (self *ServiceService) RetrieveReactionsServices() ([]entities.ServiceUsefulInformations, error) {
	return self.ServiceRepository.FindReactionsServices()
}
