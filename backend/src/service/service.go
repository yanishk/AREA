package service

import (
	"backend/src/entities"
)

type UserService interface {
	CreateUser(newEmail string, newPassword string) error
	LoginAuthentication(userEmail string, userPassword string) (string, error)
	ModifyPassword(userEmail string, newPassword entities.UserModifyPassword) error
	DeleteAccount(userEmail string) error
}

type ServiceService interface {
	GetAccesTokenFromCode(codeSpotify string) (entities.ResultToken, error)
	GetUser(emailStr string) (entities.User, error)
	CreateOrUpdateUserService(serviceId string, userId string, tokenSpotify entities.ResultToken)
	RefreshToken(refreshToken string, userUuid string) (entities.ResultToken, error)
	FindServiceByName(name string) (entities.Service, error)
	FindUserServiceByIdServiceandIdUser(idservice string, userId string) (entities.UserService, error)
}

type Service struct {
	UserService    UserService
	ServiceService ServiceService
}
