package workflow_service

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"backend/src/entities"
)

func TestCreateTaskAsana(test *testing.T) {
	test.Run("Missing Field", func(test *testing.T) {
		asana := &WorkflowService{}

		workflow := entities.Workflow{
			ReactionParam: map[string]interface{}{},
		}

		err := asana.createTaskAsana("accessToken", workflow)

		require.EqualError(test, err, errorMissingField)
	})
}

func TestCreateProjectAsana(test *testing.T) {
	test.Run("Missing Field", func(test *testing.T) {
		asana := &WorkflowService{}

		workflow := entities.Workflow{
			ReactionParam: map[string]interface{}{},
		}

		err := asana.createProjectAsana("accessToken", workflow)

		require.EqualError(test, err, errorMissingField)
	})
}

func TestCheckAsanaReactions(test *testing.T) {
	test.Run("Success", func(test *testing.T) {
		mockReactionRepo := new(MockReactionRepository)

		asana := &WorkflowService{
			ReactionRepository: mockReactionRepo,
		}

		workflow := entities.Workflow{
			ReactionId: "1",
		}

		mockReactionRepo.On("FindReactionById", workflow.ReactionId).
			Return(entities.Reaction{}, nil)

		err := asana.checkAsanaReactions(workflow)

		require.NoError(test, err)
	})

	test.Run("Fail Find Reaction", func(test *testing.T) {
		mockReactionRepo := new(MockReactionRepository)

		asana := &WorkflowService{
			ReactionRepository: mockReactionRepo,
		}

		workflow := entities.Workflow{
			ReactionId: "1",
		}

		mockReactionRepo.On("FindReactionById", workflow.ReactionId).
			Return(entities.Reaction{}, errors.New(errorRetrievingReaction))

		err := asana.checkAsanaReactions(workflow)

		require.EqualError(test, err, errorRetrievingReaction)
	})
}
