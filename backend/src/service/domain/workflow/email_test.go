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
