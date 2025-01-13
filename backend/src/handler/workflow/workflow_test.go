package workflow_handler

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"backend/src/entities"
	"backend/src/handler/middleware"
)

type MockWorkflowService struct {
	mock.Mock
}

func (m *MockWorkflowService) CreateWorkflow(userEmail, userConnectionType string, newWorkflow entities.NewWorkflow) error {
	args := m.Called(userEmail, userConnectionType, newWorkflow)
	return args.Error(0)
}

func (m *MockWorkflowService) GetUserWorkflows(email, connectionType string) ([]entities.Workflow, error) {
	args := m.Called(email, connectionType)
	return args.Get(0).([]entities.Workflow), args.Error(1)
}

func (m *MockWorkflowService) UpdateWorkflow(workflowId string, workflow entities.UpdatedWorkflow) error {
	args := m.Called(workflowId, workflow)
	return args.Error(0)
}

func (m *MockWorkflowService) DeleteWorkflow(email, connectionType, workflowId string) error {
	args := m.Called(email, connectionType, workflowId)
	return args.Error(0)
}

func (m *MockWorkflowService) CheckTimeAndDateActions() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockWorkflowService) CheckGithubActions() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockWorkflowService) CheckRedditActions() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockWorkflowService) CheckWeatherActions() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockWorkflowService) CheckNewGitlabWorkflows() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockWorkflowService) CheckNewGithubWorkflows() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockWorkflowService) CheckWebhooksWorkflows(serviceName string, request *http.Request) error {
	args := m.Called(serviceName, request)
	return args.Error(0)
}

func requestForProtected(method, url, token string, body io.Reader) *http.Request {
	req, _ := http.NewRequest(method, url, body)
	req.AddCookie(&http.Cookie{Name: "JWToken", Value: token})
	return req
}

func createToken(test *testing.T) string {
	secretKey := os.Getenv("SECRET_KEY")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":          "test@test.com",
		"connectionType": "basic",
		"exp":            time.Now().Add(time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(secretKey))
	require.NoError(test, err)

	return tokenString
}

func createMockAndRoute(isProtected bool) (*WorkflowHandler, *gin.Engine, *MockWorkflowService) {
	MockWorkflowService := new(MockWorkflowService)
	handler := &WorkflowHandler{WorkflowService: MockWorkflowService}

	router := gin.Default()
	if isProtected {
		router.Use(middleware.VerifyJWTCookie, middleware.VerifyEmailFromContext, middleware.VerifyConnectionTypeFromContext)
	}

	return handler, router, MockWorkflowService
}

func TestCreateWorkflow(test *testing.T) {
	handler, router, mock := createMockAndRoute(true)

	token := createToken(test)

	router.Use(func(c *gin.Context) {
		c.Set("email", "email")
		c.Set("connectionType", "basic")
	})
	router.POST("/workflows", handler.createWorkflow)

	test.Run("Successful", func(test *testing.T) {
		var newWorkflow entities.NewWorkflow

		mock.On("CreateWorkflow", "email", "basic", newWorkflow).
			Return(nil).Once()

		body := `{
			"service": "github",
			"apptype": "web"
		}`

		req := requestForProtected("POST", "/workflows", token, strings.NewReader(body))
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		require.Equal(test, http.StatusOK, w.Code)
		require.JSONEq(test, `{"success": "Successful workflow creation"}`, w.Body.String())

	})

	test.Run("Fail creation workflow", func(test *testing.T) {
		var newWorkflow entities.NewWorkflow

		mock.On("CreateWorkflow", "email", "basic", newWorkflow).
			Return(errors.New("Fail create workflow")).Once()

		body := `{
			"service": "github",
			"apptype": "web"
		}`

		req := requestForProtected("POST", "/workflows", token, strings.NewReader(body))
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		require.Equal(test, http.StatusInternalServerError, w.Code)
		require.JSONEq(test, `{"error": "Could not create workflow"}`, w.Body.String())

	})

	test.Run("Fail JSON Bind", func(test *testing.T) {
		var newWorkflow entities.NewWorkflow

		mock.On("CreateWorkflow", "email", "basic", newWorkflow).
			Return(nil).Once()

		req := requestForProtected("POST", "/workflows", token, nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		require.Equal(test, http.StatusBadRequest, w.Code)
		require.JSONEq(test, `{"error": "Invalid request body"}`, w.Body.String())

	})
}

func TestGetUserWorkflows(test *testing.T) {
	handler, router, mock := createMockAndRoute(true)

	token := createToken(test)

	router.Use(func(c *gin.Context) {
		c.Set("email", "email")
		c.Set("connectionType", "basic")
	})
	router.GET("/workflows", handler.getUserWorkflows)

	test.Run("Successful", func(test *testing.T) {
		mock.On("GetUserWorkflows", "email", "basic").
			Return([]entities.Workflow{}, nil).Once()

		req := requestForProtected("GET", "/workflows", token, nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		require.Equal(test, http.StatusOK, w.Code)
	})

	test.Run("Fail retrieve workflows", func(test *testing.T) {
		mock.On("GetUserWorkflows", "email", "basic").
			Return([]entities.Workflow{}, errors.New("Fail retrieve workflow")).Once()

		req := requestForProtected("GET", "/workflows", token, nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		require.Equal(test, http.StatusInternalServerError, w.Code)
		require.JSONEq(test, `{"error": "Could not retrieve user's workflows"}`, w.Body.String())
	})
}

func TestUpdateWorkflows(test *testing.T) {
	handler, router, mock := createMockAndRoute(true)

	token := createToken(test)

	router.Use(func(c *gin.Context) {
		c.Set("email", "email")
		c.Set("connectionType", "basic")
	})
	router.PUT("/workflows/:id", handler.updateWorkflow)

	test.Run("Successful", func(test *testing.T) {
		var workflow entities.UpdatedWorkflow

		body := `{
			"key": "value",
			"key": "value"
		}`

		mock.On("UpdateWorkflow", "1", workflow).
			Return(nil).Once()

		req := requestForProtected("PUT", "/workflows/1", token, strings.NewReader(body))
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		require.Equal(test, http.StatusOK, w.Code)
		require.JSONEq(test, `{"success": "Workflow successfully updated"}`, w.Body.String())
	})

	test.Run("Fail updating workflow", func(test *testing.T) {
		var workflow entities.UpdatedWorkflow

		body := `{
			"key": "value",
			"key": "value"
		}`

		mock.On("UpdateWorkflow", "1", workflow).
			Return(errors.New("Fail update workflow")).Once()

		req := requestForProtected("PUT", "/workflows/1", token, strings.NewReader(body))
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		require.Equal(test, http.StatusInternalServerError, w.Code)
		require.JSONEq(test, `{"error": "Could not update workflow"}`, w.Body.String())
	})

	test.Run("Fail JSON Bind", func(test *testing.T) {
		var workflow entities.UpdatedWorkflow

		mock.On("UpdateWorkflow", "1", workflow).
			Return(nil).Once()

		req := requestForProtected("PUT", "/workflows/1", token, nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		require.Equal(test, http.StatusBadRequest, w.Code)
		require.JSONEq(test, `{"error": "Invalid request body"}`, w.Body.String())
	})
}

func TestDeleteWorkflow(test *testing.T) {
	handler, router, mock := createMockAndRoute(true)

	token := createToken(test)

	router.Use(func(c *gin.Context) {
		c.Set("email", "email")
		c.Set("connectionType", "basic")
	})
	router.DELETE("/workflows/:id", handler.deleteWorkflow)

	test.Run("Successful", func(test *testing.T) {
		mock.On("DeleteWorkflow", "email", "basic", "1").
			Return(nil).Once()

		req := requestForProtected("DELETE", "/workflows/1", token, nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		require.Equal(test, http.StatusOK, w.Code)
		require.JSONEq(test, `{"success": "Workflow deleted"}`, w.Body.String())
	})

	test.Run("Fail delete", func(test *testing.T) {
		mock.On("DeleteWorkflow", "email", "basic", "1").
			Return(errors.New("Fail delete")).Once()

		req := requestForProtected("DELETE", "/workflows/1", token, nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		require.Equal(test, http.StatusInternalServerError, w.Code)
		require.JSONEq(test, `{"error": "Could not delete workflow"}`, w.Body.String())
	})
}
