package workflow_service

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"backend/src/entities"
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

func TestCheckNewGitlabWorkflows(test *testing.T) {
	test.Run("Fail Find Service", func(test *testing.T) {
		mockServiceServiceRepo := new(MockServiceServiceRepository)

		gitlab := &WorkflowService{
			ServiceService: mockServiceServiceRepo,
		}

		mockServiceServiceRepo.On("FindServiceByName", "Gitlab").
			Return(entities.Service{}, errors.New("Fail find service"))

		err := gitlab.CheckNewGitlabWorkflows()

		require.EqualError(test, err, "Fail find service")
	})

	test.Run("Fail Find Actions", func(test *testing.T) {
		mockServiceServiceRepo := new(MockServiceServiceRepository)
		mockActionRepo := new(MockActionRepository)

		gitlab := &WorkflowService{
			ServiceService:   mockServiceServiceRepo,
			ActionRepository: mockActionRepo,
		}

		action := entities.Action{
			Id: "1",
		}

		serviceFound := entities.Service{
			Id: "1",
		}

		mockServiceServiceRepo.On("FindServiceByName", "Gitlab").
			Return(serviceFound, nil)

		mockActionRepo.On("FindActionsByServiceId", action.Id).
			Return([]entities.Action{}, errors.New("Fail find service"))

		err := gitlab.CheckNewGitlabWorkflows()

		require.EqualError(test, err, "Fail find service")
	})

	test.Run("Success", func(test *testing.T) {
		mockServiceServiceRepo := new(MockServiceServiceRepository)
		mockActionRepo := new(MockActionRepository)

		gitlab := &WorkflowService{
			ServiceService:   mockServiceServiceRepo,
			ActionRepository: mockActionRepo,
		}

		action := entities.Action{
			Id: "1",
		}

		serviceFound := entities.Service{
			Id: "1",
		}

		mockServiceServiceRepo.On("FindServiceByName", "Gitlab").
			Return(serviceFound, nil)

		mockActionRepo.On("FindActionsByServiceId", action.Id).
			Return([]entities.Action{}, nil)

		err := gitlab.CheckNewGitlabWorkflows()

		require.NoError(test, err)
	})
}
