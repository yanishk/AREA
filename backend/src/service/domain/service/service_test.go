package service_service

import (
	"errors"
	"testing"

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
