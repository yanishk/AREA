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
	RetrieveActionsServices() ([]entities.ServiceUsefulInformations, error)
	RetrieveReactionsServices() ([]entities.ServiceUsefulInformations, error)
}

type UserServiceService interface {
	GetAccessTokenFromCode(code string, serviceName string) (entities.ResultToken, error)
	RefreshToken(refreshToken string, userUuid string, serviceName string) (entities.ResultToken, error)
	GetUser(emailStr string) (entities.User, error)
	FindServiceByName(name string) (entities.Service, error)
	FindUserServiceByIdServiceandIdUser(idUser string, idService string) (entities.UserService, error)
	CreateOrUpdateUserService(idUser string, idService string, token entities.ResultToken)
}

type Service struct {
	UserService        UserService
	ServiceService     ServiceService
	UserServiceService UserServiceService
}
