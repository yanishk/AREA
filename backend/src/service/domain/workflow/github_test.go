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
