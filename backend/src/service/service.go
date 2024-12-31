package service

import (
	"io"
	"net/http"

	"backend/src/entities"
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
	RetrieveActionsServices() ([]entities.Service, error)
	RetrieveReactionsServices() ([]entities.Service, error)
	FindServiceByName(name string) (entities.Service, error)
	FindServiceById(id string) (entities.Service, error)
	RetrieveActionsFromService(serviceName string) ([]entities.Action, error)
	RetrieveReactionsFromService(serviceName string) ([]entities.Reaction, error)
	GetGoogleRefreshTokenRequest(refreshToken string) (*http.Request, error)
	GetSpotifyRefreshTokenRequest(refreshToken string) (*http.Request, error)
	GetDiscordRefreshTokenRequest(refreshToken string) (*http.Request, error)
	GetRedditRefreshTokenRequest(refreshToken string) (*http.Request, error)
	GetAsanaRefreshTokenRequest(refreshToken string) (*http.Request, error)
	GetResultTokenFromCode(code string, serviceName string, callbackType string, appType string) (entities.ResultToken, error)
	GetUserInfoFromService(accessToken string, serviceName string) (entities.UserInfo, error)
	ExecuteRequest(request *http.Request) (*http.Response, error)
	ExecuteApiRequest(url string, method string, typeToken string, accessToken string, body io.Reader) (*http.Response, error)
	RequestToTimeApi() (entities.TimeResponse, error)
	OAuth2Service(serviceName string, callbackType string, appType string) (string, error)
}

type UserServiceService interface {
	CallApiAndRefresh(email string, connectionType string, serviceName string) (string, error)
	UpdateTokenForService(code string, serviceName string, appType string, email string, connectionType string) error
}

type WorkflowService interface {
	CreateWorkflow(userEmail string, userConnectionType string, newWorkflow entities.NewWorkflow) error
	GetUserWorkflows(email string, connectionType string) ([]entities.Workflow, error)
	UpdateWorkflow(workflowId string, workflow entities.UpdatedWorkflow) error
	DeleteWorkflow(email string, connectionType string, workflowId string) error
	CallTimeApi() error
	CallRedditApi() error
	CallWeatherApi() error
}

type AboutService interface {
	GetAboutServer(about entities.About) (entities.About, error)
}

type Service struct {
	UserService        UserService
	ServiceService     ServiceService
	UserServiceService UserServiceService
	WorkflowService    WorkflowService
	AboutService       AboutService
}
