package userservice_handler

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

type MockUserServiceService struct {
	mock.Mock
}

func (m *MockUserServiceService) RetrieveUserServiceAuthenticationStatus(email, connectionType, serviceName string) (bool, error) {
	args := m.Called(email, connectionType, serviceName)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserServiceService) CallApiAndRefresh(email, connectionType, serviceName string) (string, error) {
	args := m.Called(email, connectionType, serviceName)
	return args.String(0), args.Error(1)
}

func (m *MockUserServiceService) UpdateTokenForService(code, serviceName, appType, email, connectionType string) error {
	args := m.Called(code, serviceName, appType, email, connectionType)
	return args.Error(0)
}

func (m *MockUserServiceService) RetrieveGithubUserRepositories(email, connectionType string) ([]entities.GithubRepository, error) {
	args := m.Called(email, connectionType)
	return args.Get(0).([]entities.GithubRepository), args.Error(1)
}

func (m *MockUserServiceService) RetrieveGitlabUserProjects(email, connectionType string) ([]entities.GitlabProject, error) {
	args := m.Called(email, connectionType)
	return args.Get(0).([]entities.GitlabProject), args.Error(1)
}

func (m *MockUserServiceService) RetrieveDiscordUserServers(email, connectionType string) ([]map[string]interface{}, error) {
	args := m.Called(email, connectionType)
	return args.Get(0).([]map[string]interface{}), args.Error(1)
}

func (m *MockUserServiceService) RetrieveAsanaUserWorkspaces(email, connectionType string) (entities.AsanaWorkspacesInfo, error) {
	args := m.Called(email, connectionType)
	return args.Get(0).(entities.AsanaWorkspacesInfo), args.Error(1)
}

func (m *MockUserServiceService) RetrieveAsanaWorkspaceAssignees(email, connectionType, workspaceId string) (entities.AsanaWorkspacesInfo, error) {
	args := m.Called(email, connectionType, workspaceId)
	return args.Get(0).(entities.AsanaWorkspacesInfo), args.Error(1)
}

func (m *MockUserServiceService) RetrieveAsanaWorkspaceProjects(email, connectionType, workspaceId string) (entities.AsanaWorkspacesInfo, error) {
	args := m.Called(email, connectionType, workspaceId)
	return args.Get(0).(entities.AsanaWorkspacesInfo), args.Error(1)
}

func (m *MockUserServiceService) RetrieveAsanaWorkspaceTags(email, connectionType, workspaceId string) (entities.AsanaWorkspacesInfo, error) {
	args := m.Called(email, connectionType, workspaceId)
	return args.Get(0).(entities.AsanaWorkspacesInfo), args.Error(1)
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

func createMockAndRoute(isProtected bool) (*UserServiceHandler, *gin.Engine, *MockUserServiceService) {
	mockUserServiceService := new(MockUserServiceService)
	handler := &UserServiceHandler{UserServiceService: mockUserServiceService}

	router := gin.Default()
	if isProtected {
		router.Use(middleware.VerifyJWTCookie, middleware.VerifyEmailFromContext, middleware.VerifyConnectionTypeFromContext)
	}

	return handler, router, mockUserServiceService
}

func TestServiceCallback(test *testing.T) {
	handler, router, mockUserServiceService := createMockAndRoute(true)

	token := createToken(test)

	router.Use(func(c *gin.Context) {
		c.Set("email", "email")
		c.Set("connectionType", "basic")
	})
	router.POST("/service-callback", handler.serviceCallback)

	test.Run("Successful", func(test *testing.T) {
		mockUserServiceService.On("UpdateTokenForService", "code", "github", "web", "email", "basic").
			Return(nil).Once()

		body := `{
			"service": "github",
			"apptype": "web"
		}`

		req := requestForProtected("POST", "/service-callback?code=code", token, strings.NewReader(body))
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		require.Equal(test, http.StatusOK, w.Code)
		require.JSONEq(test, `{"success": "Token generated"}`, w.Body.String())

	})

	test.Run("Invalid code", func(test *testing.T) {
		mockUserServiceService.On("UpdateTokenForService", "code", "github", "web", "email", "basic").
			Return(nil).Once()

		body := `{
			"service": "github",
			"apptype": "web"
		}`

		req := requestForProtected("POST", "/service-callback?code=", token, strings.NewReader(body))
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		require.Equal(test, http.StatusBadRequest, w.Code)
		require.JSONEq(test, `{"error": "Invalid code authorization"}`, w.Body.String())

	})

	test.Run("Invalid app", func(test *testing.T) {
		mockUserServiceService.On("UpdateTokenForService", "code", "github", "fail", "email", "basic").
			Return(nil).Once()

		body := `{
			"service": "github",
			"apptype": "fail"
		}`

		req := requestForProtected("POST", "/service-callback?code=code", token, strings.NewReader(body))
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		require.Equal(test, http.StatusBadRequest, w.Code)
		require.JSONEq(test, `{"error": "Invalid app type"}`, w.Body.String())

	})

	test.Run("Fail update token", func(test *testing.T) {
		handler, router, mockUserServiceService := createMockAndRoute(true)

		token := createToken(test)

		router.Use(func(c *gin.Context) {
			c.Set("email", "email")
			c.Set("connectionType", "basic")
		})
		router.POST("/service-callback", handler.serviceCallback)

		mockUserServiceService.On("UpdateTokenForService", "code", "github", "web", "email", "basic").
			Return(errors.New("Failed to update token")).Once()

		body := `{
			"service": "github",
			"apptype": "web"
		}`

		req := requestForProtected("POST", "/service-callback?code=code", token, strings.NewReader(body))
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		require.Equal(test, http.StatusInternalServerError, w.Code)
		require.JSONEq(test, `{"error": "Failed to update token"}`, w.Body.String())
	})

	test.Run("Fail JSON Bind", func(test *testing.T) {
		mockUserServiceService.On("UpdateTokenForService", "code", "github", "web", "email", "basic").
			Return(nil).Once()

		req := requestForProtected("POST", "/service-callback?code=code", token, nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		require.Equal(test, http.StatusBadRequest, w.Code)
		require.JSONEq(test, `{"error": "Invalid request body"}`, w.Body.String())

	})
}

func TestGetUserServiceAuthenticationStatus(test *testing.T) {
	handler, router, mockUserServiceService := createMockAndRoute(true)

	token := createToken(test)

	router.Use(func(c *gin.Context) {
		c.Set("email", "email")
		c.Set("connectionType", "basic")
	})
	router.GET("/service-authentication-status", handler.getUserServiceAuthenticationStatus)

	test.Run("User is Authenticated", func(test *testing.T) {
		mockUserServiceService.On("RetrieveUserServiceAuthenticationStatus", "email", "basic", "Github").
			Return(true, nil).Once()

		req := requestForProtected("GET", "/service-authentication-status?service=Github", token, nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		require.Equal(test, http.StatusOK, w.Code)
		require.JSONEq(test, `{"authenticated": true}`, w.Body.String())
	})

	test.Run("User is not Authenticated", func(test *testing.T) {
		mockUserServiceService.On("RetrieveUserServiceAuthenticationStatus", "email", "basic", "Github").
			Return(false, nil).Once()

		req := requestForProtected("GET", "/service-authentication-status?service=Github", token, nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		require.Equal(test, http.StatusOK, w.Code)
		require.JSONEq(test, `{"authenticated": false}`, w.Body.String())
	})

	test.Run("Fail retrieving user's services", func(test *testing.T) {
		mockUserServiceService.On("RetrieveUserServiceAuthenticationStatus", "email", "basic", "Github").
			Return(false, errors.New("Could not retrieve user's services")).Once()

		req := requestForProtected("GET", "/service-authentication-status?service=Github", token, nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		require.Equal(test, http.StatusBadRequest, w.Code)
		require.JSONEq(test, `{"error": "Could not retrieve user's services"}`, w.Body.String())
	})
}

func TestGetGithubUserRepositories(test *testing.T) {
	handler, router, mockUserServiceService := createMockAndRoute(true)

	token := createToken(test)

	router.Use(func(c *gin.Context) {
		c.Set("email", "email")
		c.Set("connectionType", "basic")
	})
	router.GET("/github/user/repositories", handler.getGithubUserRepositories)

	test.Run("Successful", func(test *testing.T) {
		mockUserServiceService.On("RetrieveGithubUserRepositories", "email", "basic").
			Return([]entities.GithubRepository{}, nil).Once()

		req := requestForProtected("GET", "/github/user/repositories", token, nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		require.Equal(test, http.StatusOK, w.Code)
	})

	test.Run("Fail retrieve repo", func(test *testing.T) {
		mockUserServiceService.On("RetrieveGithubUserRepositories", "email", "basic").
			Return([]entities.GithubRepository{}, errors.New("Fail retrieve repo")).Once()

		req := requestForProtected("GET", "/github/user/repositories", token, nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		require.Equal(test, http.StatusInternalServerError, w.Code)
		require.JSONEq(test, `{"error": "Could not retrieve the user's repositories"}`, w.Body.String())
	})
}

func TestGetGitlabUserProjects(test *testing.T) {
	handler, router, mockUserServiceService := createMockAndRoute(true)

	token := createToken(test)

	router.Use(func(c *gin.Context) {
		c.Set("email", "email")
		c.Set("connectionType", "basic")
	})
	router.GET("/gitlab/user/projects", handler.getGitlabUserProjects)

	test.Run("Successful", func(test *testing.T) {
		mockUserServiceService.On("RetrieveGitlabUserProjects", "email", "basic").
			Return([]entities.GitlabProject{}, nil).Once()

		req := requestForProtected("GET", "/gitlab/user/projects", token, nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		require.Equal(test, http.StatusOK, w.Code)
	})

	test.Run("Fail retrieve project", func(test *testing.T) {
		mockUserServiceService.On("RetrieveGitlabUserProjects", "email", "basic").
			Return([]entities.GitlabProject{}, errors.New("Fail retrieve project")).Once()

		req := requestForProtected("GET", "/gitlab/user/projects", token, nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		require.Equal(test, http.StatusInternalServerError, w.Code)
		require.JSONEq(test, `{"error": "Could not retrieve the user's projects"}`, w.Body.String())
	})
}

func TestGetDiscordUserServers(test *testing.T) {
	handler, router, mockUserServiceService := createMockAndRoute(true)

	token := createToken(test)

	router.Use(func(c *gin.Context) {
		c.Set("email", "email")
		c.Set("connectionType", "basic")
	})
	router.GET("/discord/user/servers", handler.getDiscordUserServers)

	test.Run("Successful", func(test *testing.T) {
		mockUserServiceService.On("RetrieveDiscordUserServers", "email", "basic").
			Return([]map[string]interface{}{}, nil).Once()

		req := requestForProtected("GET", "/discord/user/servers", token, nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		require.Equal(test, http.StatusOK, w.Code)
	})

	test.Run("Fail to get servers", func(test *testing.T) {
		mockUserServiceService.On("RetrieveDiscordUserServers", "email", "basic").
			Return([]map[string]interface{}{}, errors.New("Fail retrieve servers")).Once()

		req := requestForProtected("GET", "/discord/user/servers", token, nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		require.Equal(test, http.StatusInternalServerError, w.Code)
		require.JSONEq(test, `{"error": "Could not retrieve the user's servers"}`, w.Body.String())
	})
}

func TestGetAsanaUserWorkspaces(test *testing.T) {
	handler, router, mockUserServiceService := createMockAndRoute(true)

	token := createToken(test)

	router.Use(func(c *gin.Context) {
		c.Set("email", "email")
		c.Set("connectionType", "basic")
	})
	router.GET("/asana/user/workspaces", handler.getAsanaUserWorkspaces)

	test.Run("Successful", func(test *testing.T) {
		mockWorkspaces := entities.AsanaWorkspacesInfo{
			Data: []struct {
				Gid  string `json:"gid"`
				Name string `json:"name"`
			}{
				{Gid: "gid", Name: "name"},
			},
		}

		mockUserServiceService.On("RetrieveAsanaUserWorkspaces", "email", "basic").
			Return(mockWorkspaces, nil).Once()

		req := requestForProtected("GET", "/asana/user/workspaces", token, nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		require.Equal(test, http.StatusOK, w.Code)
	})

	test.Run("Empty Data", func(test *testing.T) {
		mockUserServiceService.On("RetrieveAsanaUserWorkspaces", "email", "basic").
			Return(entities.AsanaWorkspacesInfo{}, nil).Once()

		req := requestForProtected("GET", "/asana/user/workspaces", token, nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		require.Equal(test, http.StatusInternalServerError, w.Code)
		require.JSONEq(test, `{"error": "Could not retrieve the user's workspaces"}`, w.Body.String())
	})

	test.Run("Fail get Workspace", func(test *testing.T) {
		mockWorkspaces := entities.AsanaWorkspacesInfo{
			Data: []struct {
				Gid  string `json:"gid"`
				Name string `json:"name"`
			}{
				{Gid: "gid", Name: "name"},
			},
		}

		mockUserServiceService.On("RetrieveAsanaUserWorkspaces", "email", "basic").
			Return(mockWorkspaces, errors.New("Fail get workspaces")).Once()

		req := requestForProtected("GET", "/asana/user/workspaces", token, nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		require.Equal(test, http.StatusInternalServerError, w.Code)
		require.JSONEq(test, `{"error": "Could not retrieve the user's workspaces"}`, w.Body.String())
	})
}

func TestGetAsanaWorkspaceAssignees(test *testing.T) {
	handler, router, mockUserServiceService := createMockAndRoute(true)

	token := createToken(test)

	router.Use(func(c *gin.Context) {
		c.Set("email", "email")
		c.Set("connectionType", "basic")
	})
	router.GET("/asana/workspace/assignees", handler.getAsanaWorkspaceAssignees)

	test.Run("Successful", func(test *testing.T) {
		mockUserServiceService.On("RetrieveAsanaWorkspaceAssignees", "email", "basic", "id").
			Return(entities.AsanaWorkspacesInfo{}, nil).Once()

		req := requestForProtected("GET", "/asana/workspace/assignees?id=id", token, nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		require.Equal(test, http.StatusOK, w.Code)
	})

	test.Run("Fail get assignees", func(test *testing.T) {
		mockUserServiceService.On("RetrieveAsanaWorkspaceAssignees", "email", "basic", "id").
			Return(entities.AsanaWorkspacesInfo{}, errors.New("Fail retrieve assignees")).Once()

		req := requestForProtected("GET", "/asana/workspace/assignees?id=id", token, nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		require.Equal(test, http.StatusInternalServerError, w.Code)
		require.JSONEq(test, `{"error": "Could not retrieve workspace's assignees"}`, w.Body.String())
	})
}

func TestGetAsanaWorkspaceProjects(test *testing.T) {
	handler, router, mockUserServiceService := createMockAndRoute(true)

	token := createToken(test)

	router.Use(func(c *gin.Context) {
		c.Set("email", "email")
		c.Set("connectionType", "basic")
	})
	router.GET("/asana/workspace/projects", handler.getAsanaWorkspaceProjects)

	test.Run("Successful", func(test *testing.T) {
		mockUserServiceService.On("RetrieveAsanaWorkspaceProjects", "email", "basic", "id").
			Return(entities.AsanaWorkspacesInfo{}, nil).Once()

		req := requestForProtected("GET", "/asana/workspace/projects?id=id", token, nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		require.Equal(test, http.StatusOK, w.Code)
	})

	test.Run("Fail retrieve project", func(test *testing.T) {
		mockUserServiceService.On("RetrieveAsanaWorkspaceProjects", "email", "basic", "id").
			Return(entities.AsanaWorkspacesInfo{}, errors.New("Fail retrieve project")).Once()

		req := requestForProtected("GET", "/asana/workspace/projects?id=id", token, nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		require.Equal(test, http.StatusInternalServerError, w.Code)
		require.JSONEq(test, `{"error": "Could not retrieve workspace's projects"}`, w.Body.String())
	})
}

func TestGetAsanaWorkspaceTags(test *testing.T) {
	handler, router, mockUserServiceService := createMockAndRoute(true)

	token := createToken(test)

	router.Use(func(c *gin.Context) {
		c.Set("email", "email")
		c.Set("connectionType", "basic")
	})
	router.GET("/asana/workspace/tags", handler.getAsanaWorkspaceTags)

	test.Run("Successful", func(test *testing.T) {
		mockUserServiceService.On("RetrieveAsanaWorkspaceTags", "email", "basic", "id").
			Return(entities.AsanaWorkspacesInfo{}, nil).Once()

		req := requestForProtected("GET", "/asana/workspace/tags?id=id", token, nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		require.Equal(test, http.StatusOK, w.Code)
	})

	test.Run("Fail retrieve tags", func(test *testing.T) {
		mockUserServiceService.On("RetrieveAsanaWorkspaceTags", "email", "basic", "id").
			Return(entities.AsanaWorkspacesInfo{}, errors.New("Fail retrieve tags")).Once()

		req := requestForProtected("GET", "/asana/workspace/tags?id=id", token, nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		require.Equal(test, http.StatusInternalServerError, w.Code)
		require.JSONEq(test, `{"error": "Could not retrieve workspace's tags"}`, w.Body.String())
	})
}
