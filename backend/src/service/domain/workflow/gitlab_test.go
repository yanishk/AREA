package workflow_service

import (
	"backend/src/entities"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCheckGitlabWebhooksWorkflow(test *testing.T) {
	test.Run("Fail Get Action Param", func(test *testing.T) {
		gitlab := &WorkflowService{}

		workflow := entities.Workflow{}

		err := gitlab.checkGitlabWebhooksWorkflow(workflow, 1)

		require.EqualError(test, err, errorMissingField)
	})

	test.Run("Success", func(test *testing.T) {
		gitlab := &WorkflowService{}

		workflow := entities.Workflow{
			ActionParam: map[string]interface{}{
				"project": "value",
			},
		}

		err := gitlab.checkGitlabWebhooksWorkflow(workflow, 1)

		require.NoError(test, err)
	})
}
