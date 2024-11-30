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
	GettUser(emailStr string) (entities.User, error)
	CreateOrUpdateSpotify(userId string, tokenSpotify entities.ResultToken)
	FindUserServiceByName(name string, userId string) (entities.Service, error)
	RefreshToken(refreshToken string, userUuid string) (entities.ResultToken, error)
}

type Service struct {
	UserService    UserService
	ServiceService ServiceService
}
