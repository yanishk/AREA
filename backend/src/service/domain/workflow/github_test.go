package workflow_service

import (
	"backend/src/entities"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExecuteGithubRequest(test *testing.T) {
	test.Run("Missing Field", func(test *testing.T) {
		github := &WorkflowService{}

		_, err := github.executeGithubRequest(entities.Workflow{}, "method", "accessToken", "endPoint")

		require.EqualError(test, err, errorMissingField)
	})
}
