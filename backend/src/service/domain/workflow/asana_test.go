package workflow_service

import (
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
