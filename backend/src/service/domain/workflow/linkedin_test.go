package workflow_service

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"backend/src/entities"
)

func TestPublishTextLinkedinPost(test *testing.T) {
	test.Run("Missing Field", func(test *testing.T) {
		linkedin := &WorkflowService{}

		workflow := entities.Workflow{
			ReactionParam: map[string]interface{}{"key": "value"},
		}

		err := linkedin.publishTextLinkedinPost("urn", "accessToken", workflow)

		require.EqualError(test, err, errorMissingField)
	})
}

func TestPublishURLLinkedinPost(test *testing.T) {
	test.Run("Missing Field", func(test *testing.T) {
		linkedin := &WorkflowService{}

		workflow := entities.Workflow{
			ReactionParam: map[string]interface{}{"key": "value"},
		}

		err := linkedin.publishURLLinkedinPost("urn", "accessToken", workflow)

		require.EqualError(test, err, errorMissingField)
	})
}

func TestCheckLinkedinReactions(test *testing.T) {
	test.Run("Fail Find Reactions", func(test *testing.T) {
		mockReactionRepo := new(MockReactionRepository)

		linkedin := &WorkflowService{
			ReactionRepository: mockReactionRepo,
		}

		workflow := entities.Workflow{
			ReactionId: "1",
		}

		mockReactionRepo.On("FindReactionById", workflow.ReactionId).
			Return(entities.Reaction{}, errors.New(errorRetrievingReaction))

		err := linkedin.checkLinkedinReactions(workflow)

		require.EqualError(test, err, errorRetrievingReaction)
	})

	test.Run("Success", func(test *testing.T) {
		mockReactionRepo := new(MockReactionRepository)

		linkedin := &WorkflowService{
			ReactionRepository: mockReactionRepo,
		}

		workflow := entities.Workflow{
			ReactionId: "1",
		}

		mockReactionRepo.On("FindReactionById", workflow.ReactionId).
			Return(entities.Reaction{}, nil)

		err := linkedin.checkLinkedinReactions(workflow)

		require.NoError(test, err)
	})
}
