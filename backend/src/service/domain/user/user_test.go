package user_service

import (
	"errors"
	"io"
	"net/http"
	"testing"

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

type MockWorkflowRepository struct {
	mock.Mock
}

func (m *MockWorkflowRepository) CreateWorkflow(name, ownerId, actionId, reactionId string, actionParam, reactionParam, actionData map[string]interface{}) error {
	args := m.Called(name, ownerId, actionId, reactionId, actionParam, reactionParam, actionData)
	return args.Error(0)
}

func (m *MockWorkflowRepository) FindWorkflowById(id string) (entities.Workflow, error) {
	args := m.Called(id)
	return args.Get(0).(entities.Workflow), args.Error(1)
}

func (m *MockWorkflowRepository) FindWorkflowsByActionId(actionId string) ([]entities.Workflow, error) {
	args := m.Called(actionId)
	return args.Get(0).([]entities.Workflow), args.Error(1)
}

func (m *MockWorkflowRepository) FindWorkflowsByOwnerId(ownerId string) ([]entities.Workflow, error) {
	args := m.Called(ownerId)
	return args.Get(0).([]entities.Workflow), args.Error(1)
}

func (m *MockWorkflowRepository) UpdateWorkflow(id string, updatedWorkflow entities.Workflow) error {
	args := m.Called(id, updatedWorkflow)
	return args.Error(0)
}

func (m *MockWorkflowRepository) DeleteWorkflow(id, ownerId string) error {
	args := m.Called(id, ownerId)
	return args.Error(0)
}

func (m *MockWorkflowRepository) DeleteWorkflowByOwnerId(ownerId string) error {
	args := m.Called(ownerId)
	return args.Error(0)
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

func TestCreateUser(test *testing.T) {
	test.Run("Successful No Basic", func(test *testing.T) {
		mockUserRepo := new(MockUserRepository)
		mockServiceRepo := new(MockServiceRepository)

		userService := &UserService{
			UserRepository:    mockUserRepo,
			ServiceRepository: mockServiceRepo,
		}

		mockServiceRepo.On("FindServiceByName", "Spotify").
			Return(entities.Service{}, nil)

		mockUserRepo.On("CreateUser", "test@test.com", "", "Spotify").
			Return(nil)

		err := userService.CreateUser("test@test.com", "test", "Spotify")

		require.NoError(test, err)
	})

	test.Run("Fail Find Service", func(test *testing.T) {
		mockUserRepo := new(MockUserRepository)
		mockServiceRepo := new(MockServiceRepository)

		userService := &UserService{
			UserRepository:    mockUserRepo,
			ServiceRepository: mockServiceRepo,
		}

		mockServiceRepo.On("FindServiceByName", "Spotify").
			Return(entities.Service{}, errors.New("Fail Find Service"))

		err := userService.CreateUser("test@test.com", "test", "Spotify")

		require.EqualError(test, err, "Connection type doesn't exist")
	})

	test.Run("Fail Create User", func(test *testing.T) {
		mockUserRepo := new(MockUserRepository)
		mockServiceRepo := new(MockServiceRepository)

		userService := &UserService{
			UserRepository:    mockUserRepo,
			ServiceRepository: mockServiceRepo,
		}

		mockServiceRepo.On("FindServiceByName", "Spotify").
			Return(entities.Service{}, nil)

		mockUserRepo.On("CreateUser", "test@test.com", "", "Spotify").
			Return(errors.New("Fail Create User"))

		err := userService.CreateUser("test@test.com", "test", "Spotify")

		require.EqualError(test, err, "Email address already used")
	})
}

func TestCreateToken(test *testing.T) {
	_, err := createToken("email", "basic")

	require.NoError(test, err)
}

func TestLoginAuthentication(test *testing.T) {
	test.Run("Successful", func(test *testing.T) {
		var user entities.User

		mockUserRepo := new(MockUserRepository)

		userService := &UserService{
			UserRepository: mockUserRepo,
		}

		user.Email = "test@test.com"
		user.Password = "$2a$12$tlM/vFPpczFORp7v.jrJfuZ9sz0/hAuADl86YDdohIDujKwCSq08y"
		user.ConnectionType = "basic"

		mockUserRepo.On("FindUserByEmail", "test@test.com", "basic").
			Return(user, nil)

		_, err := userService.LoginAuthentication("test@test.com", "test", "basic")

		require.NoError(test, err)
	})

	test.Run("Fail Connection type", func(test *testing.T) {
		var user entities.User

		mockUserRepo := new(MockUserRepository)

		userService := &UserService{
			UserRepository: mockUserRepo,
		}

		user.Email = "test@test.com"
		user.Password = "$2a$12$tlM/vFPpczFORp7v.jrJfuZ9sz0/hAuADl86YDdohIDujKwCSq08y"
		user.ConnectionType = "basic"

		mockUserRepo.On("FindUserByEmail", "test@test.com", "basic").
			Return(user, nil)

		_, err := userService.LoginAuthentication("test@test.com", "test", "fail")

		require.EqualError(test, err, "Wrong password")
	})

	test.Run("Fail Find User", func(test *testing.T) {
		var user entities.User

		mockUserRepo := new(MockUserRepository)

		userService := &UserService{
			UserRepository: mockUserRepo,
		}

		user.Email = "test@test.com"
		user.Password = "$2a$12$tlM/vFPpczFORp7v.jrJfuZ9sz0/hAuADl86YDdohIDujKwCSq08y"
		user.ConnectionType = "basic"

		mockUserRepo.On("FindUserByEmail", "test@test.com", "basic").
			Return(user, errors.New("Fail Find User"))

		_, err := userService.LoginAuthentication("test@test.com", "test", "basic")

		require.EqualError(test, err, "Could not find requested user")
	})

	test.Run("Invalid Password", func(test *testing.T) {
		var user entities.User

		mockUserRepo := new(MockUserRepository)

		userService := &UserService{
			UserRepository: mockUserRepo,
		}

		user.Email = "test@test.com"
		user.Password = "$2a$12$tlM/vFPpczFORp7v.jrJfuZ9sz0/hAuADl86YDdohIDujKwCSq08y"
		user.ConnectionType = "basic"

		mockUserRepo.On("FindUserByEmail", "test@test.com", "basic").
			Return(user, nil)

		_, err := userService.LoginAuthentication("test@test.com", "fail", "basic")

		require.EqualError(test, err, "Wrong password")
	})
}

func TestLoginWithService(test *testing.T) {
	test.Run("Successful", func(test *testing.T) {
		var userInfo entities.UserInfo
		var user entities.User
		var resultToken entities.ResultToken

		mockUserRepo := new(MockUserRepository)
		mockServiceRepo := new(MockServiceRepository)
		mockServiceServiceRepo := new(MockServiceServiceRepository)

		userService := &UserService{
			UserRepository:    mockUserRepo,
			ServiceRepository: mockServiceRepo,
			ServiceService:    mockServiceServiceRepo,
		}

		userInfo.Email = "test@test.com"
		user.Email = "test@test.com"

		resultToken.AccessToken = "accessToken"

		mockServiceRepo.On("FindServiceByName", "service").
			Return(entities.Service{}, nil)

		mockServiceServiceRepo.On("GetResultTokenFromCode", "code", "service", "login", "web").
			Return(resultToken, nil)

		mockServiceServiceRepo.On("GetUserInfoFromService", "accessToken", "service").
			Return(userInfo, nil)

		mockUserRepo.On("FindUserByEmail", "test@test.com", "service").
			Return(user, nil)

		mockServiceRepo.On("FindServiceByName", "service").
			Return(entities.Service{}, nil)

		mockUserRepo.On("CreateUser", "test@test.com", "", "service").
			Return(nil)

		mockUserRepo.On("FindUserByEmail", "test@test.com", "service").
			Return(user, nil)

		_, err := userService.LoginWithService("code", "service", "web")

		require.NoError(test, err)
	})

	test.Run("Fail Find Service", func(test *testing.T) {
		mockServiceRepo := new(MockServiceRepository)

		userService := &UserService{
			ServiceRepository: mockServiceRepo,
		}

		mockServiceRepo.On("FindServiceByName", "service").
			Return(entities.Service{}, errors.New("Fail Find Service"))

		_, err := userService.LoginWithService("code", "service", "web")

		require.EqualError(test, err, "Fail Find Service")
	})

	test.Run("Fail Get Token From Code", func(test *testing.T) {
		var resultToken entities.ResultToken

		mockServiceRepo := new(MockServiceRepository)
		mockServiceServiceRepo := new(MockServiceServiceRepository)

		userService := &UserService{
			ServiceRepository: mockServiceRepo,
			ServiceService:    mockServiceServiceRepo,
		}

		resultToken.AccessToken = "accessToken"

		mockServiceRepo.On("FindServiceByName", "service").
			Return(entities.Service{}, nil)

		mockServiceServiceRepo.On("GetResultTokenFromCode", "code", "service", "login", "web").
			Return(resultToken, errors.New("Fail get token from code"))

		_, err := userService.LoginWithService("code", "service", "web")

		require.EqualError(test, err, "Fail get token from code")
	})

	test.Run("Fail Get User Info", func(test *testing.T) {
		var userInfo entities.UserInfo
		var resultToken entities.ResultToken

		mockServiceRepo := new(MockServiceRepository)
		mockServiceServiceRepo := new(MockServiceServiceRepository)

		userService := &UserService{
			ServiceRepository: mockServiceRepo,
			ServiceService:    mockServiceServiceRepo,
		}

		userInfo.Email = "test@test.com"

		resultToken.AccessToken = "accessToken"

		mockServiceRepo.On("FindServiceByName", "service").
			Return(entities.Service{}, nil)

		mockServiceServiceRepo.On("GetResultTokenFromCode", "code", "service", "login", "web").
			Return(resultToken, nil)

		mockServiceServiceRepo.On("GetUserInfoFromService", "accessToken", "service").
			Return(userInfo, errors.New("Fail get user info"))

		_, err := userService.LoginWithService("code", "service", "web")

		require.EqualError(test, err, "Fail get user info")
	})

	test.Run("Email address already used", func(test *testing.T) {
		var userInfo entities.UserInfo
		var user entities.User
		var resultToken entities.ResultToken

		mockUserRepo := new(MockUserRepository)
		mockServiceRepo := new(MockServiceRepository)
		mockServiceServiceRepo := new(MockServiceServiceRepository)

		userService := &UserService{
			UserRepository:    mockUserRepo,
			ServiceRepository: mockServiceRepo,
			ServiceService:    mockServiceServiceRepo,
		}

		userInfo.Email = "test@test.com"
		user.Email = "test@test.com"

		resultToken.AccessToken = "accessToken"

		mockServiceRepo.On("FindServiceByName", "service").
			Return(entities.Service{}, nil)

		mockServiceServiceRepo.On("GetResultTokenFromCode", "code", "service", "login", "web").
			Return(resultToken, nil)

		mockServiceServiceRepo.On("GetUserInfoFromService", "accessToken", "service").
			Return(userInfo, nil)

		mockUserRepo.On("FindUserByEmail", "test@test.com", "service").
			Return(user, errors.New("Fail find user"))

		mockServiceRepo.On("FindServiceByName", "service").
			Return(entities.Service{}, nil)

		mockUserRepo.On("CreateUser", "test@test.com", "", "service").
			Return(errors.New("Fail create user"))

		_, err := userService.LoginWithService("code", "service", "web")

		require.EqualError(test, err, "Email address already used")
	})
}

func TestGetUser(test *testing.T) {
	test.Run("Successful", func(test *testing.T) {
		mockUserRepo := new(MockUserRepository)

		userService := &UserService{
			UserRepository: mockUserRepo,
		}

		mockUserRepo.On("FindUserByEmail", "test@test.com", "basic").
			Return(entities.User{}, nil)

		_, err := userService.GetUser("test@test.com", "basic")

		require.NoError(test, err)
	})

	test.Run("Failure", func(test *testing.T) {
		mockUserRepo := new(MockUserRepository)

		userService := &UserService{
			UserRepository: mockUserRepo,
		}

		mockUserRepo.On("FindUserByEmail", "test@test.com", "basic").
			Return(entities.User{}, errors.New("Fail find user"))

		_, err := userService.GetUser("test@test.com", "basic")

		require.EqualError(test, err, "Fail find user")
	})
}

func TestModifyPassword(test *testing.T) {
	test.Run("Could not modify password", func(test *testing.T) {
		var modifyUser entities.UserModifyPassword

		userService := &UserService{}

		modifyUser.OldPassword = "test"
		modifyUser.Password = "new"

		err := userService.ModifyPassword("test@test.com", "service", modifyUser)

		require.EqualError(test, err, "Could not modify the password")
	})

	test.Run("Fail find user", func(test *testing.T) {
		var modifyUser entities.UserModifyPassword

		mockUserRepo := new(MockUserRepository)

		userService := &UserService{
			UserRepository: mockUserRepo,
		}

		modifyUser.OldPassword = "test"
		modifyUser.Password = "new"

		mockUserRepo.On("FindUserByEmail", "test@test.com", "basic").
			Return(entities.User{}, errors.New("Fail find user"))

		err := userService.ModifyPassword("test@test.com", "basic", modifyUser)

		require.EqualError(test, err, "Could not find requested user")
	})

	test.Run("Incorrect old password", func(test *testing.T) {
		var foundUser entities.User
		var modifyUser entities.UserModifyPassword

		mockUserRepo := new(MockUserRepository)

		userService := &UserService{
			UserRepository: mockUserRepo,
		}

		foundUser.Email = "test@test.com"
		foundUser.Password = "$2a$12$FcfBupWu3kUpoW.Y9UNI..YOWOaaY2B4m9jkcb7EjK61Y78T6gKRu"

		modifyUser.OldPassword = "fail"
		modifyUser.Password = "new"

		mockUserRepo.On("FindUserByEmail", "test@test.com", "basic").
			Return(foundUser, nil)

		err := userService.ModifyPassword("test@test.com", "basic", modifyUser)

		require.EqualError(test, err, "Old password is incorrect")
	})
}

func TestDeleteAccount(test *testing.T) {
	test.Run("Successful", func(test *testing.T) {
		var foundUser entities.User

		mockUserRepo := new(MockUserRepository)
		mockWorkflowRepo := new(MockWorkflowRepository)
		mockUserServiceRepo := new(MockUserServiceRepository)

		userService := &UserService{
			UserRepository:        mockUserRepo,
			WorkflowRepository:    mockWorkflowRepo,
			UserServiceRepository: mockUserServiceRepo,
		}

		foundUser.Email = "test@test.com"
		foundUser.Id = "1"

		mockUserRepo.On("FindUserByEmail", "test@test.com", "basic").
			Return(foundUser, nil)

		mockWorkflowRepo.On("DeleteWorkflowByOwnerId", foundUser.Id).
			Return(nil)

		mockUserServiceRepo.On("DeleteUserServiceByUserId", foundUser.Id).
			Return(nil)

		mockUserRepo.On("DeleteUser", "test@test.com", "basic").
			Return(nil)

		err := userService.DeleteAccount("test@test.com", "basic")

		require.NoError(test, err)
	})

	test.Run("Fail find user", func(test *testing.T) {
		var foundUser entities.User

		mockUserRepo := new(MockUserRepository)

		userService := &UserService{
			UserRepository: mockUserRepo,
		}

		foundUser.Email = "test@test.com"
		foundUser.Id = "1"

		mockUserRepo.On("FindUserByEmail", "test@test.com", "basic").
			Return(foundUser, errors.New("Fail find user"))

		err := userService.DeleteAccount("test@test.com", "basic")

		require.EqualError(test, err, "Fail find user")
	})

	test.Run("Fail Delete Workflow", func(test *testing.T) {
		var foundUser entities.User

		mockUserRepo := new(MockUserRepository)
		mockWorkflowRepo := new(MockWorkflowRepository)

		userService := &UserService{
			UserRepository:     mockUserRepo,
			WorkflowRepository: mockWorkflowRepo,
		}

		foundUser.Email = "test@test.com"
		foundUser.Id = "1"

		mockUserRepo.On("FindUserByEmail", "test@test.com", "basic").
			Return(foundUser, nil)

		mockWorkflowRepo.On("DeleteWorkflowByOwnerId", foundUser.Id).
			Return(errors.New("Fail delete workflow"))

		err := userService.DeleteAccount("test@test.com", "basic")

		require.EqualError(test, err, "Fail delete workflow")
	})

	test.Run("Fail deleter user service", func(test *testing.T) {
		var foundUser entities.User

		mockUserRepo := new(MockUserRepository)
		mockWorkflowRepo := new(MockWorkflowRepository)
		mockUserServiceRepo := new(MockUserServiceRepository)

		userService := &UserService{
			UserRepository:        mockUserRepo,
			WorkflowRepository:    mockWorkflowRepo,
			UserServiceRepository: mockUserServiceRepo,
		}

		foundUser.Email = "test@test.com"
		foundUser.Id = "1"

		mockUserRepo.On("FindUserByEmail", "test@test.com", "basic").
			Return(foundUser, nil)

		mockWorkflowRepo.On("DeleteWorkflowByOwnerId", foundUser.Id).
			Return(nil)

		mockUserServiceRepo.On("DeleteUserServiceByUserId", foundUser.Id).
			Return(errors.New("Fail delete user service"))

		err := userService.DeleteAccount("test@test.com", "basic")

		require.EqualError(test, err, "Fail delete user service")
	})

	test.Run("Fail delete user", func(test *testing.T) {
		var foundUser entities.User

		mockUserRepo := new(MockUserRepository)
		mockWorkflowRepo := new(MockWorkflowRepository)
		mockUserServiceRepo := new(MockUserServiceRepository)

		userService := &UserService{
			UserRepository:        mockUserRepo,
			WorkflowRepository:    mockWorkflowRepo,
			UserServiceRepository: mockUserServiceRepo,
		}

		foundUser.Email = "test@test.com"
		foundUser.Id = "1"

		mockUserRepo.On("FindUserByEmail", "test@test.com", "basic").
			Return(foundUser, nil)

		mockWorkflowRepo.On("DeleteWorkflowByOwnerId", foundUser.Id).
			Return(nil)

		mockUserServiceRepo.On("DeleteUserServiceByUserId", foundUser.Id).
			Return(nil)

		mockUserRepo.On("DeleteUser", "test@test.com", "basic").
			Return(errors.New("Fail delete user"))

		err := userService.DeleteAccount("test@test.com", "basic")

		require.EqualError(test, err, "Could not delete account")
	})
}

func TestFindUserById(test *testing.T) {
	test.Run("Successful", func(test *testing.T) {
		mockUserRepo := new(MockUserRepository)

		userService := &UserService{
			UserRepository: mockUserRepo,
		}

		mockUserRepo.On("FindUserById", "1").
			Return(entities.User{}, nil)

		_, err := userService.FindUserById("1")

		require.NoError(test, err)
	})
}
