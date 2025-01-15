package userservice_service

import (
	"backend/src/entities"
	"io"
	"net/http"

	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) CreateUser(email, password, connectionType string) error {
	args := m.Called(email, password, connectionType)
	return args.Error(0)
}

func (m *MockUserRepository) FindUserByEmail(email, connectionType string) (entities.User, error) {
	args := m.Called(email, connectionType)
	return args.Get(0).(entities.User), args.Error(1)
}

func (m *MockUserRepository) FindUserById(userId string) (entities.User, error) {
	args := m.Called(userId)
	return args.Get(0).(entities.User), args.Error(1)
}

func (m *MockUserRepository) UpdateUser(email, password, connectionType string) error {
	args := m.Called(email, password, connectionType)
	return args.Error(0)
}

func (m *MockUserRepository) DeleteUser(email, connectionType string) error {
	args := m.Called(email, connectionType)
	return args.Error(0)
}

type MockServiceRepository struct {
	mock.Mock
}

func (m *MockServiceRepository) CreateService(name, color, logo string) error {
	args := m.Called(name, color, logo)
	return args.Error(0)
}

func (m *MockServiceRepository) FindServiceById(id string) (entities.Service, error) {
	args := m.Called(id)
	return args.Get(0).(entities.Service), args.Error(1)
}

func (m *MockServiceRepository) FindServiceByName(name string) (entities.Service, error) {
	args := m.Called(name)
	return args.Get(0).(entities.Service), args.Error(1)
}

func (m *MockServiceRepository) FindAllServices() ([]entities.Service, error) {
	args := m.Called()
	return args.Get(0).([]entities.Service), args.Error(1)
}

func (m *MockServiceRepository) FindActionsServices() ([]entities.Service, error) {
	args := m.Called()
	return args.Get(0).([]entities.Service), args.Error(1)
}

func (m *MockServiceRepository) FindReactionsServices() ([]entities.Service, error) {
	args := m.Called()
	return args.Get(0).([]entities.Service), args.Error(1)
}

type MockUserServiceRepository struct {
	mock.Mock
}

func (m *MockUserServiceRepository) CreateUserService(userId, token, tokenRefresh, expiryDate, serviceId string) error {
	args := m.Called(userId, token, tokenRefresh, expiryDate, serviceId)
	return args.Error(0)
}

func (m *MockUserServiceRepository) FindUserServiceByServiceIdandUserId(userId, serviceId string) (entities.UserService, error) {
	args := m.Called(userId, serviceId)
	return args.Get(0).(entities.UserService), args.Error(1)
}

func (m *MockUserServiceRepository) UpdateUserServiceByServiceIdAndUserId(userId, accessToken, refreshToken, expiryDate, serviceId string) error {
	args := m.Called(userId, accessToken, refreshToken, expiryDate, serviceId)
	return args.Error(0)
}

func (m *MockUserServiceRepository) DeleteUserServiceByUserId(userId string) error {
	args := m.Called(userId)
	return args.Error(0)
}

type MockServiceServiceRepository struct {
	mock.Mock
}

func (m *MockServiceServiceRepository) RetrieveActionsServices() ([]entities.Service, error) {
	args := m.Called()
	return args.Get(0).([]entities.Service), args.Error(1)
}

func (m *MockServiceServiceRepository) RetrieveReactionsServices() ([]entities.Service, error) {
	args := m.Called()
	return args.Get(0).([]entities.Service), args.Error(1)
}

func (m *MockServiceServiceRepository) FindServiceByName(name string) (entities.Service, error) {
	args := m.Called(name)
	return args.Get(0).(entities.Service), args.Error(1)
}

func (m *MockServiceServiceRepository) FindServiceById(id string) (entities.Service, error) {
	args := m.Called(id)
	return args.Get(0).(entities.Service), args.Error(1)
}

func (m *MockServiceServiceRepository) FindServiceByActionId(id string) (entities.Service, error) {
	args := m.Called(id)
	return args.Get(0).(entities.Service), args.Error(1)
}

func (m *MockServiceServiceRepository) FindServiceByReactionId(id string) (entities.Service, error) {
	args := m.Called(id)
	return args.Get(0).(entities.Service), args.Error(1)
}

func (m *MockServiceServiceRepository) FindAllServices() ([]entities.Service, error) {
	args := m.Called()
	return args.Get(0).([]entities.Service), args.Error(1)
}

func (m *MockServiceServiceRepository) RetrieveActionsFromService(serviceName string) ([]entities.Action, error) {
	args := m.Called(serviceName)
	return args.Get(0).([]entities.Action), args.Error(1)
}

func (m *MockServiceServiceRepository) RetrieveReactionsFromService(serviceName string) ([]entities.Reaction, error) {
	args := m.Called(serviceName)
	return args.Get(0).([]entities.Reaction), args.Error(1)
}

func (m *MockServiceServiceRepository) GetGoogleRefreshTokenRequest(refreshToken string) (*http.Request, error) {
	args := m.Called(refreshToken)
	return args.Get(0).(*http.Request), args.Error(1)
}

func (m *MockServiceServiceRepository) GetSpotifyRefreshTokenRequest(refreshToken string) (*http.Request, error) {
	args := m.Called(refreshToken)
	return args.Get(0).(*http.Request), args.Error(1)
}

func (m *MockServiceServiceRepository) GetDiscordRefreshTokenRequest(refreshToken string) (*http.Request, error) {
	args := m.Called(refreshToken)
	return args.Get(0).(*http.Request), args.Error(1)
}

func (m *MockServiceServiceRepository) GetRedditRefreshTokenRequest(refreshToken string) (*http.Request, error) {
	args := m.Called(refreshToken)
	return args.Get(0).(*http.Request), args.Error(1)
}

func (m *MockServiceServiceRepository) GetAsanaRefreshTokenRequest(refreshToken string) (*http.Request, error) {
	args := m.Called(refreshToken)
	return args.Get(0).(*http.Request), args.Error(1)
}

func (m *MockServiceServiceRepository) GetDropboxRefreshTokenRequest(refreshToken string) (*http.Request, error) {
	args := m.Called(refreshToken)
	return args.Get(0).(*http.Request), args.Error(1)
}

func (m *MockServiceServiceRepository) GetGitlabRefreshTokenRequest(refreshToken string) (*http.Request, error) {
	args := m.Called(refreshToken)
	return args.Get(0).(*http.Request), args.Error(1)
}

func (m *MockServiceServiceRepository) ExecuteRequest(request *http.Request) (*http.Response, error) {
	args := m.Called(request)
	return args.Get(0).(*http.Response), args.Error(1)
}

func (m *MockServiceServiceRepository) ExecuteApiRequest(url, method, typeToken, accessToken string, body io.Reader) (*http.Response, error) {
	args := m.Called(url, method, typeToken, accessToken, body)
	return args.Get(0).(*http.Response), args.Error(1)
}

func (m *MockServiceServiceRepository) GetResultTokenFromCode(code, serviceName, callbackType, appType string) (entities.ResultToken, error) {
	args := m.Called(code, serviceName, callbackType, appType)
	return args.Get(0).(entities.ResultToken), args.Error(1)
}

func (m *MockServiceServiceRepository) GetUserInfoFromService(accessToken, serviceName string) (entities.UserInfo, error) {
	args := m.Called(accessToken, serviceName)
	return args.Get(0).(entities.UserInfo), args.Error(1)
}

func (m *MockServiceServiceRepository) RequestToTimeApi() (entities.TimeResponse, error) {
	args := m.Called()
	return args.Get(0).(entities.TimeResponse), args.Error(1)
}

func (m *MockServiceServiceRepository) OAuth2Service(serviceName, callbackType, appType string) (string, error) {
	args := m.Called(serviceName, callbackType, appType)
	return args.String(0), args.Error(1)
}

func (m *MockServiceServiceRepository) RequestGithubUserRepositories(accessToken string) ([]entities.GithubRepository, error) {
	args := m.Called(accessToken)
	return args.Get(0).([]entities.GithubRepository), args.Error(1)
}

func (m *MockServiceServiceRepository) RequestGitlabUserProjects(accessToken string) ([]entities.GitlabProject, error) {
	args := m.Called(accessToken)
	return args.Get(0).([]entities.GitlabProject), args.Error(1)
}

func (m *MockServiceServiceRepository) RetrieveDiscordGuildChannels(guildId string) ([]map[string]interface{}, error) {
	args := m.Called(guildId)
	return args.Get(0).([]map[string]interface{}), args.Error(1)
}
