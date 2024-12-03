package service_service

import (
	"backend/src/storage"
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
