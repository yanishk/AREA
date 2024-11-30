package domain

import (
	"backend/src/service"
	service_service "backend/src/service/domain/service"
	user_service "backend/src/service/domain/user"
	"backend/src/storage"
)

func New(repositories *storage.Repository) *service.Service {
	return &service.Service{
		UserService:    user_service.NewUserService(repositories.UserRepository),
		ServiceService: service_service.NewServiceService(repositories.ServiceRepository, repositories.UserRepository),
	}
}
