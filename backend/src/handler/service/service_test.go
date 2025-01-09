package service_handler

import (
	"backend/src/entities"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockServiceService struct {
	mock.Mock
}

func (m *MockServiceService) OAuth2Service(serviceName, callbackType, appType string) (string, error) {
	args := m.Called(serviceName, callbackType, appType)
	return args.String(0), args.Error(1)
}

func (m *MockServiceService) RetrieveActionsServices() ([]entities.Service, error) {
	return nil, nil
}

func (m *MockServiceService) RetrieveReactionsServices() ([]entities.Service, error) {
	return nil, nil
}

func (m *MockServiceService) FindServiceByName(name string) (entities.Service, error) {
	var test entities.Service
	return test, nil
}

func (m *MockServiceService) FindServiceById(id string) (entities.Service, error) {
	args := m.Called(id)
	return args.Get(0).(entities.Service), args.Error(1)
}

func (m *MockServiceService) FindServiceByActionId(id string) (entities.Service, error) {
	args := m.Called(id)
	return args.Get(0).(entities.Service), args.Error(1)
}

func (m *MockServiceService) FindServiceByReactionId(id string) (entities.Service, error) {
	args := m.Called(id)
	return args.Get(0).(entities.Service), args.Error(1)
}

func (m *MockServiceService) RetrieveActionsFromService(serviceName string) ([]entities.Action, error) {
	args := m.Called(serviceName)
	return args.Get(0).([]entities.Action), args.Error(1)
}

func (m *MockServiceService) RetrieveReactionsFromService(serviceName string) ([]entities.Reaction, error) {
	args := m.Called(serviceName)
	return args.Get(0).([]entities.Reaction), args.Error(1)
}

func (m *MockServiceService) GetGoogleRefreshTokenRequest(refreshToken string) (*http.Request, error) {
	return nil, nil
}

func (m *MockServiceService) GetSpotifyRefreshTokenRequest(refreshToken string) (*http.Request, error) {
	return nil, nil
}

func (m *MockServiceService) GetDiscordRefreshTokenRequest(refreshToken string) (*http.Request, error) {
	return nil, nil
}

func (m *MockServiceService) GetRedditRefreshTokenRequest(refreshToken string) (*http.Request, error) {
	return nil, nil
}

func (m *MockServiceService) GetAsanaRefreshTokenRequest(refreshToken string) (*http.Request, error) {
	return nil, nil
}

func (m *MockServiceService) GetDropboxRefreshTokenRequest(refreshToken string) (*http.Request, error) {
	return nil, nil
}

func (m *MockServiceService) GetGitlabRefreshTokenRequest(refreshToken string) (*http.Request, error) {
	return nil, nil
}

func (m *MockServiceService) ExecuteRequest(request *http.Request) (*http.Response, error) {
	return nil, nil
}

func (m *MockServiceService) ExecuteApiRequest(url, method, typeToken, accessToken string, body io.Reader) (*http.Response, error) {
	return nil, nil
}

func (m *MockServiceService) GetResultTokenFromCode(code, serviceName, callbackType, appType string) (entities.ResultToken, error) {
	var test entities.ResultToken
	return test, nil
}

func (m *MockServiceService) GetUserInfoFromService(accessToken, serviceName string) (entities.UserInfo, error) {
	var test entities.UserInfo
	return test, nil
}

func (m *MockServiceService) RequestToTimeApi() (entities.TimeResponse, error) {
	var test entities.TimeResponse
	return test, nil
}

func (m *MockServiceService) RequestGithubUserRepositories(accessToken string) ([]entities.GithubRepository, error) {
	return nil, nil
}

func (m *MockServiceService) RequestGitlabUserProjects(accessToken string) ([]entities.GitlabProject, error) {
	return nil, nil
}

func (m *MockServiceService) RetrieveDiscordGuildChannels(guildId string) ([]map[string]interface{}, error) {
	args := m.Called(guildId)
	return args.Get(0).([]map[string]interface{}), args.Error(1)
}

func TestOAuth2Service(test *testing.T) {
	mockService := new(MockServiceService)
	handler := &ServiceHandler{ServiceService: mockService}

	router := gin.Default()
	router.GET("/authentication/:service/:callback-type/:app-type", handler.oauth2Service)

	test.Run("Invalid app type", func(test *testing.T) {
		req, _ := http.NewRequest("GET", "/authentication/Github/service/false", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		require.Equal(test, http.StatusBadRequest, w.Code)
		require.JSONEq(test, `{"error": "Invalid app type"}`, w.Body.String())
	})

	test.Run("Invalid callback type", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/authentication/Github/false/web", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		require.Equal(test, http.StatusBadRequest, w.Code)
		require.JSONEq(test, `{"error": "Invalid callback type"}`, w.Body.String())
	})

	test.Run("Service error", func(test *testing.T) {
		mockService.On("OAuth2Service", "false", "service", "web").Return("", fmt.Errorf("service error"))

		req, _ := http.NewRequest("GET", "/authentication/false/service/web", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		require.Equal(test, http.StatusBadRequest, w.Code)
		require.JSONEq(test, `{"error": "Unknown service"}`, w.Body.String())

		mockService.AssertCalled(test, "OAuth2Service", "false", "service", "web")
	})

	test.Run("Successful", func(test *testing.T) {
		mockService.On("OAuth2Service", "Github", "service", "web").Return("https://localhost:8080/authentication/Github/service/web", nil)

		req, _ := http.NewRequest("GET", "/authentication/Github/service/web", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		require.Equal(test, http.StatusOK, w.Code)
		require.JSONEq(test, `{"auth-url": "https://localhost:8080/authentication/Github/service/web"}`, w.Body.String())

		mockService.AssertCalled(test, "OAuth2Service", "Github", "service", "web")
	})
}
