package workflow_handler

import (
	"backend/src/entities"
	"backend/src/handler/middleware"
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
		"email":          "eva.legrand@epitech.eu",
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
}
