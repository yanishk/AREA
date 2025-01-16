package about_service

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"backend/src/entities"
)

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

func TestGetAboutAction(test *testing.T) {
	var action entities.Action

	action.Name = "test"
	action.Description = "description"

	result := getAboutAction(action)

	assert.Equal(test, action.Name, result.Name)
	assert.Equal(test, action.Description, result.Description)
}

func TestGetAboutReaction(test *testing.T) {
	var reaction entities.Reaction

	reaction.Name = "test"
	reaction.Description = "description"

	result := getAboutReaction(reaction)

	assert.Equal(test, reaction.Name, result.Name)
	assert.Equal(test, reaction.Description, result.Description)
}

func TestGetAboutActions(test *testing.T) {
	test.Run("Successful", func(test *testing.T) {
		mockActionRepo := new(MockActionRepository)

		about := &AboutService{
			ActionRepository: mockActionRepo,
		}

		mockActionRepo.On("FindActionsByServiceId", "1").
			Return([]entities.Action{}, nil)

		_, err := about.getAboutActions("1")

		require.NoError(test, err)
	})

	test.Run("Failure", func(test *testing.T) {
		mockActionRepo := new(MockActionRepository)

		about := &AboutService{
			ActionRepository: mockActionRepo,
		}

		mockActionRepo.On("FindActionsByServiceId", "1").
			Return([]entities.Action{}, errors.New("Fail find actions"))

		_, err := about.getAboutActions("1")

		require.EqualError(test, err, "Fail find actions")
	})
}

func TestGetAboutReactions(test *testing.T) {
	test.Run("Successful", func(test *testing.T) {
		mockReactionRepo := new(MockReactionRepository)

		about := &AboutService{
			ReactionRepository: mockReactionRepo,
		}

		mockReactionRepo.On("FindReactionsByServiceId", "1").
			Return([]entities.Reaction{}, nil)

		_, err := about.getAboutReactions("1")

		require.NoError(test, err)
	})

	test.Run("Failure", func(test *testing.T) {
		mockReactionRepo := new(MockReactionRepository)

		about := &AboutService{
			ReactionRepository: mockReactionRepo,
		}

		mockReactionRepo.On("FindReactionsByServiceId", "1").
			Return([]entities.Reaction{}, errors.New("Fail find reactions"))

		_, err := about.getAboutReactions("1")

		require.EqualError(test, err, "Fail find reactions")
	})
}

func TestGetAboutService(test *testing.T) {
	test.Run("Successful", func(test *testing.T) {
		var service entities.Service

		mockActionRepo := new(MockActionRepository)
		mockReactionRepo := new(MockReactionRepository)

		about := &AboutService{
			ActionRepository:   mockActionRepo,
			ReactionRepository: mockReactionRepo,
		}

		service.Id = "1"

		mockActionRepo.On("FindActionsByServiceId", "1").
			Return([]entities.Action{}, nil)

		mockReactionRepo.On("FindReactionsByServiceId", "1").
			Return([]entities.Reaction{}, nil)

		_, err := about.getAboutService(service)

		require.NoError(test, err)
	})

	test.Run("Fail action", func(test *testing.T) {
		var service entities.Service

		mockActionRepo := new(MockActionRepository)

		about := &AboutService{
			ActionRepository: mockActionRepo,
		}

		service.Id = "1"

		mockActionRepo.On("FindActionsByServiceId", "1").
			Return([]entities.Action{}, errors.New("Fail find action"))

		_, err := about.getAboutService(service)

		require.EqualError(test, err, "Fail find action")
	})

	test.Run("Fail reaction", func(test *testing.T) {
		var service entities.Service

		mockActionRepo := new(MockActionRepository)
		mockReactionRepo := new(MockReactionRepository)

		about := &AboutService{
			ActionRepository:   mockActionRepo,
			ReactionRepository: mockReactionRepo,
		}

		service.Id = "1"

		mockActionRepo.On("FindActionsByServiceId", "1").
			Return([]entities.Action{}, nil)

		mockReactionRepo.On("FindReactionsByServiceId", "1").
			Return([]entities.Reaction{}, errors.New("Fail find reaction"))

		_, err := about.getAboutService(service)

		require.EqualError(test, err, "Fail find reaction")
	})
}

func TestGetAboutServices(test *testing.T) {
	test.Run("Successful", func(test *testing.T) {
		service := []entities.Service{
			{Id: "1"},
		}

		mockActionRepo := new(MockActionRepository)
		mockReactionRepo := new(MockReactionRepository)
		mockServiceRepo := new(MockServiceRepository)

		about := &AboutService{
			ActionRepository:   mockActionRepo,
			ReactionRepository: mockReactionRepo,
			ServiceRepository:  mockServiceRepo,
		}

		mockServiceRepo.On("FindAllServices").
			Return(service, nil)

		mockActionRepo.On("FindActionsByServiceId", "1").
			Return([]entities.Action{}, nil)

		mockReactionRepo.On("FindReactionsByServiceId", "1").
			Return([]entities.Reaction{}, nil)

		_, err := about.getAboutServices()

		require.NoError(test, err)
	})

	test.Run("Fail Find Services", func(test *testing.T) {
		service := []entities.Service{
			{Id: "1"},
		}

		mockServiceRepo := new(MockServiceRepository)

		about := &AboutService{
			ServiceRepository: mockServiceRepo,
		}

		mockServiceRepo.On("FindAllServices").
			Return(service, errors.New("Fail find services"))

		_, err := about.getAboutServices()

		require.EqualError(test, err, "Fail find services")
	})

	test.Run("Fail Get About Services", func(test *testing.T) {
		service := []entities.Service{
			{Id: "1"},
		}

		mockActionRepo := new(MockActionRepository)
		mockServiceRepo := new(MockServiceRepository)

		about := &AboutService{
			ActionRepository:  mockActionRepo,
			ServiceRepository: mockServiceRepo,
		}

		mockServiceRepo.On("FindAllServices").
			Return(service, nil)

		mockActionRepo.On("FindActionsByServiceId", "1").
			Return([]entities.Action{}, errors.New("Fail get about services"))

		_, err := about.getAboutServices()

		require.EqualError(test, err, "Fail get about services")
	})
}

func TestGetAboutServer(test *testing.T) {
	test.Run("Successful", func(test *testing.T) {
		service := []entities.Service{
			{Id: "1"},
		}

		mockActionRepo := new(MockActionRepository)
		mockReactionRepo := new(MockReactionRepository)
		mockServiceRepo := new(MockServiceRepository)

		about := &AboutService{
			ActionRepository:   mockActionRepo,
			ReactionRepository: mockReactionRepo,
			ServiceRepository:  mockServiceRepo,
		}

		mockServiceRepo.On("FindAllServices").
			Return(service, nil)

		mockActionRepo.On("FindActionsByServiceId", "1").
			Return([]entities.Action{}, nil)

		mockReactionRepo.On("FindReactionsByServiceId", "1").
			Return([]entities.Reaction{}, nil)

		_, err := about.GetAboutServer(entities.About{})

		require.NoError(test, err)
	})

	test.Run("Failure", func(test *testing.T) {
		service := []entities.Service{
			{Id: "1"},
		}

		mockServiceRepo := new(MockServiceRepository)

		about := &AboutService{
			ServiceRepository: mockServiceRepo,
		}

		mockServiceRepo.On("FindAllServices").
			Return(service, errors.New("Fail find service"))

		_, err := about.GetAboutServer(entities.About{})

		require.EqualError(test, err, "Fail find service")
	})
}
