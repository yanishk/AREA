package domain

import (
	"backend/src/service"
	service_service "backend/src/service/domain/service"
	user_service "backend/src/service/domain/user"
	user_service_service "backend/src/service/domain/userservice"
	"backend/src/storage"
)

func New(repositories *storage.Repository) *service.Service {
	return &service.Service{
		UserService:        user_service.NewUserService(repositories.UserRepository),
		ServiceService:     service_service.NewServiceService(repositories.ServiceRepository, repositories.UserRepository),
		UserServiceService: user_service_service.NewUserServiceService(repositories.ServiceRepository, repositories.UserRepository, repositories.UserServiceRepository),
	}
}
