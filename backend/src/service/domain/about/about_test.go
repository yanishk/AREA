package about_service

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

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
