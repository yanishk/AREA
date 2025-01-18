package workflow_service

import (
	"backend/src/entities"
	"errors"
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

func TestCheckNewWorkflowGitlabWebhook(test *testing.T) {
	test.Run("Fail Get Action Param", func(test *testing.T) {
		gitlab := &WorkflowService{}

		workflow := entities.Workflow{}

		err := gitlab.checkNewWorkflowGitlabWebhook(workflow, "accessToken")

		require.EqualError(test, err, errorMissingField)
	})
}

func TestCheckNewWorkflowsGitlabWebhook(test *testing.T) {
	test.Run("Fail Find Workflow", func(test *testing.T) {
		mockWorkflowRepo := new(MockWorkflowRepository)

		gitlab := &WorkflowService{
			WorkflowRepository: mockWorkflowRepo,
		}

		action := entities.Action{
			Id: "1",
		}

		mockWorkflowRepo.On("FindWorkflowsByActionId", action.Id).
			Return([]entities.Workflow{}, errors.New("Fail find workflows"))

		err := gitlab.checkNewWorkflowsGitlabWebhook(action)

		require.EqualError(test, err, "Fail find workflows")
	})

	test.Run("Success", func(test *testing.T) {
		mockWorkflowRepo := new(MockWorkflowRepository)

		gitlab := &WorkflowService{
			WorkflowRepository: mockWorkflowRepo,
		}

		action := entities.Action{
			Id: "1",
		}

		mockWorkflowRepo.On("FindWorkflowsByActionId", action.Id).
			Return([]entities.Workflow{}, nil)

		err := gitlab.checkNewWorkflowsGitlabWebhook(action)

		require.NoError(test, err)
	})
}
