package workflow_service

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"backend/src/entities"
)

func TestExecuteGithubRequest(test *testing.T) {
	test.Run("Missing Field", func(test *testing.T) {
		github := &WorkflowService{}

		_, err := github.executeGithubRequest(entities.Workflow{}, "method", "accessToken", "endPoint")

		require.EqualError(test, err, errorMissingField)
	})
}

func TestCheckActionDataLen(test *testing.T) {
	test.Run("Len Exist Success Update", func(test *testing.T) {
		mockWorkflowRepo := new(MockWorkflowRepository)

		github := &WorkflowService{
			WorkflowRepository: mockWorkflowRepo,
		}

		workflow := entities.Workflow{
			Id: "1",
			ActionData: map[string]interface{}{
				"len": 1.5,
			},
		}

		mockWorkflowRepo.On("UpdateWorkflow", workflow.Id, workflow).
			Return(nil)

		err := github.checkActionDataLen(workflow, 1.0)

		require.NoError(test, err)
	})

	test.Run("Len Exist Failure Update", func(test *testing.T) {
		mockWorkflowRepo := new(MockWorkflowRepository)

		github := &WorkflowService{
			WorkflowRepository: mockWorkflowRepo,
		}

		workflow := entities.Workflow{
			Id: "1",
			ActionData: map[string]interface{}{
				"len": 0.4,
			},
		}

		mockWorkflowRepo.On("UpdateWorkflow", workflow.Id, workflow).
			Return(errors.New("Fail update worklfow"))

		err := github.checkActionDataLen(workflow, 1.0)

		require.EqualError(test, err, "Fail update worklfow")
	})
}

func TestCheckGithubNewRepositoryAction(test *testing.T) {
	test.Run("Fail Request Repositories", func(test *testing.T) {
		mockServiceServiceRepo := new(MockServiceServiceRepository)

		github := &WorkflowService{
			ServiceService: mockServiceServiceRepo,
		}

		mockServiceServiceRepo.On("RequestGithubUserRepositories", "accessToken").
			Return([]entities.GithubRepository{}, errors.New("Fail request repositories"))

		err := github.checkGithubNewRepositoryAction(entities.Workflow{}, "accessToken")

		require.EqualError(test, err, "Fail request repositories")
	})
}

func TestCheckGithubNewIssueAction(test *testing.T) {
	test.Run("Fail Request Repositories", func(test *testing.T) {
		mockServiceServiceRepo := new(MockServiceServiceRepository)

		github := &WorkflowService{
			ServiceService: mockServiceServiceRepo,
		}

		mockResponse := &http.Response{
			Body: io.NopCloser(strings.NewReader(`{
				"title": "test"}`,
			)),
		}

		mockServiceServiceRepo.On("ExecuteApiRequest", githubBaseUrl+"issues", "GET", bearerType, "accessToken", nil).
			Return(mockResponse, errors.New("Fail request"))

		err := github.checkGithubNewIssueAction(entities.Workflow{}, "accessToken")

		require.EqualError(test, err, "Fail request")
	})

	test.Run("Fail Unmarshal", func(test *testing.T) {
		mockServiceServiceRepo := new(MockServiceServiceRepository)

		github := &WorkflowService{
			ServiceService: mockServiceServiceRepo,
		}

		mockResponse := &http.Response{
			Body: io.NopCloser(strings.NewReader(`{
				"title": "test"}`,
			)),
		}

		mockServiceServiceRepo.On("ExecuteApiRequest", githubBaseUrl+"issues", "GET", bearerType, "accessToken", nil).
			Return(mockResponse, nil)

		err := github.checkGithubNewIssueAction(entities.Workflow{}, "accessToken")

		require.Error(test, err)
	})
}

func TestCheckWorkflowsWithGithubActions(test *testing.T) {
	test.Run("Fail Find Workflows", func(test *testing.T) {
		mockWorkflowRepo := new(MockWorkflowRepository)

		github := &WorkflowService{
			WorkflowRepository: mockWorkflowRepo,
		}

		action := entities.Action{
			Id:   "1",
			Name: "test",
		}

		mockWorkflowRepo.On("FindWorkflowsByActionId", action.Id).
			Return([]entities.Workflow{}, errors.New("Fail find worklfows"))

		err := github.checkWorkflowsWithGithubActions(action)

		require.EqualError(test, err, "Fail find worklfows")
	})

	test.Run("Success", func(test *testing.T) {
		mockWorkflowRepo := new(MockWorkflowRepository)

		github := &WorkflowService{
			WorkflowRepository: mockWorkflowRepo,
		}

		action := entities.Action{
			Id:   "1",
			Name: "test",
		}

		mockWorkflowRepo.On("FindWorkflowsByActionId", action.Id).
			Return([]entities.Workflow{}, nil)

		err := github.checkWorkflowsWithGithubActions(action)

		require.NoError(test, err)
	})
}

func TestCheckGithubActions(test *testing.T) {
	test.Run("Fail Find Service", func(test *testing.T) {
		mockServiceServiceRepo := new(MockServiceServiceRepository)

		github := &WorkflowService{
			ServiceService: mockServiceServiceRepo,
		}

		mockServiceServiceRepo.On("FindServiceByName", "Github").
			Return(entities.Service{}, errors.New("Fail find service"))

		err := github.CheckGithubActions()

		require.EqualError(test, err, "Fail find service")
	})

	test.Run("Fail Find Actions", func(test *testing.T) {
		mockServiceServiceRepo := new(MockServiceServiceRepository)
		mockActionRepo := new(MockActionRepository)

		github := &WorkflowService{
			ServiceService:   mockServiceServiceRepo,
			ActionRepository: mockActionRepo,
		}

		service := entities.Service{
			Id: "1",
		}

		mockServiceServiceRepo.On("FindServiceByName", "Github").
			Return(service, nil)

		mockActionRepo.On("FindActionsByServiceId", service.Id).
			Return([]entities.Action{}, errors.New("Fail find actions"))

		err := github.CheckGithubActions()

		require.EqualError(test, err, "Fail find actions")
	})

	test.Run("Fail Find Actions", func(test *testing.T) {
		mockServiceServiceRepo := new(MockServiceServiceRepository)
		mockActionRepo := new(MockActionRepository)

		github := &WorkflowService{
			ServiceService:   mockServiceServiceRepo,
			ActionRepository: mockActionRepo,
		}

		service := entities.Service{
			Id: "1",
		}

		mockServiceServiceRepo.On("FindServiceByName", "Github").
			Return(service, nil)

		mockActionRepo.On("FindActionsByServiceId", service.Id).
			Return([]entities.Action{}, nil)

		err := github.CheckGithubActions()

		require.NoError(test, err)
	})
}

func TestCheckGithubWebhooksWorkflow(test *testing.T) {
	test.Run("Fail Get Action Param", func(test *testing.T) {
		github := &WorkflowService{}

		err := github.checkGithubWebhooksWorkflow(entities.Workflow{}, "webhookRepo")

		require.EqualError(test, err, errorMissingField)
	})
}

func TestIsGithubRepositoryWebhookPresent(test *testing.T) {
	test.Run("Fail Execute Request", func(test *testing.T) {
		mockServiceServiceRepo := new(MockServiceServiceRepository)

		github := &WorkflowService{
			ServiceService: mockServiceServiceRepo,
		}

		getWebhooksUrl := githubBaseUrl + githubRepositoryEndpoint + "repo" + "/hooks"

		mockServiceServiceRepo.On("ExecuteApiRequest", getWebhooksUrl, "GET", bearerType, "accessToken", nil).
			Return(&http.Response{}, errors.New("Fail Execute API"))

		_, err := github.isGithubRepositoryWebhookPresent("repo", "accessToken")

		require.EqualError(test, err, "Fail Execute API")
	})

	test.Run("Fail decode", func(test *testing.T) {
		mockServiceServiceRepo := new(MockServiceServiceRepository)

		github := &WorkflowService{
			ServiceService: mockServiceServiceRepo,
		}

		getWebhooksUrl := githubBaseUrl + githubRepositoryEndpoint + "repo" + "/hooks"

		mockResponse := &http.Response{
			Body: io.NopCloser(strings.NewReader(`{"name": "test"}`)),
		}

		mockServiceServiceRepo.On("ExecuteApiRequest", getWebhooksUrl, "GET", bearerType, "accessToken", nil).
			Return(mockResponse, nil)

		_, err := github.isGithubRepositoryWebhookPresent("repo", "accessToken")

		require.Error(test, err)
	})
}

func TestCheckNewWorkflowGithubWebhook(test *testing.T) {
	test.Run("Get Action Param Fail", func(test *testing.T) {
		github := &WorkflowService{}

		err := github.checkNewWorkflowGithubWebhook(entities.Workflow{}, "accessToken")

		require.EqualError(test, err, errorMissingField)
	})
}

func TestCheckNewWorkflowsGithubWebhook(test *testing.T) {
	test.Run("Fail Find Workflows", func(test *testing.T) {
		mockWorkflowRepo := new(MockWorkflowRepository)

		github := &WorkflowService{
			WorkflowRepository: mockWorkflowRepo,
		}

		action := entities.Action{
			Id: "1",
		}

		mockWorkflowRepo.On("FindWorkflowsByActionId", action.Id).
			Return([]entities.Workflow{}, errors.New("Fail find workflows"))

		err := github.checkNewWorkflowsGithubWebhook(action)

		require.EqualError(test, err, "Fail find workflows")
	})

	test.Run("Success", func(test *testing.T) {
		mockWorkflowRepo := new(MockWorkflowRepository)

		github := &WorkflowService{
			WorkflowRepository: mockWorkflowRepo,
		}

		action := entities.Action{
			Id: "1",
		}

		mockWorkflowRepo.On("FindWorkflowsByActionId", action.Id).
			Return([]entities.Workflow{}, nil)

		err := github.checkNewWorkflowsGithubWebhook(action)

		require.NoError(test, err)
	})
}

func TestCheckNewGithubWorkflows(test *testing.T) {
	test.Run("Success", func(test *testing.T) {
		mockServiceServiceRepo := new(MockServiceServiceRepository)
		mockActionRepo := new(MockActionRepository)

		github := &WorkflowService{
			ServiceService:   mockServiceServiceRepo,
			ActionRepository: mockActionRepo,
		}

		service := entities.Service{
			Id: "1",
		}

		mockServiceServiceRepo.On("FindServiceByName", "Github").
			Return(service, nil)

		mockActionRepo.On("FindActionsByServiceId", service.Id).
			Return([]entities.Action{}, nil)

		err := github.CheckNewGithubWorkflows()

		require.NoError(test, err)
	})

	test.Run("Fail Find Actions", func(test *testing.T) {
		mockServiceServiceRepo := new(MockServiceServiceRepository)
		mockActionRepo := new(MockActionRepository)

		github := &WorkflowService{
			ServiceService:   mockServiceServiceRepo,
			ActionRepository: mockActionRepo,
		}

		service := entities.Service{
			Id: "1",
		}

		mockServiceServiceRepo.On("FindServiceByName", "Github").
			Return(service, nil)

		mockActionRepo.On("FindActionsByServiceId", service.Id).
			Return([]entities.Action{}, errors.New("Fail find actions"))

		err := github.CheckNewGithubWorkflows()

		require.EqualError(test, err, "Fail find actions")
	})

	test.Run("Fail Find Actions", func(test *testing.T) {
		mockServiceServiceRepo := new(MockServiceServiceRepository)

		github := &WorkflowService{
			ServiceService: mockServiceServiceRepo,
		}

		mockServiceServiceRepo.On("FindServiceByName", "Github").
			Return(entities.Service{}, errors.New("Fail find service"))

		err := github.CheckNewGithubWorkflows()

		require.EqualError(test, err, "Fail find service")
	})
}
