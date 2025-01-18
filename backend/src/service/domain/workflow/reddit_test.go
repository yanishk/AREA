package workflow_service

import (
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
