package workflow_service

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"backend/src/entities"
)

func TestPostRedditComment(test *testing.T) {
	test.Run("Missing Field", func(t *testing.T) {
		mockReactionRepo := new(MockReactionRepository)

		reddit := &WorkflowService{
			ReactionRepository: mockReactionRepo,
		}

		err := reddit.postRedditComment("accessToken", entities.Workflow{})

		require.EqualError(test, err, errorMissingField)
	})
}

func TestUpvoteRedditPost(test *testing.T) {
	test.Run("Missing Field", func(t *testing.T) {
		mockReactionRepo := new(MockReactionRepository)

		reddit := &WorkflowService{
			ReactionRepository: mockReactionRepo,
		}

		err := reddit.upvoteRedditPost("accessToken", entities.Workflow{})

		require.EqualError(test, err, errorMissingField)
	})
}

func TestDownvoteRedditPost(test *testing.T) {
	test.Run("Missing Field", func(t *testing.T) {
		mockReactionRepo := new(MockReactionRepository)

		reddit := &WorkflowService{
			ReactionRepository: mockReactionRepo,
		}

		err := reddit.downvoteRedditPost("accessToken", entities.Workflow{})

		require.EqualError(test, err, errorMissingField)
	})
}

func TestSubmitRedditPost(test *testing.T) {
	test.Run("Missing Field", func(t *testing.T) {
		mockReactionRepo := new(MockReactionRepository)

		reddit := &WorkflowService{
			ReactionRepository: mockReactionRepo,
		}

		err := reddit.submitRedditPost("accessToken", entities.Workflow{})

		require.EqualError(test, err, errorMissingField)
	})
}

func TestCheckRedditReactions(test *testing.T) {
	workflow := entities.Workflow{
		ReactionId: "1",
	}

	test.Run("Success", func(t *testing.T) {
		mockReactionRepo := new(MockReactionRepository)

		reddit := &WorkflowService{
			ReactionRepository: mockReactionRepo,
		}

		reactionFound := entities.Reaction{
			Name: "Test",
		}

		mockReactionRepo.On("FindReactionById", workflow.ReactionId).
			Return(reactionFound, nil)

		err := reddit.checkRedditReactions(workflow)

		require.NoError(test, err)
	})

	test.Run("Fail Find Reaction", func(t *testing.T) {
		mockReactionRepo := new(MockReactionRepository)

		reddit := &WorkflowService{
			ReactionRepository: mockReactionRepo,
		}

		mockReactionRepo.On("FindReactionById", workflow.ReactionId).
			Return(entities.Reaction{}, errors.New(errorRetrievingReaction))

		err := reddit.checkRedditReactions(workflow)

		require.EqualError(test, err, errorRetrievingReaction)
	})
}
