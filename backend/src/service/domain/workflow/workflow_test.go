package workflow_service

import (
	"backend/src/entities"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
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

func TestConcatanateSlashToPath(t *testing.T) {
	name := "test.go"
	filepath := "path"
	expected := "/" + filepath + "/" + name

	res := concatanateSlashToPath(filepath, name)
	if res != expected {
		t.Fatalf(`Got : %s but Expected -> %s`, res, expected)
	}
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

		err := service.CreateWorkflow("test@test.com", "basic", entities.NewWorkflow{})
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
