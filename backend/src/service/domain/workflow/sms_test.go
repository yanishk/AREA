package workflow_service

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"backend/src/entities"
)

func TestSendAnSMS(test *testing.T) {
	test.Run("Missing Field", func(test *testing.T) {
		spotify := &WorkflowService{}

		workflow := entities.Workflow{
			ReactionParam: map[string]interface{}{"key": "value"},
		}

		err := spotify.sendAnSMS(workflow)

		require.EqualError(test, err, errorMissingField)
	})
}

func TestCheckSMSReactions(test *testing.T) {
	test.Run("Success", func(test *testing.T) {
		mockReactionRepo := new(MockReactionRepository)

		spotify := &WorkflowService{
			ReactionRepository: mockReactionRepo,
		}

		workflow := entities.Workflow{
			ReactionId: "1",
		}

		mockReactionRepo.On("FindReactionById", workflow.ReactionId).
			Return(entities.Reaction{}, nil)

		err := spotify.checkSMSReactions(workflow)

		require.NoError(test, err)
	})
	test.Run("Fail Find Reaction", func(test *testing.T) {
		mockReactionRepo := new(MockReactionRepository)

		spotify := &WorkflowService{
			ReactionRepository: mockReactionRepo,
		}

		workflow := entities.Workflow{
			ReactionId: "1",
		}

		mockReactionRepo.On("FindReactionById", workflow.ReactionId).
			Return(entities.Reaction{}, errors.New(errorRetrievingReaction))

		err := spotify.checkSMSReactions(workflow)

		require.EqualError(test, err, errorRetrievingReaction)
	})
}
