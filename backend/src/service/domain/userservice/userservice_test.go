package userservice_service

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"backend/src/entities"
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

func TestGetUser(test *testing.T) {
	test.Run("Successful", func(test *testing.T) {
		var user entities.User

		mockUserRepo := new(MockUserRepository)
		userService := &UserServiceService{
			UserRepository: mockUserRepo,
		}

		user.Email = "test@test.com"

		mockUserRepo.On("FindUserByEmail", "test@test.com", "basic").
			Return(user, nil)

		_, err := userService.GetUser("test@test.com", "basic")
		require.NoError(test, err)
	})

	test.Run("User not found", func(test *testing.T) {
		var user entities.User

		mockUserRepo := new(MockUserRepository)
		userService := &UserServiceService{
			UserRepository: mockUserRepo,
		}

		user.Email = "test@test.com"

		mockUserRepo.On("FindUserByEmail", "test@test.com", "basic").
			Return(user, errors.New("User not found"))

		_, err := userService.GetUser("test@test.com", "basic")
		require.EqualError(test, err, "Could not find requested user")
	})
}

func TestRetrieveUserServiceAuthenticationStatus(test *testing.T) {
	test.Run("Successful", func(test *testing.T) {
		var user entities.User
		var service entities.Service

		mockUserRepo := new(MockUserRepository)
		mockServiceRepo := new(MockServiceRepository)
		mockUserServiceRepo := new(MockUserServiceRepository)

		userService := &UserServiceService{
			UserRepository:        mockUserRepo,
			ServiceRepository:     mockServiceRepo,
			UserServiceRepository: mockUserServiceRepo,
		}

		user.Email = "test@test.com"
		user.Id = "1"

		service.Id = "2"
		service.IsAuthNeeded = true

		mockUserRepo.On("FindUserByEmail", "test@test.com", "basic").
			Return(user, nil)

		mockServiceRepo.On("FindServiceByName", "Github").
			Return(service, nil)

		mockUserServiceRepo.On("FindUserServiceByServiceIdandUserId", "1", "2").
			Return(entities.UserService{}, nil)

		_, err := userService.RetrieveUserServiceAuthenticationStatus("test@test.com", "basic", "Github")
		require.NoError(test, err)
	})

	test.Run("User not found", func(test *testing.T) {
		var user entities.User
		var service entities.Service

		mockUserRepo := new(MockUserRepository)
		mockServiceRepo := new(MockServiceRepository)
		mockUserServiceRepo := new(MockUserServiceRepository)

		userService := &UserServiceService{
			UserRepository:        mockUserRepo,
			ServiceRepository:     mockServiceRepo,
			UserServiceRepository: mockUserServiceRepo,
		}

		user.Email = "test@test.com"
		user.Id = "1"

		service.Id = "2"
		service.IsAuthNeeded = true

		mockUserRepo.On("FindUserByEmail", "test@test.com", "basic").
			Return(user, errors.New("Could not find requested user"))

		mockServiceRepo.On("FindServiceByName", "Github").
			Return(service, nil)

		mockUserServiceRepo.On("FindUserServiceByServiceIdandUserId", "1", "2").
			Return(entities.UserService{}, nil)

		_, err := userService.RetrieveUserServiceAuthenticationStatus("test@test.com", "basic", "Github")
		require.EqualError(test, err, "Could not find requested user")
	})

	test.Run("Unknown service", func(test *testing.T) {
		var user entities.User
		var service entities.Service

		mockUserRepo := new(MockUserRepository)
		mockServiceRepo := new(MockServiceRepository)
		mockUserServiceRepo := new(MockUserServiceRepository)

		userService := &UserServiceService{
			UserRepository:        mockUserRepo,
			ServiceRepository:     mockServiceRepo,
			UserServiceRepository: mockUserServiceRepo,
		}

		user.Email = "test@test.com"
		user.Id = "1"

		service.Id = "2"
		service.IsAuthNeeded = true

		mockUserRepo.On("FindUserByEmail", "test@test.com", "basic").
			Return(user, nil)

		mockServiceRepo.On("FindServiceByName", "Github").
			Return(service, errors.New("Unknown service"))

		mockUserServiceRepo.On("FindUserServiceByServiceIdandUserId", "1", "2").
			Return(entities.UserService{}, nil)

		_, err := userService.RetrieveUserServiceAuthenticationStatus("test@test.com", "basic", "Github")
		require.EqualError(test, err, "Unknown service")
	})

	test.Run("User service not found", func(test *testing.T) {
		var user entities.User
		var service entities.Service

		mockUserRepo := new(MockUserRepository)
		mockServiceRepo := new(MockServiceRepository)
		mockUserServiceRepo := new(MockUserServiceRepository)

		userService := &UserServiceService{
			UserRepository:        mockUserRepo,
			ServiceRepository:     mockServiceRepo,
			UserServiceRepository: mockUserServiceRepo,
		}

		user.Email = "test@test.com"
		user.Id = "1"

		service.Id = "2"
		service.IsAuthNeeded = true

		mockUserRepo.On("FindUserByEmail", "test@test.com", "basic").
			Return(user, nil)

		mockServiceRepo.On("FindServiceByName", "Github").
			Return(service, nil)

		mockUserServiceRepo.On("FindUserServiceByServiceIdandUserId", "1", "2").
			Return(entities.UserService{}, errors.New("User Service not found"))

		answer, err := userService.RetrieveUserServiceAuthenticationStatus("test@test.com", "basic", "Github")

		require.NoError(test, err)
		require.False(test, answer)
	})
}

func TestRefreshToken(test *testing.T) {
	mockService := entities.Service{
		Id:   "1",
		Name: "Google",
	}

	test.Run("Successful", func(test *testing.T) {
		mockServiceService := new(MockServiceServiceRepository)
		mockServiceRepo := new(MockServiceRepository)
		mockUserServiceRepo := new(MockUserServiceRepository)

		userService := &UserServiceService{
			ServiceService:        mockServiceService,
			ServiceRepository:     mockServiceRepo,
			UserServiceRepository: mockUserServiceRepo,
		}

		mockRequest := &http.Request{}
		mockResponse := &http.Response{
			Body: io.NopCloser(strings.NewReader(`{
				"access_token": "accessToken",
				"refresh_token": "refreshToken",
				"expires_in": 3600}`,
			)),
		}

		mockServiceService.On("GetGoogleRefreshTokenRequest", "refreshToken").
			Return(mockRequest, nil)

		mockServiceService.On("ExecuteRequest", mockRequest).
			Return(mockResponse, nil)

		mockServiceRepo.On("FindServiceByName", "Google").
			Return(mockService, nil)

		mockUserServiceRepo.On("UpdateUserServiceByServiceIdAndUserId", "1", "accessToken", "refreshToken", mock.Anything, mockService.Id).
			Return(nil)

		_, err := userService.refreshToken("refreshToken", "1", "Google")

		require.NoError(test, err)
	})

	test.Run("Unknown service", func(test *testing.T) {
		mockServiceService := new(MockServiceServiceRepository)
		mockServiceRepo := new(MockServiceRepository)
		mockUserServiceRepo := new(MockUserServiceRepository)

		userService := &UserServiceService{
			ServiceService:        mockServiceService,
			ServiceRepository:     mockServiceRepo,
			UserServiceRepository: mockUserServiceRepo,
		}

		_, err := userService.refreshToken("refreshToken", "1", "false")

		require.EqualError(test, err, "Unknown service")
	})

	test.Run("Request error", func(test *testing.T) {
		mockServiceService := new(MockServiceServiceRepository)
		mockServiceRepo := new(MockServiceRepository)
		mockUserServiceRepo := new(MockUserServiceRepository)

		userService := &UserServiceService{
			ServiceService:        mockServiceService,
			ServiceRepository:     mockServiceRepo,
			UserServiceRepository: mockUserServiceRepo,
		}

		mockRequest := &http.Request{}

		mockServiceService.On("GetGoogleRefreshTokenRequest", "refreshToken").
			Return(mockRequest, errors.New("request error"))

		_, err := userService.refreshToken("refreshToken", "1", "Google")

		require.EqualError(test, err, "request error")
	})

	test.Run("Fail Execute", func(test *testing.T) {
		mockServiceService := new(MockServiceServiceRepository)
		mockServiceRepo := new(MockServiceRepository)
		mockUserServiceRepo := new(MockUserServiceRepository)

		userService := &UserServiceService{
			ServiceService:        mockServiceService,
			ServiceRepository:     mockServiceRepo,
			UserServiceRepository: mockUserServiceRepo,
		}

		mockRequest := &http.Request{}
		mockResponse := &http.Response{
			Body: io.NopCloser(strings.NewReader(`{
			"access_token": "accessToken",
			"refresh_token": "refreshToken",
			"expires_in": 3600}`,
			)),
		}

		mockServiceService.On("GetGoogleRefreshTokenRequest", "refreshToken").
			Return(mockRequest, nil)

		mockServiceService.On("ExecuteRequest", mockRequest).
			Return(mockResponse, errors.New("Fail execute"))

		_, err := userService.refreshToken("refreshToken", "1", "Google")

		require.EqualError(test, err, "Fail execute")
	})

	test.Run("Fail Decode", func(t *testing.T) {
		mockServiceService := new(MockServiceServiceRepository)
		mockServiceRepo := new(MockServiceRepository)
		mockUserServiceRepo := new(MockUserServiceRepository)

		userService := &UserServiceService{
			ServiceService:        mockServiceService,
			ServiceRepository:     mockServiceRepo,
			UserServiceRepository: mockUserServiceRepo,
		}

		mockRequest := &http.Request{}
		mockResponse := &http.Response{
			Body: io.NopCloser(strings.NewReader(`invalid json`)),
		}

		mockServiceService.On("GetGoogleRefreshTokenRequest", "refreshToken").
			Return(mockRequest, nil)

		mockServiceService.On("ExecuteRequest", mockRequest).
			Return(mockResponse, nil)

		_, err := userService.refreshToken("refreshToken", "1", "Google")

		require.Error(test, err)
	})

	test.Run("Fail Find Service", func(t *testing.T) {
		mockServiceService := new(MockServiceServiceRepository)
		mockServiceRepo := new(MockServiceRepository)
		mockUserServiceRepo := new(MockUserServiceRepository)

		userService := &UserServiceService{
			ServiceService:        mockServiceService,
			ServiceRepository:     mockServiceRepo,
			UserServiceRepository: mockUserServiceRepo,
		}

		mockRequest := &http.Request{}
		mockResponse := &http.Response{
			Body: io.NopCloser(strings.NewReader(`{
			"access_token": "accessToken",
			"refresh_token": "refreshToken",
			"expires_in": 3600}`,
			)),
		}

		mockServiceService.On("GetGoogleRefreshTokenRequest", "refreshToken").
			Return(mockRequest, nil)

		mockServiceService.On("ExecuteRequest", mockRequest).
			Return(mockResponse, nil)

		mockServiceRepo.On("FindServiceByName", "Google").
			Return(mockService, errors.New("Fail find service"))

		_, err := userService.refreshToken("refreshToken", "1", "Google")

		require.EqualError(test, err, "Fail find service")
	})

	test.Run("Fail update", func(t *testing.T) {
		mockServiceService := new(MockServiceServiceRepository)
		mockServiceRepo := new(MockServiceRepository)
		mockUserServiceRepo := new(MockUserServiceRepository)

		userService := &UserServiceService{
			ServiceService:        mockServiceService,
			ServiceRepository:     mockServiceRepo,
			UserServiceRepository: mockUserServiceRepo,
		}

		mockRequest := &http.Request{}
		mockResponse := &http.Response{
			Body: io.NopCloser(strings.NewReader(`{
			"access_token": "accessToken",
			"refresh_token": "refreshToken",
			"expires_in": 3600}`,
			)),
		}

		mockServiceService.On("GetGoogleRefreshTokenRequest", "refreshToken").
			Return(mockRequest, nil)

		mockServiceService.On("ExecuteRequest", mockRequest).
			Return(mockResponse, nil)

		mockServiceRepo.On("FindServiceByName", "Google").
			Return(mockService, nil)

		mockUserServiceRepo.On("UpdateUserServiceByServiceIdAndUserId", "1", "accessToken", "refreshToken", mock.Anything, mockService.Id).
			Return(errors.New("Fail update service"))

		_, err := userService.refreshToken("refreshToken", "1", "Google")

		require.EqualError(test, err, "Fail update service")
	})
}

func TestCallApiAndRefresh(test *testing.T) {
	test.Run("User not found", func(test *testing.T) {
		mockUserRepo := new(MockUserRepository)
		mockServiceRepo := new(MockServiceRepository)
		mockUserServiceRepo := new(MockUserServiceRepository)
		mockServiceService := new(MockServiceServiceRepository)

		userService := &UserServiceService{
			UserRepository:        mockUserRepo,
			ServiceRepository:     mockServiceRepo,
			UserServiceRepository: mockUserServiceRepo,
			ServiceService:        mockServiceService,
		}

		mockUserRepo.On("FindUserByEmail", "test@test.com", "basic").
			Return(entities.User{}, errors.New("user not found"))

		_, err := userService.CallApiAndRefresh("test@test.com", "basic", "Google")

		require.EqualError(test, err, "Could not find requested user")
	})

	test.Run("Fail Retrieving Service", func(t *testing.T) {
		var user entities.User

		mockUserRepo := new(MockUserRepository)
		mockServiceRepo := new(MockServiceRepository)
		mockUserServiceRepo := new(MockUserServiceRepository)
		mockServiceService := new(MockServiceServiceRepository)

		userService := &UserServiceService{
			UserRepository:        mockUserRepo,
			ServiceRepository:     mockServiceRepo,
			UserServiceRepository: mockUserServiceRepo,
			ServiceService:        mockServiceService,
		}

		user.Email = "test@test.com"
		user.Id = "1"

		mockUserRepo.On("FindUserByEmail", "test@test.com", "basic").
			Return(user, nil)

		mockServiceRepo.On("FindServiceByName", "Google").
			Return(entities.Service{}, errors.New("service not found"))

		_, err := userService.CallApiAndRefresh("test@test.com", "basic", "Google")

		require.EqualError(t, err, "service not found")
	})

	test.Run("Fail Retrieving User service", func(test *testing.T) {
		var user entities.User
		var service entities.Service

		mockUserRepo := new(MockUserRepository)
		mockServiceRepo := new(MockServiceRepository)
		mockUserServiceRepo := new(MockUserServiceRepository)
		mockServiceService := new(MockServiceServiceRepository)

		userService := &UserServiceService{
			UserRepository:        mockUserRepo,
			ServiceRepository:     mockServiceRepo,
			UserServiceRepository: mockUserServiceRepo,
			ServiceService:        mockServiceService,
		}

		user.Email = "test@test.com"
		user.Id = "1"

		service.Id = "1"

		mockUserRepo.On("FindUserByEmail", "test@test.com", "basic").
			Return(user, nil)

		mockServiceRepo.On("FindServiceByName", "Google").
			Return(service, nil)

		mockUserServiceRepo.On("FindUserServiceByServiceIdandUserId", "1", "1").
			Return(entities.UserService{}, errors.New("user service not found"))

		_, err := userService.CallApiAndRefresh("test@test.com", "basic", "Google")

		require.EqualError(test, err, "user service not found")
	})

	test.Run("Refresh successful", func(test *testing.T) {
		var user entities.User
		var service entities.Service
		var serviceOfUser entities.UserService

		mockUserRepo := new(MockUserRepository)
		mockServiceRepo := new(MockServiceRepository)
		mockUserServiceRepo := new(MockUserServiceRepository)
		mockServiceService := new(MockServiceServiceRepository)

		userService := &UserServiceService{
			UserRepository:        mockUserRepo,
			ServiceRepository:     mockServiceRepo,
			UserServiceRepository: mockUserServiceRepo,
			ServiceService:        mockServiceService,
		}

		user.Email = "test@test.com"
		user.Id = "1"

		service.Id = "1"

		serviceOfUser.AccessToken = "accessToken"
		serviceOfUser.RefreshToken = "refreshToken"
		serviceOfUser.ExpiryDate = time.Now().Add(-time.Hour).Format(formattingDate)

		mockUserRepo.On("FindUserByEmail", "test@test.com", "basic").
			Return(user, nil)

		mockServiceRepo.On("FindServiceByName", "Google").
			Return(service, nil)

		mockUserServiceRepo.On("FindUserServiceByServiceIdandUserId", "1", "1").
			Return(serviceOfUser, nil)

		mockRequest := &http.Request{}
		mockResponse := &http.Response{
			Body: io.NopCloser(strings.NewReader(`{
				"access_token": "accessToken",
				"refresh_token": "refreshToken",
				"expires_in": 3600
			}`)),
		}

		mockServiceService.On("GetGoogleRefreshTokenRequest", "refreshToken").
			Return(mockRequest, nil)

		mockServiceService.On("ExecuteRequest", mockRequest).
			Return(mockResponse, nil)

		mockUserServiceRepo.On("UpdateUserServiceByServiceIdAndUserId", "1", "accessToken", "refreshToken", mock.Anything, "1").
			Return(nil)

		_, err := userService.CallApiAndRefresh("test@test.com", "basic", "Google")

		require.NoError(test, err)
	})

	test.Run("Fail refresh", func(test *testing.T) {
		var user entities.User
		var service entities.Service
		var serviceOfUser entities.UserService

		mockUserRepo := new(MockUserRepository)
		mockServiceRepo := new(MockServiceRepository)
		mockUserServiceRepo := new(MockUserServiceRepository)
		mockServiceService := new(MockServiceServiceRepository)

		userService := &UserServiceService{
			UserRepository:        mockUserRepo,
			ServiceRepository:     mockServiceRepo,
			UserServiceRepository: mockUserServiceRepo,
			ServiceService:        mockServiceService,
		}

		user.Email = "test@test.com"
		user.Id = "1"

		service.Id = "1"

		serviceOfUser.AccessToken = "accessToken"
		serviceOfUser.RefreshToken = "refreshToken"
		serviceOfUser.ExpiryDate = time.Now().Add(-time.Hour).Format(formattingDate)

		mockUserRepo.On("FindUserByEmail", "test@test.com", "basic").
			Return(user, nil)

		mockServiceRepo.On("FindServiceByName", "Google").
			Return(service, nil)

		mockUserServiceRepo.On("FindUserServiceByServiceIdandUserId", "1", "1").
			Return(serviceOfUser, nil)

		mockRequest := &http.Request{}
		mockResponse := &http.Response{
			Body: io.NopCloser(strings.NewReader(`{
					"access_token": "accessToken",
					"refresh_token": "refreshToken",
					"expires_in": 3600
				}`)),
		}

		mockServiceService.On("GetGoogleRefreshTokenRequest", "refreshToken").
			Return(mockRequest, nil)

		mockServiceService.On("ExecuteRequest", mock.Anything).
			Return(mockResponse, errors.New("refresh token failed"))

		_, err := userService.CallApiAndRefresh("test@test.com", "basic", "Google")

		require.EqualError(test, err, "refresh token failed")
	})

	test.Run("Token not expired", func(test *testing.T) {
		var user entities.User
		var service entities.Service
		var serviceOfUser entities.UserService

		mockUserRepo := new(MockUserRepository)
		mockServiceRepo := new(MockServiceRepository)
		mockUserServiceRepo := new(MockUserServiceRepository)
		mockServiceService := new(MockServiceServiceRepository)

		userService := &UserServiceService{
			UserRepository:        mockUserRepo,
			ServiceRepository:     mockServiceRepo,
			UserServiceRepository: mockUserServiceRepo,
			ServiceService:        mockServiceService,
		}

		user.Email = "test@test.com"
		user.Id = "1"

		service.Id = "1"

		serviceOfUser.AccessToken = "accessToken"
		serviceOfUser.RefreshToken = "refreshToken"
		serviceOfUser.ExpiryDate = time.Now().Add(time.Hour).Format(formattingDate)

		mockUserRepo.On("FindUserByEmail", "test@test.com", "basic").
			Return(user, nil)

		mockServiceRepo.On("FindServiceByName", "Google").
			Return(service, nil)

		mockUserServiceRepo.On("FindUserServiceByServiceIdandUserId", "1", "1").
			Return(serviceOfUser, nil)

		_, err := userService.CallApiAndRefresh("test@test.com", "basic", "Google")

		require.NoError(test, err)
	})
}

func TestUpdateTokenForService(test *testing.T) {
	test.Run("Successful", func(test *testing.T) {
		var user entities.User
		var service entities.Service
		var token entities.ResultToken

		mockUserRepo := new(MockUserRepository)
		mockServiceRepo := new(MockServiceRepository)
		mockUserServiceRepo := new(MockUserServiceRepository)
		mockServiceService := new(MockServiceServiceRepository)

		userService := &UserServiceService{
			UserRepository:        mockUserRepo,
			ServiceRepository:     mockServiceRepo,
			UserServiceRepository: mockUserServiceRepo,
			ServiceService:        mockServiceService,
		}

		user.Email = "test@test.com"
		user.Id = "1"

		service.Id = "2"
		service.Name = "Google"

		token.AccessToken = "accessToken"
		token.RefreshToken = "refreshToken"
		token.ExpiresIn = 3600

		mockServiceRepo.On("FindServiceByName", "Google").
			Return(service, nil)

		mockServiceService.On("GetResultTokenFromCode", "code", "Google", "service", "web").
			Return(token, nil)

		mockUserRepo.On("FindUserByEmail", "test@test.com", "basic").
			Return(user, nil)

		mockServiceRepo.On("FindServiceByName", "Google").
			Return(service, nil)

		mockUserServiceRepo.On("FindUserServiceByServiceIdandUserId", "1", "2").
			Return(entities.UserService{}, errors.New("user service not found"))

		mockUserServiceRepo.On("CreateUserService", "1", "accessToken", "refreshToken", mock.Anything, "2").
			Return(nil)

		err := userService.UpdateTokenForService("code", "Google", "web", "test@test.com", "basic")
		require.NoError(test, err)
	})

	test.Run("Service not found", func(test *testing.T) {
		mockUserRepo := new(MockUserRepository)
		mockServiceRepo := new(MockServiceRepository)
		mockUserServiceRepo := new(MockUserServiceRepository)
		mockServiceService := new(MockServiceServiceRepository)

		userService := &UserServiceService{
			UserRepository:        mockUserRepo,
			ServiceRepository:     mockServiceRepo,
			UserServiceRepository: mockUserServiceRepo,
			ServiceService:        mockServiceService,
		}

		mockServiceRepo.On("FindServiceByName", "Google").
			Return(entities.Service{}, fmt.Errorf("service not found"))

		err := userService.UpdateTokenForService("code", "Google", "web", "test@test.com", "basic")
		require.EqualError(test, err, "service not found")
	})

	test.Run("Fail Token from Code", func(test *testing.T) {
		var service entities.Service

		mockUserRepo := new(MockUserRepository)
		mockServiceRepo := new(MockServiceRepository)
		mockUserServiceRepo := new(MockUserServiceRepository)
		mockServiceService := new(MockServiceServiceRepository)

		userService := &UserServiceService{
			UserRepository:        mockUserRepo,
			ServiceRepository:     mockServiceRepo,
			UserServiceRepository: mockUserServiceRepo,
			ServiceService:        mockServiceService,
		}

		mockServiceRepo.On("FindServiceByName", "Google").
			Return(service, nil)

		mockServiceService.On("GetResultTokenFromCode", "code", "Google", "service", "web").
			Return(entities.ResultToken{}, errors.New("token retrieval failed"))

		err := userService.UpdateTokenForService("code", "Google", "web", "test@test.com", "basic")

		require.EqualError(test, err, "token retrieval failed")
	})

	test.Run("User not found", func(test *testing.T) {
		var service entities.Service
		var token entities.ResultToken

		mockUserRepo := new(MockUserRepository)
		mockServiceRepo := new(MockServiceRepository)
		mockUserServiceRepo := new(MockUserServiceRepository)
		mockServiceService := new(MockServiceServiceRepository)

		userService := &UserServiceService{
			UserRepository:        mockUserRepo,
			ServiceRepository:     mockServiceRepo,
			UserServiceRepository: mockUserServiceRepo,
			ServiceService:        mockServiceService,
		}

		mockServiceRepo.On("FindServiceByName", "Google").
			Return(service, nil)

		mockServiceService.On("GetResultTokenFromCode", "code", "Google", "service", "web").
			Return(token, nil)

		mockUserRepo.On("FindUserByEmail", "test@test.com", "basic").
			Return(entities.User{}, errors.New("user not found"))

		err := userService.UpdateTokenForService("code", "Google", "web", "test@test.com", "basic")

		require.EqualError(test, err, "Could not find requested user")
	})
}

func TestRetrieveGithubUserRepositories(test *testing.T) {
	var user entities.User
	var service entities.Service
	var serviceOfUser entities.UserService

	user.Email = "test@test.com"
	user.Id = "1"

	service.Id = "1"

	serviceOfUser.AccessToken = "accessToken"
	serviceOfUser.RefreshToken = "refreshToken"
	serviceOfUser.ExpiryDate = time.Now().Add(time.Hour).Format(formattingDate)

	test.Run("Successful", func(test *testing.T) {
		mockUserRepo := new(MockUserRepository)
		mockServiceRepo := new(MockServiceRepository)
		mockUserServiceRepo := new(MockUserServiceRepository)
		mockServiceService := new(MockServiceServiceRepository)

		userService := &UserServiceService{
			UserRepository:        mockUserRepo,
			ServiceRepository:     mockServiceRepo,
			UserServiceRepository: mockUserServiceRepo,
			ServiceService:        mockServiceService,
		}

		mockUserRepo.On("FindUserByEmail", "test@test.com", "basic").
			Return(user, nil)

		mockServiceRepo.On("FindServiceByName", "Github").
			Return(service, nil)

		mockUserServiceRepo.On("FindUserServiceByServiceIdandUserId", "1", "1").
			Return(serviceOfUser, nil)

		mockServiceService.On("RequestGithubUserRepositories", "accessToken").
			Return([]entities.GithubRepository{{}}, nil)

		_, err := userService.RetrieveGithubUserRepositories("test@test.com", "basic")

		require.NoError(test, err)
	})

	test.Run("Call api fail", func(test *testing.T) {
		mockUserRepo := new(MockUserRepository)
		mockServiceRepo := new(MockServiceRepository)
		mockUserServiceRepo := new(MockUserServiceRepository)
		mockServiceService := new(MockServiceServiceRepository)

		userService := &UserServiceService{
			UserRepository:        mockUserRepo,
			ServiceRepository:     mockServiceRepo,
			UserServiceRepository: mockUserServiceRepo,
			ServiceService:        mockServiceService,
		}

		mockUserRepo.On("FindUserByEmail", "test@test.com", "basic").
			Return(user, nil)

		mockServiceRepo.On("FindServiceByName", "Github").
			Return(service, nil)

		mockUserServiceRepo.On("FindUserServiceByServiceIdandUserId", "1", "1").
			Return(serviceOfUser, errors.New("Call api refresh fail"))

		mockServiceService.On("RequestGithubUserRepositories", "accessToken").
			Return([]entities.GithubRepository{{}}, nil)

		_, err := userService.RetrieveGithubUserRepositories("test@test.com", "basic")

		require.EqualError(test, err, "Call api refresh fail")
	})
}

func TestRetrieveGitlabUserProjects(test *testing.T) {
	var user entities.User
	var service entities.Service
	var serviceOfUser entities.UserService

	user.Email = "test@test.com"
	user.Id = "1"

	service.Id = "1"

	serviceOfUser.AccessToken = "accessToken"
	serviceOfUser.RefreshToken = "refreshToken"
	serviceOfUser.ExpiryDate = time.Now().Add(time.Hour).Format(formattingDate)

	test.Run("Successful", func(test *testing.T) {
		mockUserRepo := new(MockUserRepository)
		mockServiceRepo := new(MockServiceRepository)
		mockUserServiceRepo := new(MockUserServiceRepository)
		mockServiceService := new(MockServiceServiceRepository)

		userService := &UserServiceService{
			UserRepository:        mockUserRepo,
			ServiceRepository:     mockServiceRepo,
			UserServiceRepository: mockUserServiceRepo,
			ServiceService:        mockServiceService,
		}

		mockUserRepo.On("FindUserByEmail", "test@test.com", "basic").
			Return(user, nil)

		mockServiceRepo.On("FindServiceByName", "Gitlab").
			Return(service, nil)

		mockUserServiceRepo.On("FindUserServiceByServiceIdandUserId", "1", "1").
			Return(serviceOfUser, nil)

		mockServiceService.On("RequestGitlabUserProjects", "accessToken").
			Return([]entities.GitlabProject{{}}, nil)

		_, err := userService.RetrieveGitlabUserProjects("test@test.com", "basic")

		require.NoError(test, err)
	})

	test.Run("Call api fail", func(test *testing.T) {
		mockUserRepo := new(MockUserRepository)
		mockServiceRepo := new(MockServiceRepository)
		mockUserServiceRepo := new(MockUserServiceRepository)
		mockServiceService := new(MockServiceServiceRepository)

		userService := &UserServiceService{
			UserRepository:        mockUserRepo,
			ServiceRepository:     mockServiceRepo,
			UserServiceRepository: mockUserServiceRepo,
			ServiceService:        mockServiceService,
		}

		mockUserRepo.On("FindUserByEmail", "test@test.com", "basic").
			Return(user, nil)

		mockServiceRepo.On("FindServiceByName", "Gitlab").
			Return(service, nil)

		mockUserServiceRepo.On("FindUserServiceByServiceIdandUserId", "1", "1").
			Return(serviceOfUser, errors.New("Call api refresh fail"))

		mockServiceService.On("RequestGitlabUserProjects", "accessToken").
			Return([]entities.GitlabProject{{}}, nil)

		_, err := userService.RetrieveGitlabUserProjects("test@test.com", "basic")

		require.EqualError(test, err, "Call api refresh fail")
	})
}

func TestRetrieveDiscordUserServers(test *testing.T) {
	var user entities.User
	var service entities.Service
	var serviceOfUser entities.UserService

	user.Email = "test@test.com"
	user.Id = "1"

	service.Id = "1"

	serviceOfUser.AccessToken = "accessToken"
	serviceOfUser.RefreshToken = "refreshToken"
	serviceOfUser.ExpiryDate = time.Now().Add(time.Hour).Format(formattingDate)

	test.Run("Successful", func(test *testing.T) {
		mockResponse := &http.Response{
			Body: io.NopCloser(strings.NewReader(`
				[
					{
						"id":"460797762129231873",
						"name":"Tek-KMU",
						"icon":"0fbf50ba2c34c12324417f02b25e27ef",
						"banner":null,
						"owner":false,
						"permissions":"2222085186641473",
						"features":["CHANNEL_ICON_EMOJIS_GENERATED","SOUNDBOARD"]
					}
				]`,
			))}

		mockUserRepo := new(MockUserRepository)
		mockServiceRepo := new(MockServiceRepository)
		mockUserServiceRepo := new(MockUserServiceRepository)
		mockServiceService := new(MockServiceServiceRepository)

		userService := &UserServiceService{
			UserRepository:        mockUserRepo,
			ServiceRepository:     mockServiceRepo,
			UserServiceRepository: mockUserServiceRepo,
			ServiceService:        mockServiceService,
		}

		mockUserRepo.On("FindUserByEmail", "test@test.com", "basic").
			Return(user, nil)

		mockServiceRepo.On("FindServiceByName", "Discord").
			Return(service, nil)

		mockUserServiceRepo.On("FindUserServiceByServiceIdandUserId", "1", "1").
			Return(serviceOfUser, nil)

		mockServiceService.On("ExecuteApiRequest", "https://discord.com/api/v10/users/@me/guilds", "GET", "Bearer ", "accessToken", nil).
			Return(mockResponse, nil)

		_, err := userService.RetrieveDiscordUserServers("test@test.com", "basic")

		require.NoError(test, err)
	})

	test.Run("Fail API Refresh", func(test *testing.T) {
		mockUserRepo := new(MockUserRepository)
		mockServiceRepo := new(MockServiceRepository)
		mockUserServiceRepo := new(MockUserServiceRepository)

		userService := &UserServiceService{
			UserRepository:        mockUserRepo,
			ServiceRepository:     mockServiceRepo,
			UserServiceRepository: mockUserServiceRepo,
		}

		mockUserRepo.On("FindUserByEmail", "test@test.com", "basic").
			Return(user, nil)

		mockServiceRepo.On("FindServiceByName", "Discord").
			Return(service, nil)

		mockUserServiceRepo.On("FindUserServiceByServiceIdandUserId", "1", "1").
			Return(serviceOfUser, errors.New("Fail API Request"))

		_, err := userService.RetrieveDiscordUserServers("test@test.com", "basic")

		require.EqualError(test, err, "Fail API Request")
	})

	test.Run("Successful", func(test *testing.T) {
		mockResponse := &http.Response{
			Body: io.NopCloser(strings.NewReader(`
				[
					{
						"id":"460797762129231873",
						"name":"Tek-KMU",
						"icon":"0fbf50ba2c34c12324417f02b25e27ef",
						"banner":null,
						"owner":false,
						"permissions":"2222085186641473",
						"features":["CHANNEL_ICON_EMOJIS_GENERATED","SOUNDBOARD"]
					}
				]`,
			))}

		mockUserRepo := new(MockUserRepository)
		mockServiceRepo := new(MockServiceRepository)
		mockUserServiceRepo := new(MockUserServiceRepository)
		mockServiceService := new(MockServiceServiceRepository)

		userService := &UserServiceService{
			UserRepository:        mockUserRepo,
			ServiceRepository:     mockServiceRepo,
			UserServiceRepository: mockUserServiceRepo,
			ServiceService:        mockServiceService,
		}

		mockUserRepo.On("FindUserByEmail", "test@test.com", "basic").
			Return(user, nil)

		mockServiceRepo.On("FindServiceByName", "Discord").
			Return(service, nil)

		mockUserServiceRepo.On("FindUserServiceByServiceIdandUserId", "1", "1").
			Return(serviceOfUser, nil)

		mockServiceService.On("ExecuteApiRequest", "https://discord.com/api/v10/users/@me/guilds", "GET", "Bearer ", "accessToken", nil).
			Return(mockResponse, errors.New("Fail execute API"))

		_, err := userService.RetrieveDiscordUserServers("test@test.com", "basic")

		require.EqualError(test, err, "Fail execute API")
	})

	test.Run("Fail Decode", func(test *testing.T) {
		mockResponse := &http.Response{
			Body: io.NopCloser(strings.NewReader(`invalid json`)),
		}

		mockUserRepo := new(MockUserRepository)
		mockServiceRepo := new(MockServiceRepository)
		mockUserServiceRepo := new(MockUserServiceRepository)
		mockServiceService := new(MockServiceServiceRepository)

		userService := &UserServiceService{
			UserRepository:        mockUserRepo,
			ServiceRepository:     mockServiceRepo,
			UserServiceRepository: mockUserServiceRepo,
			ServiceService:        mockServiceService,
		}

		mockUserRepo.On("FindUserByEmail", "test@test.com", "basic").
			Return(user, nil)

		mockServiceRepo.On("FindServiceByName", "Discord").
			Return(service, nil)

		mockUserServiceRepo.On("FindUserServiceByServiceIdandUserId", "1", "1").
			Return(serviceOfUser, nil)

		mockServiceService.On("ExecuteApiRequest", "https://discord.com/api/v10/users/@me/guilds", "GET", "Bearer ", "accessToken", nil).
			Return(mockResponse, nil)

		_, err := userService.RetrieveDiscordUserServers("test@test.com", "basic")

		require.Error(test, err)
	})
}

func TestRetrieveAsanaUserWorkspaces(test *testing.T) {
	var user entities.User
	var service entities.Service
	var serviceOfUser entities.UserService

	user.Email = "test@test.com"
	user.Id = "1"

	service.Id = "1"

	serviceOfUser.AccessToken = "accessToken"
	serviceOfUser.RefreshToken = "refreshToken"
	serviceOfUser.ExpiryDate = time.Now().Add(time.Hour).Format(formattingDate)

	test.Run("Successful", func(test *testing.T) {
		mockResponse := &http.Response{
			Body: io.NopCloser(strings.NewReader(`
				{
					"data": [
						{
							"gid": "12345",
							"resource_type": "workspace",
							"name": "My Company Workspace"
						}
					],
					"next_page": {
						"offset": "eyJ0eXAiOJiKV1iQLCJhbGciOiJIUzI1NiJ9",
						"path": "/tasks/12345/attachments?limit=2&offset=eyJ0eXAiOJiKV1iQLCJhbGciOiJIUzI1NiJ9",
						"uri": "https://app.asana.com/api/1.0/tasks/12345/attachments?limit=2&offset=eyJ0eXAiOJiKV1iQLCJhbGciOiJIUzI1NiJ9"
					}
				}`,
			)),
		}

		mockUserRepo := new(MockUserRepository)
		mockServiceRepo := new(MockServiceRepository)
		mockUserServiceRepo := new(MockUserServiceRepository)
		mockServiceService := new(MockServiceServiceRepository)

		userService := &UserServiceService{
			UserRepository:        mockUserRepo,
			ServiceRepository:     mockServiceRepo,
			UserServiceRepository: mockUserServiceRepo,
			ServiceService:        mockServiceService,
		}

		mockUserRepo.On("FindUserByEmail", "test@test.com", "basic").
			Return(user, nil)

		mockServiceRepo.On("FindServiceByName", "Asana").
			Return(service, nil)

		mockUserServiceRepo.On("FindUserServiceByServiceIdandUserId", "1", "1").
			Return(serviceOfUser, nil)

		mockServiceService.On("ExecuteApiRequest", "https://app.asana.com/api/1.0/workspaces/", "GET", "Bearer ", "accessToken", nil).
			Return(mockResponse, nil)

		_, err := userService.RetrieveAsanaUserWorkspaces("test@test.com", "basic")

		require.NoError(test, err)
	})

	test.Run("Call API Refresh", func(test *testing.T) {
		mockUserRepo := new(MockUserRepository)
		mockServiceRepo := new(MockServiceRepository)
		mockUserServiceRepo := new(MockUserServiceRepository)

		userService := &UserServiceService{
			UserRepository:        mockUserRepo,
			ServiceRepository:     mockServiceRepo,
			UserServiceRepository: mockUserServiceRepo,
		}

		mockUserRepo.On("FindUserByEmail", "test@test.com", "basic").
			Return(user, nil)

		mockServiceRepo.On("FindServiceByName", "Asana").
			Return(service, nil)

		mockUserServiceRepo.On("FindUserServiceByServiceIdandUserId", "1", "1").
			Return(serviceOfUser, errors.New("Fail Api Refresh"))

		_, err := userService.RetrieveAsanaUserWorkspaces("test@test.com", "basic")

		require.EqualError(test, err, "Fail Api Refresh")
	})

	test.Run("Fail Execute Request", func(test *testing.T) {
		mockResponse := &http.Response{
			Body: io.NopCloser(strings.NewReader(`
				{
					"data": [
						{
							"gid": "12345",
							"resource_type": "workspace",
							"name": "My Company Workspace"
						}
					],
					"next_page": {
						"offset": "eyJ0eXAiOJiKV1iQLCJhbGciOiJIUzI1NiJ9",
						"path": "/tasks/12345/attachments?limit=2&offset=eyJ0eXAiOJiKV1iQLCJhbGciOiJIUzI1NiJ9",
						"uri": "https://app.asana.com/api/1.0/tasks/12345/attachments?limit=2&offset=eyJ0eXAiOJiKV1iQLCJhbGciOiJIUzI1NiJ9"
					}
				}`,
			)),
		}

		mockUserRepo := new(MockUserRepository)
		mockServiceRepo := new(MockServiceRepository)
		mockUserServiceRepo := new(MockUserServiceRepository)
		mockServiceService := new(MockServiceServiceRepository)

		userService := &UserServiceService{
			UserRepository:        mockUserRepo,
			ServiceRepository:     mockServiceRepo,
			UserServiceRepository: mockUserServiceRepo,
			ServiceService:        mockServiceService,
		}

		mockUserRepo.On("FindUserByEmail", "test@test.com", "basic").
			Return(user, nil)

		mockServiceRepo.On("FindServiceByName", "Asana").
			Return(service, nil)

		mockUserServiceRepo.On("FindUserServiceByServiceIdandUserId", "1", "1").
			Return(serviceOfUser, nil)

		mockServiceService.On("ExecuteApiRequest", "https://app.asana.com/api/1.0/workspaces/", "GET", "Bearer ", "accessToken", nil).
			Return(mockResponse, errors.New("Fail Execute Request"))

		_, err := userService.RetrieveAsanaUserWorkspaces("test@test.com", "basic")

		require.EqualError(test, err, "Fail Execute Request")
	})

	test.Run("Fail Decode", func(test *testing.T) {
		mockResponse := &http.Response{
			Body: io.NopCloser(strings.NewReader(`invalid json`)),
		}

		mockUserRepo := new(MockUserRepository)
		mockServiceRepo := new(MockServiceRepository)
		mockUserServiceRepo := new(MockUserServiceRepository)
		mockServiceService := new(MockServiceServiceRepository)

		userService := &UserServiceService{
			UserRepository:        mockUserRepo,
			ServiceRepository:     mockServiceRepo,
			UserServiceRepository: mockUserServiceRepo,
			ServiceService:        mockServiceService,
		}

		mockUserRepo.On("FindUserByEmail", "test@test.com", "basic").
			Return(user, nil)

		mockServiceRepo.On("FindServiceByName", "Asana").
			Return(service, nil)

		mockUserServiceRepo.On("FindUserServiceByServiceIdandUserId", "1", "1").
			Return(serviceOfUser, nil)

		mockServiceService.On("ExecuteApiRequest", "https://app.asana.com/api/1.0/workspaces/", "GET", "Bearer ", "accessToken", nil).
			Return(mockResponse, nil)

		_, err := userService.RetrieveAsanaUserWorkspaces("test@test.com", "basic")

		require.Error(test, err)
	})

}

func TestDecodeRequiredWorkspaceInfo(test *testing.T) {
	var user entities.User
	var service entities.Service
	var serviceOfUser entities.UserService

	user.Email = "test@test.com"
	user.Id = "1"

	service.Id = "1"

	serviceOfUser.AccessToken = "accessToken"
	serviceOfUser.RefreshToken = "refreshToken"
	serviceOfUser.ExpiryDate = time.Now().Add(time.Hour).Format(formattingDate)

	test.Run("Successful", func(test *testing.T) {
		mockResponse := &http.Response{
			Body: io.NopCloser(strings.NewReader(`
				{"data":[{"gid":"1209091057287055","name":"Cross-functional project plan","resource_type":"project"}]}
			`)),
		}

		mockUserRepo := new(MockUserRepository)
		mockServiceRepo := new(MockServiceRepository)
		mockUserServiceRepo := new(MockUserServiceRepository)
		mockServiceService := new(MockServiceServiceRepository)

		userService := &UserServiceService{
			UserRepository:        mockUserRepo,
			ServiceRepository:     mockServiceRepo,
			UserServiceRepository: mockUserServiceRepo,
			ServiceService:        mockServiceService,
		}

		mockUserRepo.On("FindUserByEmail", "test@test.com", "basic").
			Return(user, nil)

		mockServiceRepo.On("FindServiceByName", "Asana").
			Return(service, nil)

		mockUserServiceRepo.On("FindUserServiceByServiceIdandUserId", "1", "1").
			Return(serviceOfUser, nil)

		mockServiceService.On("ExecuteApiRequest", "https://app.asana.com/api/1.0/workspaces/", "GET", "Bearer ", "accessToken", nil).
			Return(mockResponse, nil)

		_, err := userService.decodeRequiredWorkspaceInfo("test@test.com", "basic", "https://app.asana.com/api/1.0/workspaces/")

		require.NoError(test, err)
	})

	test.Run("Fail Call API", func(test *testing.T) {

		mockUserRepo := new(MockUserRepository)
		mockServiceRepo := new(MockServiceRepository)
		mockUserServiceRepo := new(MockUserServiceRepository)
		mockServiceService := new(MockServiceServiceRepository)

		userService := &UserServiceService{
			UserRepository:        mockUserRepo,
			ServiceRepository:     mockServiceRepo,
			UserServiceRepository: mockUserServiceRepo,
			ServiceService:        mockServiceService,
		}

		mockUserRepo.On("FindUserByEmail", "test@test.com", "basic").
			Return(user, nil)

		mockServiceRepo.On("FindServiceByName", "Asana").
			Return(service, nil)

		mockUserServiceRepo.On("FindUserServiceByServiceIdandUserId", "1", "1").
			Return(serviceOfUser, errors.New("Fail API Refresh"))

		_, err := userService.decodeRequiredWorkspaceInfo("test@test.com", "basic", "https://app.asana.com/api/1.0/workspaces/")

		require.EqualError(test, err, "Fail API Refresh")
	})

	test.Run("Fail Execute API", func(test *testing.T) {
		mockResponse := &http.Response{
			Body: io.NopCloser(strings.NewReader(`
				{"data":[{"gid":"1209091057287055","name":"Cross-functional project plan","resource_type":"project"}]}
			`)),
		}

		mockUserRepo := new(MockUserRepository)
		mockServiceRepo := new(MockServiceRepository)
		mockUserServiceRepo := new(MockUserServiceRepository)
		mockServiceService := new(MockServiceServiceRepository)

		userService := &UserServiceService{
			UserRepository:        mockUserRepo,
			ServiceRepository:     mockServiceRepo,
			UserServiceRepository: mockUserServiceRepo,
			ServiceService:        mockServiceService,
		}

		mockUserRepo.On("FindUserByEmail", "test@test.com", "basic").
			Return(user, nil)

		mockServiceRepo.On("FindServiceByName", "Asana").
			Return(service, nil)

		mockUserServiceRepo.On("FindUserServiceByServiceIdandUserId", "1", "1").
			Return(serviceOfUser, nil)

		mockServiceService.On("ExecuteApiRequest", "https://app.asana.com/api/1.0/workspaces/", "GET", "Bearer ", "accessToken", nil).
			Return(mockResponse, errors.New("Fail Execute API"))

		_, err := userService.decodeRequiredWorkspaceInfo("test@test.com", "basic", "https://app.asana.com/api/1.0/workspaces/")

		require.EqualError(test, err, "Fail Execute API")
	})

	test.Run("Fail Decode", func(test *testing.T) {
		mockResponse := &http.Response{
			Body: io.NopCloser(strings.NewReader(`invalid json`)),
		}

		mockUserRepo := new(MockUserRepository)
		mockServiceRepo := new(MockServiceRepository)
		mockUserServiceRepo := new(MockUserServiceRepository)
		mockServiceService := new(MockServiceServiceRepository)

		userService := &UserServiceService{
			UserRepository:        mockUserRepo,
			ServiceRepository:     mockServiceRepo,
			UserServiceRepository: mockUserServiceRepo,
			ServiceService:        mockServiceService,
		}

		mockUserRepo.On("FindUserByEmail", "test@test.com", "basic").
			Return(user, nil)

		mockServiceRepo.On("FindServiceByName", "Asana").
			Return(service, nil)

		mockUserServiceRepo.On("FindUserServiceByServiceIdandUserId", "1", "1").
			Return(serviceOfUser, nil)

		mockServiceService.On("ExecuteApiRequest", "https://app.asana.com/api/1.0/workspaces/", "GET", "Bearer ", "accessToken", nil).
			Return(mockResponse, errors.New("Fail Execute API"))

		_, err := userService.decodeRequiredWorkspaceInfo("test@test.com", "basic", "https://app.asana.com/api/1.0/workspaces/")

		require.Error(test, err)
	})
}

func TestRetrieveAsanaWorkspaceAssignees(test *testing.T) {
	var user entities.User
	var service entities.Service
	var serviceOfUser entities.UserService

	user.Email = "test@test.com"
	user.Id = "1"

	service.Id = "1"

	serviceOfUser.AccessToken = "accessToken"
	serviceOfUser.RefreshToken = "refreshToken"
	serviceOfUser.ExpiryDate = time.Now().Add(time.Hour).Format(formattingDate)

	test.Run("Successful", func(test *testing.T) {
		mockResponse := &http.Response{
			Body: io.NopCloser(strings.NewReader(`
				{"data":[{"gid":"1209091057287055","name":"Cross-functional project plan","resource_type":"project"}]}
			`)),
		}

		mockUserRepo := new(MockUserRepository)
		mockServiceRepo := new(MockServiceRepository)
		mockUserServiceRepo := new(MockUserServiceRepository)
		mockServiceService := new(MockServiceServiceRepository)

		userService := &UserServiceService{
			UserRepository:        mockUserRepo,
			ServiceRepository:     mockServiceRepo,
			UserServiceRepository: mockUserServiceRepo,
			ServiceService:        mockServiceService,
		}

		mockUserRepo.On("FindUserByEmail", "test@test.com", "basic").
			Return(user, nil)

		mockServiceRepo.On("FindServiceByName", "Asana").
			Return(service, nil)

		mockUserServiceRepo.On("FindUserServiceByServiceIdandUserId", "1", "1").
			Return(serviceOfUser, nil)

		mockServiceService.On("ExecuteApiRequest", "https://app.asana.com/api/1.0/workspaces/1/users", "GET", "Bearer ", "accessToken", nil).
			Return(mockResponse, nil)

		_, err := userService.RetrieveAsanaWorkspaceAssignees("test@test.com", "basic", "1")

		require.NoError(test, err)
	})
}

func TestRetrieveAsanaWorkspaceProjects(test *testing.T) {
	var user entities.User
	var service entities.Service
	var serviceOfUser entities.UserService

	user.Email = "test@test.com"
	user.Id = "1"

	service.Id = "1"

	serviceOfUser.AccessToken = "accessToken"
	serviceOfUser.RefreshToken = "refreshToken"
	serviceOfUser.ExpiryDate = time.Now().Add(time.Hour).Format(formattingDate)

	test.Run("Successful", func(test *testing.T) {
		mockResponse := &http.Response{
			Body: io.NopCloser(strings.NewReader(`
				{"data":[{"gid":"1209091057287055","name":"Cross-functional project plan","resource_type":"project"}]}
			`)),
		}

		mockUserRepo := new(MockUserRepository)
		mockServiceRepo := new(MockServiceRepository)
		mockUserServiceRepo := new(MockUserServiceRepository)
		mockServiceService := new(MockServiceServiceRepository)

		userService := &UserServiceService{
			UserRepository:        mockUserRepo,
			ServiceRepository:     mockServiceRepo,
			UserServiceRepository: mockUserServiceRepo,
			ServiceService:        mockServiceService,
		}

		mockUserRepo.On("FindUserByEmail", "test@test.com", "basic").
			Return(user, nil)

		mockServiceRepo.On("FindServiceByName", "Asana").
			Return(service, nil)

		mockUserServiceRepo.On("FindUserServiceByServiceIdandUserId", "1", "1").
			Return(serviceOfUser, nil)

		mockServiceService.On("ExecuteApiRequest", "https://app.asana.com/api/1.0/workspaces/1/projects", "GET", "Bearer ", "accessToken", nil).
			Return(mockResponse, nil)

		_, err := userService.RetrieveAsanaWorkspaceProjects("test@test.com", "basic", "1")

		require.NoError(test, err)
	})
}

func TestRetrieveAsanaWorkspaceTags(test *testing.T) {
	var user entities.User
	var service entities.Service
	var serviceOfUser entities.UserService

	user.Email = "test@test.com"
	user.Id = "1"

	service.Id = "1"

	serviceOfUser.AccessToken = "accessToken"
	serviceOfUser.RefreshToken = "refreshToken"
	serviceOfUser.ExpiryDate = time.Now().Add(time.Hour).Format(formattingDate)

	test.Run("Successful", func(test *testing.T) {
		mockResponse := &http.Response{
			Body: io.NopCloser(strings.NewReader(`
				{"data":[{"gid":"1209091057287055","name":"Cross-functional project plan","resource_type":"project"}]}
			`)),
		}

		mockUserRepo := new(MockUserRepository)
		mockServiceRepo := new(MockServiceRepository)
		mockUserServiceRepo := new(MockUserServiceRepository)
		mockServiceService := new(MockServiceServiceRepository)

		userService := &UserServiceService{
			UserRepository:        mockUserRepo,
			ServiceRepository:     mockServiceRepo,
			UserServiceRepository: mockUserServiceRepo,
			ServiceService:        mockServiceService,
		}

		mockUserRepo.On("FindUserByEmail", "test@test.com", "basic").
			Return(user, nil)

		mockServiceRepo.On("FindServiceByName", "Asana").
			Return(service, nil)

		mockUserServiceRepo.On("FindUserServiceByServiceIdandUserId", "1", "1").
			Return(serviceOfUser, nil)

		mockServiceService.On("ExecuteApiRequest", "https://app.asana.com/api/1.0/workspaces/1/tags", "GET", "Bearer ", "accessToken", nil).
			Return(mockResponse, nil)

		_, err := userService.RetrieveAsanaWorkspaceTags("test@test.com", "basic", "1")

		require.NoError(test, err)
	})
}
