package service

import (
	"io"
	"net/http"

	"backend/src/entities"
)

type UserService interface {
	CreateUser(userEmail, userPassword, userConnectionType string) error
	LoginAuthentication(userEmail, userPassword, userConnectionType string) (string, error)
	LoginWithService(code, serviceName, appType string) (string, error)
	GetUser(userEmail, userConnectionType string) (entities.UserInfos, error)
	ModifyPassword(userEmail, userConnectionType string, newPassword entities.UserModifyPassword) error
	DeleteAccount(userEmail, userConnectionType string) error
	FindUserById(userId string) (entities.User, error)
}

type ServiceService interface {
	RetrieveActionsServices() ([]entities.Service, error)
	RetrieveReactionsServices() ([]entities.Service, error)
	FindServiceByName(name string) (entities.Service, error)
	FindServiceById(id string) (entities.Service, error)
	FindServiceByActionId(id string) (entities.Service, error)
	FindServiceByReactionId(id string) (entities.Service, error)
	FindAllServices() ([]entities.Service, error)
	RetrieveActionsFromService(serviceName string) ([]entities.Action, error)
	RetrieveReactionsFromService(serviceName string) ([]entities.Reaction, error)
	GetGoogleRefreshTokenRequest(refreshToken string) (*http.Request, error)
	GetSpotifyRefreshTokenRequest(refreshToken string) (*http.Request, error)
	GetDiscordRefreshTokenRequest(refreshToken string) (*http.Request, error)
	GetRedditRefreshTokenRequest(refreshToken string) (*http.Request, error)
	GetAsanaRefreshTokenRequest(refreshToken string) (*http.Request, error)
	GetDropboxRefreshTokenRequest(refreshToken string) (*http.Request, error)
	GetGitlabRefreshTokenRequest(refreshToken string) (*http.Request, error)
	ExecuteRequest(request *http.Request) (*http.Response, error)
	ExecuteApiRequest(url, method, typeToken, accessToken string, body io.Reader) (*http.Response, error)
	GetResultTokenFromCode(code, serviceName, callbackType, appType string) (entities.ResultToken, error)
	GetUserInfoFromService(accessToken, serviceName string) (entities.UserInfo, error)
	RequestToTimeApi() (entities.TimeResponse, error)
	OAuth2Service(serviceName, callbackType, appType string) (string, error)
	RequestGithubUserRepositories(accessToken string) ([]entities.GithubRepository, error)
	RequestGitlabUserProjects(accessToken string) ([]entities.GitlabProject, error)
	RetrieveDiscordGuildChannels(guildId string) ([]map[string]interface{}, error)
}

type UserServiceService interface {
	RetrieveUserServiceAuthenticationStatus(email, connectionType, serviceName string) (bool, error)
	CallApiAndRefresh(email, connectionType, serviceName string) (string, error)
	UpdateTokenForService(code, serviceName, appType, email, connectionType string) error
	RetrieveGithubUserRepositories(email, connectionType string) ([]entities.GithubRepository, error)
	RetrieveGitlabUserProjects(email, connectionType string) ([]entities.GitlabProject, error)
	RetrieveDiscordUserServers(email, connectionType string) ([]map[string]interface{}, error)
	RetrieveAsanaUserWorkspaces(email, connectionType string) (entities.AsanaWorkspacesInfo, error)
	RetrieveAsanaWorkspaceAssignees(email, connectionType, workspaceId string) (entities.AsanaWorkspacesInfo, error)
	RetrieveAsanaWorkspaceProjects(email, connectionType, workspaceId string) (entities.AsanaWorkspacesInfo, error)
	RetrieveAsanaWorkspaceTags(email, connectionType, workspaceId string) (entities.AsanaWorkspacesInfo, error)
}

type WorkflowService interface {
	CreateWorkflow(userEmail, userConnectionType string, newWorkflow entities.NewWorkflow) error
	GetUserWorkflows(email, connectionType string) ([]entities.Workflow, error)
	UpdateWorkflow(workflowId string, workflow entities.UpdatedWorkflow) error
	DeleteWorkflow(email, connectionType, workflowId string) error
	CheckTimeAndDateActions() error
	CheckGithubActions() error
	CheckRedditActions() error
	CheckWeatherActions() error
	CheckNewGitlabWorkflows() error
	CheckNewGithubWorkflows() error
	CheckWebhooksWorkflows(serviceName string, request *http.Request) error
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
