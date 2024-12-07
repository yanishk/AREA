package service

import (
	"backend/src/entities"

	"github.com/gin-gonic/gin"
)

type UserService interface {
	CreateUser(newEmail string, newPassword string) error
	LoginAuthentication(userEmail string, userPassword string) (string, error)
	ModifyPassword(userEmail string, newPassword entities.UserModifyPassword) error
	DeleteAccount(userEmail string) error
	FindUserById(idUser string) (entities.User, error)
}

type ServiceService interface {
	FindServiceByName(name string) (entities.Service, error)
	FindTriggerByIdService(idService string) ([]entities.Trigger, error)
	FindWorkflowByIdTrigger(idTrigger string) ([]entities.Workflow, error)
	FindReactionByID(id string) (entities.Action, error)
}

type UserServiceService interface {
	GetAccesTokenFromCode(codeSpotify string) (entities.ResultToken, error)
	RefreshToken(refreshToken string, userUuid string) (entities.ResultToken, error)
	GetUser(emailStr string) (entities.User, error)
	FindServiceByName(name string) (entities.Service, error)
	FindUserServiceByIdServiceandIdUser(idservice string, userId string) (entities.UserService, error)
	CreateOrUpdateUserService(serviceId string, userId string, tokenSpotify entities.ResultToken)
	GetEmailFromContext(context *gin.Context) (string, bool)
	SpotifyCallApi(email string) (string, error)
}

type Service struct {
	UserService        UserService
	ServiceService     ServiceService
	UserServiceService UserServiceService
}
