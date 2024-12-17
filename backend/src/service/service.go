package service

import (
	"backend/src/entities"

	"github.com/gin-gonic/gin"
)

type UserService interface {
	CreateUser(userEmail string, userPassword string, userConnectionType string) error
	LoginAuthentication(userEmail string, userPassword string, userConnectionType string) (string, error)
	ModifyPassword(userEmail string, newPassword entities.UserModifyPassword) error
	DeleteAccount(userEmail string, userConnectionType string) error
	FindUserById(userId string) (entities.User, error)
}

type ServiceService interface {
	RetrieveActionsServices() ([]entities.ServiceUsefulInformations, error)
	RetrieveReactionsServices() ([]entities.ServiceUsefulInformations, error)
	FindServiceByName(name string) (entities.Service, error)
	FindWorkflowByActionId(actionId string) ([]entities.Workflow, error)
	FindReactionByID(id string) (entities.Reaction, error)
	FindActionByServiceId(serviceName string) ([]entities.Action, error)
	RetrieveReactionsFromService(serviceName string) ([]entities.Reaction, error)
}

type UserServiceService interface {
	GetEmailFromContext(context *gin.Context) (string, bool)
	CallApiAndRefresh(email string, connectionType string, serviceName string) (string, error)
	UpdateTokenForService(code string, serviceName string, email string, connectionType string) error
}

type WorkflowService interface {
	CreateWorkflow(userEmail string, userConnectionType string, newWorkflow entities.NewWorkflow) error
	ActivateDeactivateWorkflow(workflowActivation entities.ChangeActivationWorkflow) error
	GetUserWorkflows(email string, connectionType string) ([]entities.Workflow, error)
}

type Service struct {
	UserService        UserService
	ServiceService     ServiceService
	UserServiceService UserServiceService
	WorkflowService    WorkflowService
}
