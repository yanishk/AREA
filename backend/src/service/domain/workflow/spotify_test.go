package workflow_service

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"backend/src/entities"
)

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

func TestGetSpotifyUrl(test *testing.T) {
	workflow := entities.Workflow{
		ReactionParam: map[string]interface{}{"key": "value"},
	}
	baseUrl := "base"
	endUrl := "endurl"

	test.Run("Success", func(test *testing.T) {
		parameterName := "key"

		expectedUrl := baseUrl + workflow.ReactionParam[parameterName].(string) + endUrl

		res, err := getSpotifyUrl(workflow, baseUrl, parameterName, endUrl)

		require.NoError(test, err)
		assert.Equal(test, expectedUrl, res)
	})

	test.Run("Failure", func(test *testing.T) {
		_, err := getSpotifyUrl(workflow, baseUrl, "fail", endUrl)

		require.EqualError(test, err, "Incorrect parameter")
	})
}

func TestExecuteSpotifyRequest(test *testing.T) {
	workflow := entities.Workflow{
		ReactionParam: map[string]interface{}{"key": "value"},
	}

	test.Run("Failure", func(test *testing.T) {
		spotify := &WorkflowService{}
		spotifyParam := entities.SpotifyRequestParameters{
			BaseUrl:       "baseUrl",
			EndUrl:        "endUrl",
			ParameterName: "Fail",
		}

		_, err := spotify.executeSpotifyRequest(workflow, spotifyParam, "accessToken")

		require.EqualError(test, err, "Incorrect parameter")
	})
}

func TestCheckSpotifyReactions(test *testing.T) {
	workflow := entities.Workflow{
		Id:         "1",
		ReactionId: "1",
	}

	test.Run("Fail Retrieve Reaction", func(test *testing.T) {
		mockReactionRepo := new(MockReactionRepository)

		spotify := &WorkflowService{
			ReactionRepository: mockReactionRepo,
		}

		mockReactionRepo.On("FindReactionById", workflow.ReactionId).
			Return(entities.Reaction{}, errors.New(errorRetrievingReaction))

		err := spotify.checkSpotifyReactions(workflow)

		require.EqualError(test, err, errorRetrievingReaction)
	})

	test.Run("Unknown Reaction", func(test *testing.T) {
		mockReactionRepo := new(MockReactionRepository)

		spotify := &WorkflowService{
			ReactionRepository: mockReactionRepo,
		}

		mockReactionRepo.On("FindReactionById", workflow.ReactionId).
			Return(entities.Reaction{}, nil)

		err := spotify.checkSpotifyReactions(workflow)

		require.EqualError(test, err, "Unknown reaction")
	})
}
