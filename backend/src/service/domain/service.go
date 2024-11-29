package domain

import (
	"backend/src/storage"
	"backend/src/service"
	"backend/src/service/domain/user"
)

func New(repositories *storage.Repository) *service.Service {
	return &service.Service{
		UserService: user_service.NewUserService(repositories.UserRepository),
	}
}
