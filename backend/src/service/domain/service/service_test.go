package service_service

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"backend/src/entities"
)

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

type MockReactionRepository struct {
	mock.Mock
}

func (m *MockReactionRepository) CreateReaction(name, description, serviceId string, nbParam int) error {
	args := m.Called(name, description, serviceId, nbParam)
	return args.Error(0)
}

func (m *MockReactionRepository) FindReactionById(id string) (entities.Reaction, error) {
	args := m.Called(id)
	return args.Get(0).(entities.Reaction), args.Error(1)
}

func (m *MockReactionRepository) FindReactionByName(name string) (entities.Reaction, error) {
	args := m.Called(name)
	return args.Get(0).(entities.Reaction), args.Error(1)
}

func (m *MockReactionRepository) FindReactionsByServiceId(serviceId string) ([]entities.Reaction, error) {
	args := m.Called(serviceId)
	return args.Get(0).([]entities.Reaction), args.Error(1)
}

func TestFindServiceByName(test *testing.T) {
	mockServiceRepo := new(MockServiceRepository)

	serviceservice := &ServiceService{
		ServiceRepository: mockServiceRepo,
	}

	mockServiceRepo.On("FindServiceByName", "service").
		Return(entities.Service{}, nil)

	_, err := serviceservice.FindServiceByName("service")

	require.NoError(test, err)
}

func TestFindServiceById(test *testing.T) {
	mockServiceRepo := new(MockServiceRepository)

	serviceservice := &ServiceService{
		ServiceRepository: mockServiceRepo,
	}

	mockServiceRepo.On("FindServiceById", "id").
		Return(entities.Service{}, nil)

	_, err := serviceservice.FindServiceById("id")

	require.NoError(test, err)
}

func TestFindServiceByActionId(test *testing.T) {
	test.Run("Successful", func(test *testing.T) {
		var action entities.Action

		mockServiceRepo := new(MockServiceRepository)
		mockActionRepo := new(MockActionRepository)

		serviceservice := &ServiceService{
			ServiceRepository: mockServiceRepo,
			ActionRepository:  mockActionRepo,
		}

		action.ServiceId = "id"

		mockActionRepo.On("FindActionById", "id").
			Return(action, nil)

		mockServiceRepo.On("FindServiceById", "id").
			Return(entities.Service{}, nil)

		_, err := serviceservice.FindServiceByActionId("id")

		require.NoError(test, err)
	})

	test.Run("Find action failure", func(test *testing.T) {
		mockServiceRepo := new(MockServiceRepository)
		mockActionRepo := new(MockActionRepository)

		serviceservice := &ServiceService{
			ServiceRepository: mockServiceRepo,
			ActionRepository:  mockActionRepo,
		}

		mockActionRepo.On("FindActionById", "id").
			Return(entities.Action{}, errors.New("Fail find action"))

		_, err := serviceservice.FindServiceByActionId("id")

		require.EqualError(test, err, "Fail find action")
	})
}

func TestFindServiceByReactionId(test *testing.T) {
	test.Run("Successful", func(test *testing.T) {
		var reaction entities.Reaction

		mockServiceRepo := new(MockServiceRepository)
		mockReactionRepo := new(MockReactionRepository)

		serviceservice := &ServiceService{
			ServiceRepository:  mockServiceRepo,
			ReactionRepository: mockReactionRepo,
		}

		reaction.ServiceId = "id"

		mockReactionRepo.On("FindReactionById", "id").
			Return(reaction, nil)

		mockServiceRepo.On("FindServiceById", "id").
			Return(entities.Service{}, nil)

		_, err := serviceservice.FindServiceByReactionId("id")

		require.NoError(test, err)
	})

	test.Run("Find action failure", func(test *testing.T) {
		mockServiceRepo := new(MockServiceRepository)
		mockReactionRepo := new(MockReactionRepository)

		serviceservice := &ServiceService{
			ServiceRepository:  mockServiceRepo,
			ReactionRepository: mockReactionRepo,
		}

		mockReactionRepo.On("FindReactionById", "id").
			Return(entities.Reaction{}, errors.New("Fail find reaction"))

		_, err := serviceservice.FindServiceByReactionId("id")

		require.EqualError(test, err, "Fail find reaction")
	})
}

func TestFindAllServices(test *testing.T) {
	mockServiceRepo := new(MockServiceRepository)

	serviceservice := &ServiceService{
		ServiceRepository: mockServiceRepo,
	}

	mockServiceRepo.On("FindAllServices").
		Return([]entities.Service{}, nil)

	_, err := serviceservice.FindAllServices()

	require.NoError(test, err)
}

func TestRetrieveActionsFromService(test *testing.T) {
	test.Run("Successful", func(test *testing.T) {
		var service entities.Service

		mockServiceRepo := new(MockServiceRepository)
		mockActionRepo := new(MockActionRepository)

		serviceservice := &ServiceService{
			ServiceRepository: mockServiceRepo,
			ActionRepository:  mockActionRepo,
		}

		service.Id = "id"

		mockServiceRepo.On("FindServiceByName", "service").
			Return(service, nil)

		mockActionRepo.On("FindActionsByServiceId", service.Id).
			Return([]entities.Action{}, nil)

		_, err := serviceservice.RetrieveActionsFromService("service")

		require.NoError(test, err)
	})

	test.Run("Fail Find Service", func(test *testing.T) {
		mockServiceRepo := new(MockServiceRepository)

		serviceservice := &ServiceService{
			ServiceRepository: mockServiceRepo,
		}

		mockServiceRepo.On("FindServiceByName", "service").
			Return(entities.Service{}, errors.New("Fail find service"))

		_, err := serviceservice.RetrieveActionsFromService("service")

		require.EqualError(test, err, "Fail find service")
	})

	test.Run("Fail Find Actions", func(test *testing.T) {
		var service entities.Service

		mockServiceRepo := new(MockServiceRepository)
		mockActionRepo := new(MockActionRepository)

		serviceservice := &ServiceService{
			ServiceRepository: mockServiceRepo,
			ActionRepository:  mockActionRepo,
		}

		service.Id = "id"

		mockServiceRepo.On("FindServiceByName", "service").
			Return(service, nil)

		mockActionRepo.On("FindActionsByServiceId", service.Id).
			Return([]entities.Action{}, errors.New("Fail find actions"))

		_, err := serviceservice.RetrieveActionsFromService("service")

		require.EqualError(test, err, "Fail find actions")
	})
}

func TestRetrieveReactionsFromService(test *testing.T) {
	test.Run("Successful", func(test *testing.T) {
		var service entities.Service

		mockServiceRepo := new(MockServiceRepository)
		mockReactionRepo := new(MockReactionRepository)

		serviceservice := &ServiceService{
			ServiceRepository:  mockServiceRepo,
			ReactionRepository: mockReactionRepo,
		}

		service.Id = "id"

		mockServiceRepo.On("FindServiceByName", "service").
			Return(service, nil)

		mockReactionRepo.On("FindReactionsByServiceId", service.Id).
			Return([]entities.Reaction{}, nil)

		_, err := serviceservice.RetrieveReactionsFromService("service")

		require.NoError(test, err)
	})

	test.Run("Fail Find Service", func(test *testing.T) {
		mockServiceRepo := new(MockServiceRepository)

		serviceservice := &ServiceService{
			ServiceRepository: mockServiceRepo,
		}

		mockServiceRepo.On("FindServiceByName", "service").
			Return(entities.Service{}, errors.New("Fail find service"))

		_, err := serviceservice.RetrieveReactionsFromService("service")

		require.EqualError(test, err, "Fail find service")
	})

	test.Run("Fail Find Actions", func(test *testing.T) {
		var service entities.Service

		mockServiceRepo := new(MockServiceRepository)
		mockReactionRepo := new(MockReactionRepository)

		serviceservice := &ServiceService{
			ServiceRepository:  mockServiceRepo,
			ReactionRepository: mockReactionRepo,
		}

		service.Id = "id"

		mockServiceRepo.On("FindServiceByName", "service").
			Return(service, nil)

		mockReactionRepo.On("FindReactionsByServiceId", service.Id).
			Return([]entities.Reaction{}, errors.New("Fail find reactions"))

		_, err := serviceservice.RetrieveReactionsFromService("service")

		require.EqualError(test, err, "Fail find reactions")
	})
}

func TestRetrieveActionsServices(test *testing.T) {
	mockServiceRepo := new(MockServiceRepository)

	serviceservice := &ServiceService{
		ServiceRepository: mockServiceRepo,
	}

	mockServiceRepo.On("FindActionsServices").
		Return([]entities.Service{}, nil)

	_, err := serviceservice.RetrieveActionsServices()

	require.NoError(test, err)
}

func TestRetrieveReactionsServices(test *testing.T) {
	mockServiceRepo := new(MockServiceRepository)

	serviceservice := &ServiceService{
		ServiceRepository: mockServiceRepo,
	}

	mockServiceRepo.On("FindReactionsServices").
		Return([]entities.Service{}, nil)

	_, err := serviceservice.RetrieveReactionsServices()

	require.NoError(test, err)
}

func TestGetCallbackAndClientId(test *testing.T) {
	test.Run("Login", func(test *testing.T) {

		serviceName := "SPOTIFY"

		callback, clientId := getCallbackAndClientId("login", "SPOTIFY", true)

		resCallback := serviceName + "_LOGIN_CALLBACK"
		resClientId := serviceName + "_LOGIN_CLIENT_ID"

		assert.Equal(test, os.Getenv(resCallback), callback)
		assert.Equal(test, os.Getenv(resClientId), clientId)
	})

	test.Run("Service", func(test *testing.T) {

		serviceName := "SPOTIFY"

		callback, clientId := getCallbackAndClientId("service", "SPOTIFY", true)

		resCallback := serviceName + "_SERVICE_CALLBACK"
		resClientId := serviceName + "_SERVICE_CLIENT_ID"

		assert.Equal(test, os.Getenv(resCallback), callback)
		assert.Equal(test, os.Getenv(resClientId), clientId)
	})
}
