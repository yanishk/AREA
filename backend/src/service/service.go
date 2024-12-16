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
	FindUserById(userId string) (entities.User, error)
	FindUserByEmail(email string) (entities.User, error)
}

type ServiceService interface {
	RetrieveActionsServices() ([]entities.ServiceUsefulInformations, error)
	RetrieveReactionsServices() ([]entities.ServiceUsefulInformations, error)
	FindServiceByName(name string) (entities.Service, error)
	FindReactionByServiceId(serviceId string) ([]entities.Reaction, error)
	FindWorkflowByActionId(actionId string) ([]entities.Workflow, error)
	FindReactionByID(id string) (entities.Reaction, error)
	FindActionByServiceId(serviceName string) ([]entities.Action, error)
}

type UserServiceService interface {
	GetAccessTokenFromCode(code string, serviceName string) (entities.ResultToken, error)
	RefreshToken(refreshToken string, userUuid string, serviceName string) (entities.ResultToken, error)
	GetUser(emailStr string) (entities.User, error)
	FindServiceByName(name string) (entities.Service, error)
	FindUserServiceByServiceIdandUserId(userId string, serviceId string) (entities.UserService, error)
	CreateOrUpdateUserService(userId string, serviceId string, token entities.ResultToken)
	GetEmailFromContext(context *gin.Context) (string, bool)
	CallApiAndRefresh(email string, serviceName string) (string, error)
}

type WorkflowService interface {
	CreateWorkflow(userEmail string, newWorkflow entities.NewWorkflow) error
	FindWorkflowByUserId(userId string) ([]entities.Workflow, error)
	ActivateDeactivateWorkflow(workflowId string, isActivated bool) error
	FindWorkflowById(workflowId string) (entities.Workflow, error)
}

type Service struct {
	UserService        UserService
	ServiceService     ServiceService
	UserServiceService UserServiceService
	WorkflowService    WorkflowService
}
