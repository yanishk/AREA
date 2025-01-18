package workflow_service

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"backend/src/entities"
)

func TestPostDiscordMessage(test *testing.T) {
	test.Run("Missing Field", func(test *testing.T) {
		discord := &WorkflowService{}

		workflow := entities.Workflow{
			ReactionParam: map[string]interface{}{
				"key": "value",
			},
		}

		err := discord.postDiscordMessage("tokenBot", workflow)

		require.EqualError(test, err, errorMissingField)
	})
}

func TestCreateDiscordThread(test *testing.T) {
	test.Run("Missing Field", func(test *testing.T) {
		discord := &WorkflowService{}

		workflow := entities.Workflow{
			ReactionParam: map[string]interface{}{
				"key": "value",
			},
		}

		err := discord.createDiscordThread("tokenBot", workflow)

		require.EqualError(test, err, errorMissingField)
	})
}

func TestCheckDiscordReactions(test *testing.T) {
	test.Run("Fail Find Reaction", func(test *testing.T) {
		mockReactionRepo := new(MockReactionRepository)

		discord := &WorkflowService{
			ReactionRepository: mockReactionRepo,
		}

		workflow := entities.Workflow{
			ReactionId: "1",
		}

		mockReactionRepo.On("FindReactionById", workflow.ReactionId).
			Return(entities.Reaction{}, errors.New(errorRetrievingReaction))

		err := discord.checkDiscordReactions(workflow)

		require.EqualError(test, err, errorRetrievingReaction)
	})

	test.Run("Success", func(test *testing.T) {
		mockReactionRepo := new(MockReactionRepository)

		discord := &WorkflowService{
			ReactionRepository: mockReactionRepo,
		}

		workflow := entities.Workflow{
			ReactionId: "1",
		}

		mockReactionRepo.On("FindReactionById", workflow.ReactionId).
			Return(entities.Reaction{}, nil)

		err := discord.checkDiscordReactions(workflow)

		require.EqualError(test, err, "Unknown reaction")
	})
}
