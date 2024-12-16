package service_service

import (
	"backend/src/entities"
	"backend/src/storage"
)

type ServiceService struct {
	ServiceRepository  storage.ServiceRepository
	UserRepository     storage.UserRepository
	ActionRepository   storage.ActionRepository
	WorkflowRepository storage.WorkflowRepository
	ReactionRepository storage.ReactionRepository
}

func NewServiceService(ServiceRepository storage.ServiceRepository, UserRepository storage.UserRepository,
	ActionRepository storage.ActionRepository, WorkflowRepository storage.WorkflowRepository,
	ReactionRepository storage.ReactionRepository) *ServiceService {
	return &ServiceService{
		ServiceRepository:  ServiceRepository,
		UserRepository:     UserRepository,
		ActionRepository:   ActionRepository,
		WorkflowRepository: WorkflowRepository,
		ReactionRepository: ReactionRepository,
	}
}

func (self *ServiceService) FindServiceByName(name string) (entities.Service, error) {
	return self.ServiceRepository.FindServiceByName(name)
}

func (self *ServiceService) FindActionByServiceId(serviceName string) ([]entities.Action, error) {
	foundService, errFoundService := self.ServiceRepository.FindServiceByName(serviceName)
	if errFoundService != nil {
		return nil, errFoundService
	}

	actions, errFindActions := self.ActionRepository.FindActionByServiceId(foundService.Id)
	if errFindActions != nil {
		return nil, errFindActions
	}
	return actions, nil
}

func (self *ServiceService) FindWorkflowByActionId(actionId string) ([]entities.Workflow, error) {
	return self.WorkflowRepository.FindWorkflowByActionId(actionId)
}

func (self *ServiceService) FindReactionByServiceId(serviceId string) ([]entities.Reaction, error) {
	return self.ReactionRepository.FindReactionByServiceId(serviceId)
}

func (self *ServiceService) FindReactionByID(id string) (entities.Reaction, error) {
	return self.ReactionRepository.FindReactionByID(id)
}

func (self *ServiceService) RetrieveActionsServices() ([]entities.ServiceUsefulInformations, error) {
	return self.ServiceRepository.FindActionsServices()
}

func (self *ServiceService) RetrieveReactionsServices() ([]entities.ServiceUsefulInformations, error) {
	return self.ServiceRepository.FindReactionsServices()
}
