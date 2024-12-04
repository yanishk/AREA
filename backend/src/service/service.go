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
}

type UserServiceService interface {
	GetAccesTokenFromCode(codeSpotify string) (entities.ResultToken, error)
	RefreshToken(refreshToken string, userUuid string) (entities.ResultToken, error)
	GetUser(emailStr string) (entities.User, error)
	FindServiceByName(name string) (entities.Service, error)
	FindUserServiceByIdServiceandIdUser(idservice string, userId string) (entities.UserService, error)
	CreateOrUpdateUserService(serviceId string, userId string, tokenSpotify entities.ResultToken)
}

type Service struct {
	UserService        UserService
	ServiceService     ServiceService
	UserServiceService UserServiceService
}
