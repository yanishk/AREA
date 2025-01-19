package workflow_service

import (
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"backend/src/entities"
)

type MockServiceServiceRepository struct {
	mock.Mock
}

func (m *MockServiceServiceRepository) OAuth2Service(serviceName, callbackType, appType string) (string, error) {
	args := m.Called(serviceName, callbackType, appType)
	return args.String(0), args.Error(1)
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

func (m *MockServiceServiceRepository) RetrieveActionsFromService(serviceName string) ([]entities.Action, error) {
	args := m.Called(serviceName)
	return args.Get(0).([]entities.Action), args.Error(1)
}

func (m *MockServiceServiceRepository) RetrieveReactionsFromService(serviceName string) ([]entities.Reaction, error) {
	args := m.Called(serviceName)
	return args.Get(0).([]entities.Reaction), args.Error(1)
}

func (m *MockServiceServiceRepository) GetGoogleRefreshTokenRequest(refreshToken string) (*http.Request, error) {
	return nil, nil
}

func (m *MockServiceServiceRepository) GetSpotifyRefreshTokenRequest(refreshToken string) (*http.Request, error) {
	return nil, nil
}

func (m *MockServiceServiceRepository) GetDiscordRefreshTokenRequest(refreshToken string) (*http.Request, error) {
	return nil, nil
}

func (m *MockServiceServiceRepository) GetRedditRefreshTokenRequest(refreshToken string) (*http.Request, error) {
	return nil, nil
}

func (m *MockServiceServiceRepository) GetAsanaRefreshTokenRequest(refreshToken string) (*http.Request, error) {
	return nil, nil
}

func (m *MockServiceServiceRepository) GetDropboxRefreshTokenRequest(refreshToken string) (*http.Request, error) {
	return nil, nil
}

func (m *MockServiceServiceRepository) GetGitlabRefreshTokenRequest(refreshToken string) (*http.Request, error) {
	return nil, nil
}

func (m *MockServiceServiceRepository) ExecuteRequest(request *http.Request) (*http.Response, error) {
	return nil, nil
}

func (m *MockServiceServiceRepository) ExecuteApiRequest(url, method, typeToken, accessToken string, body io.Reader) (*http.Response, error) {
	args := m.Called(url, method, typeToken, accessToken, body)
	return args.Get(0).(*http.Response), args.Error(1)
}

func (m *MockServiceServiceRepository) GetResultTokenFromCode(code, serviceName, callbackType, appType string) (entities.ResultToken, error) {
	var test entities.ResultToken
	return test, nil
}

func (m *MockServiceServiceRepository) GetUserInfoFromService(accessToken, serviceName string) (entities.UserInfo, error) {
	var test entities.UserInfo
	return test, nil
}

func (m *MockServiceServiceRepository) RequestToTimeApi() (entities.TimeResponse, error) {
	args := m.Called()
	return args.Get(0).(entities.TimeResponse), args.Error(1)
}

func (m *MockServiceServiceRepository) RequestGithubUserRepositories(accessToken string) ([]entities.GithubRepository, error) {
	args := m.Called(accessToken)
	return args.Get(0).([]entities.GithubRepository), args.Error(1)
}

func (m *MockServiceServiceRepository) RequestGitlabUserProjects(accessToken string) ([]entities.GitlabProject, error) {
	return nil, nil
}

func (m *MockServiceServiceRepository) RetrieveDiscordGuildChannels(guildId string) ([]map[string]interface{}, error) {
	args := m.Called(guildId)
	return args.Get(0).([]map[string]interface{}), args.Error(1)
}

func (m *MockServiceServiceRepository) FindAllServices() ([]entities.Service, error) {
	args := m.Called()
	return args.Get(0).([]entities.Service), args.Error(1)
}

type MockActionRepository struct {
	mock.Mock
}

func (m *MockActionRepository) CreateAction(name, description, serviceId string, nbParam int) error {
	args := m.Called(name, description, serviceId, nbParam)
	return args.Error(0)
}

func (m *MockActionRepository) FindActionById(id string) (entities.Action, error) {
	args := m.Called(id)
	return args.Get(0).(entities.Action), args.Error(1)
}

func (m *MockActionRepository) FindActionByName(name string) (entities.Action, error) {
	args := m.Called(name)
	return args.Get(0).(entities.Action), args.Error(1)
}

func (m *MockActionRepository) FindActionsByServiceId(serviceId string) ([]entities.Action, error) {
	args := m.Called(serviceId)
	return args.Get(0).([]entities.Action), args.Error(1)
}

func (m *MockActionRepository) FindActionByNameAndServiceId(name, serviceId string) (entities.Action, error) {
	args := m.Called(name, serviceId)
	return args.Get(0).(entities.Action), args.Error(1)
}

func TestCheckWorkflowsWithWeatherActions(test *testing.T) {
	action := entities.Action{
		Id:          "1",
		Name:        "action",
		Description: "description",
		ServiceId:   "2",
		NbParam:     0,
	}

	workflowEntities := []entities.Workflow{
		{Name: "test 1", IsActivated: true, ActionId: "1"},
		{Name: "test 2", IsActivated: false, ActionId: "2"},
	}
	test.Run("Success", func(test *testing.T) {
		mockWorkflowRepo := new(MockWorkflowRepository)

		workflow := &WorkflowService{
			WorkflowRepository: mockWorkflowRepo,
		}

		mockWorkflowRepo.On("FindWorkflowsByActionId", action.Id).
			Return(workflowEntities, nil)

		err := workflow.checkWorkflowsWithWeatherActions(action)

		require.NoError(test, err)
	})

	test.Run("Fail find workflow", func(test *testing.T) {
		mockWorkflowRepo := new(MockWorkflowRepository)

		workflow := &WorkflowService{
			WorkflowRepository: mockWorkflowRepo,
		}

		mockWorkflowRepo.On("FindWorkflowsByActionId", action.Id).
			Return(workflowEntities, errors.New("Fail find workflows"))

		err := workflow.checkWorkflowsWithWeatherActions(action)

		require.EqualError(test, err, "Fail find workflows")
	})
}

func TestCheckWeatherActions(test *testing.T) {
	test.Run("Success", func(test *testing.T) {
		serviceFound := entities.Service{
			Id: "1",
		}

		actionFound := []entities.Action{
			{Name: "Action 1", Id: "1"},
		}

		workflowFound := []entities.Workflow{
			{Name: "test 1", IsActivated: true, ActionId: "1"},
		}

		mockServiceServiecRepo := new(MockServiceServiceRepository)
		mockActionRepo := new(MockActionRepository)
		mockWorkflowRepo := new(MockWorkflowRepository)

		workflow := &WorkflowService{
			ServiceService:     mockServiceServiecRepo,
			ActionRepository:   mockActionRepo,
			WorkflowRepository: mockWorkflowRepo,
		}

		mockServiceServiecRepo.On("FindServiceByName", "FreeWeather").
			Return(serviceFound, nil)

		mockActionRepo.On("FindActionsByServiceId", serviceFound.Id).
			Return(actionFound, nil)

		mockWorkflowRepo.On("FindWorkflowsByActionId", actionFound[0].Id).
			Return(workflowFound, nil)

		err := workflow.CheckWeatherActions()

		require.NoError(test, err)
	})

	test.Run("Fail Find Service", func(test *testing.T) {
		mockServiceServiecRepo := new(MockServiceServiceRepository)

		workflow := &WorkflowService{
			ServiceService: mockServiceServiecRepo,
		}

		mockServiceServiecRepo.On("FindServiceByName", "FreeWeather").
			Return(entities.Service{}, errors.New("Fail find service"))

		err := workflow.CheckWeatherActions()

		require.EqualError(test, err, "Fail find service")
	})

	test.Run("Fail Find Actions", func(test *testing.T) {
		serviceFound := entities.Service{
			Id: "1",
		}

		mockServiceServiecRepo := new(MockServiceServiceRepository)
		mockActionRepo := new(MockActionRepository)

		workflow := &WorkflowService{
			ServiceService:   mockServiceServiecRepo,
			ActionRepository: mockActionRepo,
		}

		mockServiceServiecRepo.On("FindServiceByName", "FreeWeather").
			Return(serviceFound, nil)

		mockActionRepo.On("FindActionsByServiceId", serviceFound.Id).
			Return([]entities.Action{}, errors.New("Fail find actions"))

		err := workflow.CheckWeatherActions()

		require.EqualError(test, err, "Fail find actions")
	})

	test.Run("Fail Check Workflows", func(test *testing.T) {
		serviceFound := entities.Service{
			Id: "1",
		}

		actionFound := []entities.Action{
			{Name: "Action 1", Id: "1"},
		}

		workflowFound := []entities.Workflow{
			{Name: "test 1", IsActivated: true, ActionId: "1"},
		}

		mockServiceServiecRepo := new(MockServiceServiceRepository)
		mockActionRepo := new(MockActionRepository)
		mockWorkflowRepo := new(MockWorkflowRepository)

		workflow := &WorkflowService{
			ServiceService:     mockServiceServiecRepo,
			ActionRepository:   mockActionRepo,
			WorkflowRepository: mockWorkflowRepo,
		}

		mockServiceServiecRepo.On("FindServiceByName", "FreeWeather").
			Return(serviceFound, nil)

		mockActionRepo.On("FindActionsByServiceId", serviceFound.Id).
			Return(actionFound, nil)

		mockWorkflowRepo.On("FindWorkflowsByActionId", actionFound[0].Id).
			Return(workflowFound, errors.New("Fail check workflows"))

		err := workflow.CheckWeatherActions()

		require.EqualError(test, err, "Fail check workflows")
	})
}
