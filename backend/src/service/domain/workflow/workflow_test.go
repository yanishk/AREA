package workflow_service

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"backend/src/entities"
)

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

type MockUserServiceRepository struct {
	mock.Mock
}

func (m *MockUserServiceRepository) CallApiAndRefresh(email, connectionType, serviceName string) (string, error) {
	args := m.Called(email, connectionType, serviceName)
	return args.String(0), args.Error(1)
}

func (m *MockUserServiceRepository) UpdateTokenForService(code, serviceName, appType, email, connectionType string) error {
	args := m.Called(code, serviceName, appType, email, connectionType)
	return args.Error(0)
}

func (m *MockUserServiceRepository) RetrieveGithubUserRepositories(email, connectionType string) ([]entities.GithubRepository, error) {
	args := m.Called(email, connectionType)
	return args.Get(0).([]entities.GithubRepository), args.Error(1)
}

func (m *MockUserServiceRepository) RetrieveGitlabUserProjects(email, connectionType string) ([]entities.GitlabProject, error) {
	args := m.Called(email, connectionType)
	return args.Get(0).([]entities.GitlabProject), args.Error(1)
}

func (m *MockUserServiceRepository) RetrieveDiscordUserServers(email, connectionType string) ([]map[string]interface{}, error) {
	args := m.Called(email, connectionType)
	return args.Get(0).([]map[string]interface{}), args.Error(1)
}

func (m *MockUserServiceRepository) RetrieveAsanaUserWorkspaces(email, connectionType string) (entities.AsanaWorkspacesInfo, error) {
	args := m.Called(email, connectionType)
	return args.Get(0).(entities.AsanaWorkspacesInfo), args.Error(1)
}

func (m *MockUserServiceRepository) RetrieveAsanaWorkspaceAssignees(email, connectionType, workspaceId string) (entities.AsanaWorkspacesInfo, error) {
	args := m.Called(email, connectionType, workspaceId)
	return args.Get(0).(entities.AsanaWorkspacesInfo), args.Error(1)
}

func (m *MockUserServiceRepository) RetrieveAsanaWorkspaceProjects(email, connectionType, workspaceId string) (entities.AsanaWorkspacesInfo, error) {
	args := m.Called(email, connectionType, workspaceId)
	return args.Get(0).(entities.AsanaWorkspacesInfo), args.Error(1)
}

func (m *MockUserServiceRepository) RetrieveAsanaWorkspaceTags(email, connectionType, workspaceId string) (entities.AsanaWorkspacesInfo, error) {
	args := m.Called(email, connectionType, workspaceId)
	return args.Get(0).(entities.AsanaWorkspacesInfo), args.Error(1)
}

func (m *MockUserServiceRepository) RetrieveUserServiceAuthenticationStatus(email, connectionType, serviceName string) (bool, error) {
	args := m.Called(email, connectionType, serviceName)
	return args.Bool(0), args.Error(1)
}

func TestCreateWorkflow(test *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockWorkflowRepo := new(MockWorkflowRepository)
	service := &WorkflowService{
		UserRepository:     mockUserRepo,
		WorkflowRepository: mockWorkflowRepo,
	}

	test.Run("User not found", func(test *testing.T) {
		mockUserRepo.On("FindUserByEmail", "test@test.com", "basic").
			Return(entities.User{}, errors.New("user not found")).Once()

		err := service.CreateWorkflow("test@test.com", "basic", entities.NewWorkflow{})
		require.EqualError(test, err, "user not found")
	})

	test.Run("Fail workflow creation", func(test *testing.T) {
		var user entities.User
		user.Id = "1"

		mockUserRepo.On("FindUserByEmail", "test@test.com", "basic").
			Return(user, nil).Once()

		mockWorkflowRepo.On("CreateWorkflow", "Test Workflow", "1", "1", "2", map[string]interface{}{"key": "value"}, map[string]interface{}{"key": "value"}, map[string]interface{}{"key": "value"}).
			Return(errors.New("Fail workflow creation")).Once()

		newWorkflow := entities.NewWorkflow{
			Name:          "Test Workflow",
			ActionId:      "1",
			ReactionId:    "2",
			ActionParam:   map[string]interface{}{"key": "value"},
			ReactionParam: map[string]interface{}{"key": "value"},
			ActionData:    map[string]interface{}{"key": "value"},
		}

		err := service.CreateWorkflow("test@test.com", "basic", newWorkflow)
		require.EqualError(test, err, "Fail workflow creation")
	})

	test.Run("Successful", func(test *testing.T) {
		var user entities.User
		user.Id = "1"

		mockUserRepo.On("FindUserByEmail", "test@test.com", "basic").
			Return(user, nil).Once()

		mockWorkflowRepo.On("CreateWorkflow", "Test Workflow", "1", "1", "2", map[string]interface{}{"key": "value"}, map[string]interface{}{"key": "value"}, map[string]interface{}{"key": "value"}).
			Return(nil).Once()

		newWorkflow := entities.NewWorkflow{
			Name:          "Test Workflow",
			ActionId:      "1",
			ReactionId:    "2",
			ActionParam:   map[string]interface{}{"key": "value"},
			ReactionParam: map[string]interface{}{"key": "value"},
			ActionData:    map[string]interface{}{"key": "value"},
		}

		err := service.CreateWorkflow("test@test.com", "basic", newWorkflow)
		require.NoError(test, err)
	})
}

func TestGetUserWorkflows(test *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockWorkflowRepo := new(MockWorkflowRepository)
	service := &WorkflowService{
		UserRepository:     mockUserRepo,
		WorkflowRepository: mockWorkflowRepo,
	}

	test.Run("User not found", func(test *testing.T) {
		mockUserRepo.On("FindUserByEmail", "test@test.com", "basic").
			Return(entities.User{}, errors.New("user not found")).Once()

		_, err := service.GetUserWorkflows("test@test.com", "basic")
		require.EqualError(test, err, "user not found")
	})

	test.Run("Fail retrieve workflow", func(test *testing.T) {
		var user entities.User
		user.Id = "1"

		mockUserRepo.On("FindUserByEmail", "test@test.com", "basic").
			Return(user, nil).Once()

		mockWorkflowRepo.On("FindWorkflowsByOwnerId", "1").
			Return([]entities.Workflow{}, errors.New("Fail retrieve workflow")).Once()

		_, err := service.GetUserWorkflows("test@test.com", "basic")
		require.EqualError(test, err, "Fail retrieve workflow")
	})

	test.Run("Successful", func(test *testing.T) {
		var user entities.User
		user.Id = "1"

		mockUserRepo.On("FindUserByEmail", "test@test.com", "basic").
			Return(user, nil).Once()

		mockWorkflowRepo.On("FindWorkflowsByOwnerId", "1").
			Return([]entities.Workflow{}, nil).Once()

		_, err := service.GetUserWorkflows("test@test.com", "basic")
		require.NoError(test, err)
	})
}

func TestUpdateWorkflow(test *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockWorkflowRepo := new(MockWorkflowRepository)
	service := &WorkflowService{
		UserRepository:     mockUserRepo,
		WorkflowRepository: mockWorkflowRepo,
	}

	test.Run("User not found", func(test *testing.T) {
		var updateWorkflow entities.UpdatedWorkflow

		mockWorkflowRepo.On("FindWorkflowById", "1").
			Return(entities.Workflow{}, errors.New("user not found")).Once()

		err := service.UpdateWorkflow("1", updateWorkflow)
		require.EqualError(test, err, "user not found")
	})

	test.Run("Fail update workflow", func(test *testing.T) {
		var updateWorkflow entities.UpdatedWorkflow

		mockWorkflowRepo.On("FindWorkflowById", "1").
			Return(entities.Workflow{}, nil).Once()

		mockWorkflowRepo.On("UpdateWorkflow", "1", entities.Workflow{}).
			Return(errors.New("Fail update workflow")).Once()

		err := service.UpdateWorkflow("1", updateWorkflow)
		require.EqualError(test, err, "Fail update workflow")
	})

	test.Run("Successful", func(test *testing.T) {
		var updateWorkflow entities.UpdatedWorkflow

		mockWorkflowRepo.On("FindWorkflowById", "1").
			Return(entities.Workflow{}, nil).Once()

		mockWorkflowRepo.On("UpdateWorkflow", "1", entities.Workflow{}).
			Return(nil).Once()

		err := service.UpdateWorkflow("1", updateWorkflow)
		require.NoError(test, err)
	})
}

func TestDeleteWorkflow(test *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockWorkflowRepo := new(MockWorkflowRepository)
	service := &WorkflowService{
		UserRepository:     mockUserRepo,
		WorkflowRepository: mockWorkflowRepo,
	}

	test.Run("User not found", func(test *testing.T) {
		mockUserRepo.On("FindUserByEmail", "test@test.com", "basic").
			Return(entities.User{}, errors.New("user not found")).Once()

		err := service.DeleteWorkflow("test@test.com", "basic", "1")
		require.EqualError(test, err, "user not found")
	})

	test.Run("Fail delete workflow", func(test *testing.T) {
		var user entities.User
		user.Id = "1"

		mockUserRepo.On("FindUserByEmail", "test@test.com", "basic").
			Return(user, nil).Once()

		mockWorkflowRepo.On("DeleteWorkflow", "1", user.Id).
			Return(errors.New("Fail delete workflow")).Once()

		err := service.DeleteWorkflow("test@test.com", "basic", "1")
		require.EqualError(test, err, "Fail delete workflow")
	})

	test.Run("Successful", func(test *testing.T) {
		var user entities.User
		user.Id = "1"

		mockUserRepo.On("FindUserByEmail", "test@test.com", "basic").
			Return(user, nil).Once()

		mockWorkflowRepo.On("DeleteWorkflow", "1", user.Id).
			Return(nil).Once()

		err := service.DeleteWorkflow("test@test.com", "basic", "1")
		require.NoError(test, err)
	})
}

func TestGetAccessToken(test *testing.T) {
	mockUserRepo := new(MockUserRepository)
	MockUserServiceRepo := new(MockUserServiceRepository)
	service := &WorkflowService{
		UserRepository:     mockUserRepo,
		UserServiceService: MockUserServiceRepo,
	}

	test.Run("User not found", func(test *testing.T) {
		var workflow entities.Workflow

		workflow.OwnerId = "1"

		mockUserRepo.On("FindUserById", "1").
			Return(entities.User{}, errors.New("user not found")).Once()

		_, err := service.getAccessToken("service", workflow)
		require.EqualError(test, err, "Error finding user")
	})

	test.Run("Fail refresh", func(test *testing.T) {
		var workflow entities.Workflow
		var user entities.User

		workflow.OwnerId = "1"
		user.Email = "test@test.com"
		user.ConnectionType = "basic"

		mockUserRepo.On("FindUserById", "1").
			Return(user, nil).Once()

		MockUserServiceRepo.On("CallApiAndRefresh", "test@test.com", "basic", "service").
			Return("accessToken", errors.New("Refresh Failed"))

		_, err := service.getAccessToken("service", workflow)
		require.EqualError(test, err, "Refresh Failed")
	})

	test.Run("Successful", func(test *testing.T) {
		mockUserRepo := new(MockUserRepository)
		MockUserServiceRepo := new(MockUserServiceRepository)
		service := &WorkflowService{
			UserRepository:     mockUserRepo,
			UserServiceService: MockUserServiceRepo,
		}

		var workflow entities.Workflow
		var user entities.User

		workflow.OwnerId = "1"
		user.Email = "test@test.com"
		user.ConnectionType = "basic"

		mockUserRepo.On("FindUserById", "1").
			Return(user, nil).Once()

		MockUserServiceRepo.On("CallApiAndRefresh", "test@test.com", "basic", "service").
			Return("accessToken", nil)

		_, err := service.getAccessToken("service", workflow)
		require.NoError(test, err)
	})
}

func TestConcatanateSlashToPath(test *testing.T) {
	name := "test.go"
	filepath := "path"
	expected := "/" + filepath + "/" + name

	res := concatanateSlashToPath(filepath, name)

	require.Equal(test, expected, res)
}

func TestGetWorkflowStringActionParam(test *testing.T) {
	test.Run("Successful", func(test *testing.T) {
		var workflow entities.Workflow
		workflow.ActionParam = map[string]interface{}{"key": "value"}

		_, err := getWorkflowStringActionParam(workflow, "key")
		require.NoError(test, err)
	})

	test.Run("ParamKey does not exist", func(test *testing.T) {
		var workflow entities.Workflow
		workflow.ActionParam = map[string]interface{}{"key": "value"}

		_, err := getWorkflowStringActionParam(workflow, "false")
		require.EqualError(test, err, "Missing required field")
	})

	test.Run("ParamKey does not exist", func(test *testing.T) {
		var workflow entities.Workflow
		workflow.ActionParam = map[string]interface{}{"key": 1}

		_, err := getWorkflowStringActionParam(workflow, "key")
		require.EqualError(test, err, "Missing required field")
	})
}
