package service_service

import (
	"backend/src/storage"
	"backend/src/entities"
)

type ServiceService struct {
	ServiceRepository storage.ServiceRepository
	UserRepository    storage.UserRepository
}

func NewServiceService(ServiceRepository storage.ServiceRepository, UserRepository storage.UserRepository) *ServiceService {
	return &ServiceService{
		ServiceRepository: ServiceRepository,
		UserRepository:    UserRepository,
	}
}

func (self *ServiceService) RetrieveActionsServices() ([]entities.ServiceUsefulInformations, error) {
	return self.ServiceRepository.FindActionsServices()
}

func (self *ServiceService) RetrieveReactionsServices() ([]entities.ServiceUsefulInformations, error) {
	return self.ServiceRepository.FindReactionsServices()
}
