package workflow_service

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"backend/src/entities"
)

func TestSendMeEmail(test *testing.T) {
	test.Run("Fail Find User", func(test *testing.T) {
		mockUserRepo := new(MockUserRepository)

		email := &WorkflowService{
			UserRepository: mockUserRepo,
		}

		workflow := entities.Workflow{
			OwnerId: "1",
			ReactionParam: map[string]interface{}{
				"key": "value",
			},
		}

		mockUserRepo.On("FindUserById", workflow.OwnerId).
			Return(entities.User{}, errors.New("Fail find user"))

		err := email.sendMeEmail(workflow)

		require.EqualError(test, err, "Error finding user")
	})

	test.Run("Missing fields", func(test *testing.T) {
		mockUserRepo := new(MockUserRepository)

		email := &WorkflowService{
			UserRepository: mockUserRepo,
		}

		workflow := entities.Workflow{
			OwnerId: "1",
			ReactionParam: map[string]interface{}{
				"key": "value",
			},
		}

		mockUserRepo.On("FindUserById", workflow.OwnerId).
			Return(entities.User{}, nil)

		err := email.sendMeEmail(workflow)

		require.EqualError(test, err, errorMissingField)
	})
}

func TestCheckSendEmailReactions(test *testing.T) {
	test.Run("Fail Find Reaction", func(test *testing.T) {
		mockReactionRepo := new(MockReactionRepository)

		email := &WorkflowService{
			ReactionRepository: mockReactionRepo,
		}

		workflow := entities.Workflow{
			ReactionId: "1",
		}

		mockReactionRepo.On("FindReactionById", workflow.ReactionId).
			Return(entities.Reaction{}, errors.New(errorRetrievingReaction))

		err := email.checkSendEmailReactions(workflow)

		require.EqualError(test, err, errorRetrievingReaction)
	})

	test.Run("Success", func(test *testing.T) {
		mockReactionRepo := new(MockReactionRepository)

		email := &WorkflowService{
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

		err := email.checkSendEmailReactions(workflow)

		require.NoError(test, err)
	})
}
