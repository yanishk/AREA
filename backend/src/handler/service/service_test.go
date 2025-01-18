package service_handler

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"backend/src/entities"
	"backend/src/handler/middleware"
)

type MockServiceService struct {
	mock.Mock
}

func (m *MockServiceService) OAuth2Service(serviceName, callbackType, appType string) (string, error) {
	args := m.Called(serviceName, callbackType, appType)
	return args.String(0), args.Error(1)
}

func (m *MockServiceService) RetrieveActionsServices() ([]entities.Service, error) {
	args := m.Called()
	return args.Get(0).([]entities.Service), args.Error(1)
}

func (m *MockServiceService) RetrieveReactionsServices() ([]entities.Service, error) {
	args := m.Called()
	return args.Get(0).([]entities.Service), args.Error(1)
}

func (m *MockServiceService) FindServiceByName(name string) (entities.Service, error) {
	args := m.Called(name)
	return args.Get(0).(entities.Service), args.Error(1)
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

func (m *MockServiceService) FindAllServices() ([]entities.Service, error) {
	args := m.Called()
	return args.Get(0).([]entities.Service), args.Error(1)
}

func requestForProtected(method, url, token string) *http.Request {
	req, _ := http.NewRequest(method, url, nil)
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

func createMockAndRoute(isProtected bool) (*ServiceHandler, *gin.Engine, *MockServiceService) {
	mockService := new(MockServiceService)
	handler := &ServiceHandler{ServiceService: mockService}

	router := gin.Default()
	if isProtected {
		router.Use(middleware.VerifyJWTCookie, middleware.VerifyEmailFromContext, middleware.VerifyConnectionTypeFromContext)
	}

	return handler, router, mockService
}

func TestOAuth2Service(test *testing.T) {
	handler, router, mockService := createMockAndRoute(false)
	router.GET("/authentication", handler.oauth2Service)

	test.Run("Successful", func(test *testing.T) {
		mockService.On("OAuth2Service", "Github", "service", "web").
			Return("https://localhost:8080/authentication/Github/service/web", nil)

		req, _ := http.NewRequest("GET", "/authentication?service=Github&callbacktype=service&apptype=web", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		require.Equal(test, http.StatusOK, w.Code)

		mockService.AssertCalled(test, "OAuth2Service", "Github", "service", "web")
	})

	test.Run("Invalid app type", func(test *testing.T) {
		req, _ := http.NewRequest("GET", "/authentication?service=Github&callbacktype=service&apptype=false", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		require.Equal(test, http.StatusBadRequest, w.Code)
		require.JSONEq(test, `{"error": "Invalid app type"}`, w.Body.String())
	})

	test.Run("Invalid callback type", func(test *testing.T) {
		req, _ := http.NewRequest("GET", "/authentication?service=Github&callbacktype=false&apptype=web", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		require.Equal(test, http.StatusBadRequest, w.Code)
		require.JSONEq(test, `{"error": "Invalid callback type"}`, w.Body.String())
	})

	test.Run("Service error", func(test *testing.T) {
		mockService.On("OAuth2Service", "false", "service", "web").
			Return("", fmt.Errorf("service error"))

		req, _ := http.NewRequest("GET", "/authentication?service=false&callbacktype=service&apptype=web", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		require.Equal(test, http.StatusBadRequest, w.Code)
		require.JSONEq(test, `{"error": "Unknown service"}`, w.Body.String())

		mockService.AssertCalled(test, "OAuth2Service", "false", "service", "web")
	})
}

func TestRetrieveAllService(test *testing.T) {
	test.Run("Successful", func(test *testing.T) {
		handler, router, mockService := createMockAndRoute(false)

		router.GET("/services", handler.retrieveAllServices)

		mockService.On("FindAllServices").
			Return([]entities.Service{}, nil)

		req, _ := http.NewRequest("GET", "/services", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		require.Equal(test, http.StatusOK, w.Code)
	})

	test.Run("Fail", func(test *testing.T) {
		handler, router, mockService := createMockAndRoute(false)

		router.GET("/services", handler.retrieveAllServices)

		mockService.On("FindAllServices").
			Return([]entities.Service{}, errors.New("Fail"))

		req, _ := http.NewRequest("GET", "/services", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		require.Equal(test, http.StatusInternalServerError, w.Code)
		require.JSONEq(test, `{"error": "Internal server error"}`, w.Body.String())
	})
}

func TestRetrieveActionsServices(test *testing.T) {
	test.Run("Successful", func(test *testing.T) {
		handler, router, mockService := createMockAndRoute(true)

		token := createToken(test)

		router.GET("/services/actions", handler.retrieveActionsServices)

		mockService.On("RetrieveActionsServices").
			Return([]entities.Service{}, nil)

		req := requestForProtected("GET", "/services/actions", token)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		require.Equal(test, http.StatusOK, w.Code)

	})

	test.Run("Fail", func(test *testing.T) {
		handler, router, mockService := createMockAndRoute(true)

		token := createToken(test)

		router.GET("/services/actions", handler.retrieveActionsServices)

		mockService.On("RetrieveActionsServices").
			Return([]entities.Service{}, errors.New("Fail"))

		req := requestForProtected("GET", "/services/actions", token)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		require.Equal(test, http.StatusInternalServerError, w.Code)
		require.JSONEq(test, `{"error": "Internal server error"}`, w.Body.String())
	})
}

func TestRetrieveReactionsServices(test *testing.T) {
	test.Run("Successful", func(test *testing.T) {
		handler, router, mockService := createMockAndRoute(true)

		token := createToken(test)

		router.GET("/services/reactions", handler.retrieveReactionsServices)

		mockService.On("RetrieveReactionsServices").
			Return([]entities.Service{}, nil)

		req := requestForProtected("GET", "/services/reactions", token)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		require.Equal(test, http.StatusOK, w.Code)

	})

	test.Run("Fail", func(test *testing.T) {
		handler, router, mockService := createMockAndRoute(true)

		token := createToken(test)

		router.GET("/services/reactions", handler.retrieveReactionsServices)

		mockService.On("RetrieveReactionsServices").
			Return([]entities.Service{}, errors.New("Fail"))

		req := requestForProtected("GET", "/services/reactions", token)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		require.Equal(test, http.StatusInternalServerError, w.Code)
		require.JSONEq(test, `{"error": "Internal server error"}`, w.Body.String())
	})
}

func TestRetrieveServiceById(test *testing.T) {
	test.Run("Successful", func(test *testing.T) {
		handler, router, mockService := createMockAndRoute(true)

		token := createToken(test)
		expectedService := entities.Service{Id: "1", Name: "service"}

		router.GET("/services/:id", handler.retrieveServiceById)

		mockService.On("FindServiceById", "1").
			Return(expectedService, nil)

		req := requestForProtected("GET", "/services/1", token)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		require.Equal(test, http.StatusOK, w.Code)

		mockService.AssertExpectations(test)
	})

	test.Run("Fail", func(test *testing.T) {
		handler, router, mockService := createMockAndRoute(true)

		token := createToken(test)
		expectedService := entities.Service{Id: "1", Name: "service"}

		router.GET("/services/:id", handler.retrieveServiceById)

		mockService.On("FindServiceById", "1").
			Return(expectedService, errors.New("Fail"))

		req := requestForProtected("GET", "/services/1", token)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		require.Equal(test, http.StatusInternalServerError, w.Code)
		require.JSONEq(test, `{"error": "Could not retrieve requested service"}`, w.Body.String())

		mockService.AssertExpectations(test)
	})
}

func TestRetrieveServiceByActionId(test *testing.T) {
	test.Run("Successful", func(test *testing.T) {
		handler, router, mockService := createMockAndRoute(true)

		token := createToken(test)

		router.GET("/services/action/:actionid", handler.retrieveServiceByActionId)

		mockService.On("FindServiceByActionId", "1").
			Return(entities.Service{}, nil)

		req := requestForProtected("GET", "/services/action/1", token)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		require.Equal(test, http.StatusOK, w.Code)

		mockService.AssertExpectations(test)
	})

	test.Run("Fail", func(test *testing.T) {
		handler, router, mockService := createMockAndRoute(true)

		token := createToken(test)

		router.GET("/services/action/:actionid", handler.retrieveServiceByActionId)

		mockService.On("FindServiceByActionId", "1").
			Return(entities.Service{}, errors.New("Fail"))

		req := requestForProtected("GET", "/services/action/1", token)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		require.Equal(test, http.StatusInternalServerError, w.Code)
		require.JSONEq(test, `{"error": "Could not retrieve requested service"}`, w.Body.String())

		mockService.AssertExpectations(test)
	})
}

func TestRetrieveServiceByReactionId(test *testing.T) {
	test.Run("Successful", func(test *testing.T) {
		handler, router, mockService := createMockAndRoute(true)

		token := createToken(test)

		router.GET("/services/reaction/:reactionid", handler.retrieveServiceByReactionId)

		mockService.On("FindServiceByReactionId", "1").
			Return(entities.Service{}, nil)

		req := requestForProtected("GET", "/services/reaction/1", token)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		require.Equal(test, http.StatusOK, w.Code)

		mockService.AssertExpectations(test)
	})

	test.Run("Fail", func(test *testing.T) {
		handler, router, mockService := createMockAndRoute(true)

		token := createToken(test)

		router.GET("/services/reaction/:reactionid", handler.retrieveServiceByReactionId)

		mockService.On("FindServiceByReactionId", "1").
			Return(entities.Service{}, errors.New("Fail"))

		req := requestForProtected("GET", "/services/reaction/1", token)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		require.Equal(test, http.StatusInternalServerError, w.Code)
		require.JSONEq(test, `{"error": "Could not retrieve requested service"}`, w.Body.String())

		mockService.AssertExpectations(test)
	})
}

func TestRetrieveActionsFromService(test *testing.T) {
	handler, router, mockService := createMockAndRoute(false)

	expectedActions := []entities.Action{
		{Id: "1", Name: "action1"},
		{Id: "2", Name: "action2"},
	}

	router.GET("/service/:service/actions", handler.retrieveActionsFromService)

	test.Run("Successful", func(test *testing.T) {
		mockService.On("RetrieveActionsFromService", "Reddit").Return(expectedActions, nil)

		req, _ := http.NewRequest("GET", "/service/Reddit/actions", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		require.Equal(test, http.StatusOK, w.Code)

		mockService.AssertExpectations(test)
	})

	test.Run("Unknown service", func(test *testing.T) {
		var action []entities.Action

		mockService.On("RetrieveActionsFromService", "false").Return(action, fmt.Errorf("service error"))

		req, _ := http.NewRequest("GET", "/service/false/actions", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		require.Equal(test, http.StatusBadRequest, w.Code)
		require.JSONEq(test, `{"error": "Unknown service"}`, w.Body.String())

		mockService.AssertExpectations(test)
	})
}

func TestRetrieveReactionsFromService(test *testing.T) {
	handler, router, mockService := createMockAndRoute(false)

	expectedReactions := []entities.Reaction{
		{Id: "1", Name: "reaction1"},
		{Id: "2", Name: "reaction2"},
	}

	router.GET("/service/:service/reactions", handler.retrieveReactionsFromService)

	test.Run("Successful", func(test *testing.T) {
		mockService.On("RetrieveReactionsFromService", "Reddit").Return(expectedReactions, nil)

		req, _ := http.NewRequest("GET", "/service/Reddit/reactions", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		require.Equal(test, http.StatusOK, w.Code)

		mockService.AssertExpectations(test)
	})

	test.Run("Unknown service", func(test *testing.T) {
		var reactions []entities.Reaction

		mockService.On("RetrieveReactionsFromService", "false").Return(reactions, fmt.Errorf("service error"))

		req, _ := http.NewRequest("GET", "/service/false/reactions", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		require.Equal(test, http.StatusBadRequest, w.Code)
		require.JSONEq(test, `{"error": "Unknown service"}`, w.Body.String())

		mockService.AssertExpectations(test)
	})
}

func TestGetDiscordGuildChannels(test *testing.T) {

	handler, router, mockService := createMockAndRoute(true)

	token := createToken(test)

	guildId := "12345"
	expectedChannels := []map[string]interface{}{
		{"name": "test1", "id": "1"},
		{"name": "test2", "id": "2"},
		{"name": "test3", "id": "3"},
	}

	router.GET("/discord/server/channels", handler.getDiscordGuildChannels)

	createRequestWithJWT := func(method, url, guildId string) *http.Request {
		req, _ := http.NewRequest(method, url+"?id="+guildId, nil)
		req.AddCookie(&http.Cookie{Name: "JWToken", Value: token})
		return req
	}

	test.Run("Successful", func(t *testing.T) {
		mockService.On("RetrieveDiscordGuildChannels", guildId).
			Return(expectedChannels, nil)

		req := createRequestWithJWT("GET", "/discord/server/channels", guildId)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)

		mockService.AssertExpectations(t)
	})

	test.Run("Fail", func(t *testing.T) {
		handler, router, mockService := createMockAndRoute(true)

		token := createToken(test)

		guildId := "12345"
		expectedChannels := []map[string]interface{}{
			{"name": "test1", "id": "1"},
			{"name": "test2", "id": "2"},
			{"name": "test3", "id": "3"},
		}

		router.GET("/discord/server/channels", handler.getDiscordGuildChannels)

		createRequestWithJWT := func(method, url, guildId string) *http.Request {
			req, _ := http.NewRequest(method, url+"?id="+guildId, nil)
			req.AddCookie(&http.Cookie{Name: "JWToken", Value: token})
			return req
		}

		mockService.On("RetrieveDiscordGuildChannels", guildId).
			Return(expectedChannels, errors.New("Fail"))

		req := createRequestWithJWT("GET", "/discord/server/channels", guildId)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		require.Equal(t, http.StatusInternalServerError, w.Code)
		require.JSONEq(test, `{"error": "Could not retrieve the server's channels"}`, w.Body.String())

		mockService.AssertExpectations(t)
	})
}
