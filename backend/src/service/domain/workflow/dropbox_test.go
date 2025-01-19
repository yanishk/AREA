package workflow_service

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"backend/src/entities"
)

func TestMoveFileOrFolder(test *testing.T) {
	test.Run("Missing Field", func(test *testing.T) {
		dropbox := &WorkflowService{}

		workflow := entities.Workflow{
			ReactionParam: map[string]interface{}{
				"key": "value",
			},
		}

		err := dropbox.moveFileOrFolder("accessToken", workflow)

		require.EqualError(test, err, errorMissingField)
	})
}

func TestCreateTextFile(test *testing.T) {
	test.Run("Missing Field", func(test *testing.T) {
		dropbox := &WorkflowService{}

		workflow := entities.Workflow{
			ReactionParam: map[string]interface{}{
				"key": "value",
			},
		}

		err := dropbox.createTextFile("accessToken", workflow)

		require.EqualError(test, err, errorMissingField)
	})
}

func TestAppendTextFile(test *testing.T) {
	test.Run("Missing Field", func(test *testing.T) {
		dropbox := &WorkflowService{}

		workflow := entities.Workflow{
			ReactionParam: map[string]interface{}{
				"key": "value",
			},
		}

		err := dropbox.appendTextFile("accessToken", workflow)

		require.EqualError(test, err, errorMissingField)
	})
}

func TestCheckDropboxReactions(test *testing.T) {
	test.Run("Fail Find Reaction", func(test *testing.T) {
		mockReactionRepo := new(MockReactionRepository)

		dropbox := &WorkflowService{
			ReactionRepository: mockReactionRepo,
		}

		workflow := entities.Workflow{
			ReactionId: "1",
		}

		mockReactionRepo.On("FindReactionById", workflow.ReactionId).
			Return(entities.Reaction{}, errors.New(errorRetrievingReaction))

		err := dropbox.checkDropboxReactions(workflow)

		require.EqualError(test, err, errorRetrievingReaction)
	})

	test.Run("Success", func(test *testing.T) {
		mockReactionRepo := new(MockReactionRepository)

		dropbox := &WorkflowService{
			ReactionRepository: mockReactionRepo,
		}

		workflow := entities.Workflow{
			ReactionId: "1",
		}

		reaction := entities.Reaction{
			Name: "test",
		}

		mockReactionRepo.On("FindReactionById", workflow.ReactionId).
			Return(reaction, nil)

		err := dropbox.checkDropboxReactions(workflow)

		require.NoError(test, err)
	})
}
