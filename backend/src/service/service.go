package service

import (
	"backend/src/entities"
)

type UserService interface {
	CreateUser(newEmail string, newPassword string) (error)
	LoginAuthentication(userEmail string, userPassword string) (string, error)
	ModifyPassword(userEmail string, newPassword entities.UserModifyPassword) error
	DeleteAccount(userEmail string) error
}

type Service struct {
	UserService UserService
}
