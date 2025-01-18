package workflow_service

import (
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
