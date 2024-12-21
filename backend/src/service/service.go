package service

import (
	"backend/src/entities"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserService interface {
	CreateUser(userEmail string, userPassword string, userConnectionType string) error
	LoginAuthentication(userEmail string, userPassword string, userConnectionType string) (string, error)
	LoginWithService(code string, serviceName string, appType string) (string, error)
	ModifyPassword(userEmail string, userConnectionType string, newPassword entities.UserModifyPassword) error
	DeleteAccount(userEmail string, userConnectionType string) error
	FindUserById(userId string) (entities.User, error)
}

type ServiceService interface {
	RetrieveActionsServices() ([]entities.ServiceUsefulInformations, error)
	RetrieveReactionsServices() ([]entities.ServiceUsefulInformations, error)
	FindServiceByName(name string) (entities.Service, error)
	FindWorkflowByActionId(actionId string) ([]entities.Workflow, error)
	FindReactionByID(id string) (entities.Reaction, error)
	RetrieveActionsFromService(serviceName string) ([]entities.Action, error)
	RetrieveReactionsFromService(serviceName string) ([]entities.Reaction, error)
	FindActionByServiceId(serviceId string) ([]entities.Action, error)
	GetGoogleRefreshTokenRequest(refreshToken string) (*http.Request, error)
	GetSpotifyRefreshTokenRequest(refreshToken string) (*http.Request, error)
	GetResultTokenFromCode(code string, serviceName string, callbackType string, appType string) (entities.ResultToken, error)
	GetUserInfoFromService(accessToken string, serviceName string) (entities.UserInfo, error)
	RequestService(request *http.Request) (*http.Response, error)
}

type UserServiceService interface {
	GetEmailFromContext(context *gin.Context) (string, bool)
	GetConnectionTypeFromContext(context *gin.Context) (string, bool)
	CallApiAndRefresh(email string, connectionType string, serviceName string) (string, error)
	UpdateTokenForService(code string, serviceName string, appType string, email string, connectionType string) error
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
