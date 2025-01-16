package user_handler

import (
	"errors"
	"fmt"
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

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) CreateUser(userEmail, userPassword, userConnectionType string) error {
	args := m.Called(userEmail, userPassword, userConnectionType)
	return args.Error(0)
}

func (m *MockUserService) LoginAuthentication(userEmail, userPassword, userConnectionType string) (string, error) {
	args := m.Called(userEmail, userPassword, userConnectionType)
	return args.String(0), args.Error(1)
}

func (m *MockUserService) LoginWithService(code, serviceName, appType string) (string, error) {
	args := m.Called(code, serviceName, appType)
	return args.String(0), args.Error(1)
}

func (m *MockUserService) GetUser(userEmail, userConnectionType string) (entities.UserInfos, error) {
	args := m.Called(userEmail, userConnectionType)
	return args.Get(0).(entities.UserInfos), args.Error(1)
}

func (m *MockUserService) ModifyPassword(userEmail, userConnectionType string, newPassword entities.UserModifyPassword) error {
	args := m.Called(userEmail, userConnectionType, newPassword)
	return args.Error(0)
}

func (m *MockUserService) DeleteAccount(userEmail, userConnectionType string) error {
	args := m.Called(userEmail, userConnectionType)
	return args.Error(0)
}

func (m *MockUserService) FindUserById(userId string) (entities.User, error) {
	args := m.Called(userId)
	return args.Get(0).(entities.User), args.Error(1)
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

func createMockAndRoute(isProtected bool) (*UserHandler, *gin.Engine, *MockUserService) {
	mockUserService := new(MockUserService)
	handler := &UserHandler{UserService: mockUserService}

	router := gin.Default()
	if isProtected {
		router.Use(middleware.VerifyJWTCookie, middleware.VerifyEmailFromContext, middleware.VerifyConnectionTypeFromContext)
	}

	return handler, router, mockUserService
}

func TestCreateUser(test *testing.T) {
	handler, router, mockUserService := createMockAndRoute(false)
	router.POST("/register", handler.createUser)

	test.Run("Successful", func(test *testing.T) {
		mockUserService.On("CreateUser", "test@test.com", "password", "basic").
			Return(nil)

		body := `{
			"email": "test@test.com",
			"password": "password"
		}`

		req, _ := http.NewRequest("POST", "/register", strings.NewReader(body))
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		require.Equal(test, http.StatusOK, w.Code)
		require.JSONEq(test, `{"success": "New user created"}`, w.Body.String())
	})

	test.Run("Fail JSON Bind", func(test *testing.T) {
		mockUserService.On("CreateUser", "test@test.com", "password", "basic").
			Return(nil)

		req, _ := http.NewRequest("POST", "/register", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		require.Equal(test, http.StatusBadRequest, w.Code)
		require.JSONEq(test, `{"error": "Invalid request body"}`, w.Body.String())
	})

	test.Run("Fail connection type", func(test *testing.T) {
		handler, router, mockUserService := createMockAndRoute(false)
		router.POST("/register", handler.createUser)

		mockUserService.On("CreateUser", "test@test.com", "password", "basic").
			Return(errors.New("Connection type doesn't exist"))

		body := `{
			"email": "test@test.com",
			"password": "password"
		}`

		req, _ := http.NewRequest("POST", "/register", strings.NewReader(body))
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		require.Equal(test, http.StatusBadRequest, w.Code)
		require.JSONEq(test, `{"error": "Connection type doesn't exist"}`, w.Body.String())
	})

	test.Run("Email already used", func(test *testing.T) {
		handler, router, mockUserService := createMockAndRoute(false)
		router.POST("/register", handler.createUser)

		mockUserService.On("CreateUser", "test@test.com", "password", "basic").
			Return(errors.New("Email address already used"))

		body := `{
			"email": "test@test.com",
			"password": "password"
		}`

		req, _ := http.NewRequest("POST", "/register", strings.NewReader(body))
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		require.Equal(test, http.StatusConflict, w.Code)
		require.JSONEq(test, `{"error": "Email address already used"}`, w.Body.String())
	})

	test.Run("Other error", func(test *testing.T) {
		handler, router, mockUserService := createMockAndRoute(false)
		router.POST("/register", handler.createUser)

		mockUserService.On("CreateUser", "test@test.com", "password", "basic").
			Return(errors.New("other error"))

		body := `{
			"email": "test@test.com",
			"password": "password"
		}`

		req, _ := http.NewRequest("POST", "/register", strings.NewReader(body))
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		require.Equal(test, http.StatusInternalServerError, w.Code)
		require.JSONEq(test, `{"error": "other error"}`, w.Body.String())
	})
}

func TestLoginAuthentication(test *testing.T) {
	handler, router, mockUserService := createMockAndRoute(false)
	router.POST("/login", handler.loginAuthentication)

	test.Run("Successful", func(test *testing.T) {
		mockUserService.On("LoginAuthentication", "test@test.com", "password", "basic").
			Return("token", nil).Once()

		body := `{
			"email": "test@test.com",
			"password": "password"
		}`

		req, _ := http.NewRequest("POST", "/login", strings.NewReader(body))
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		require.Equal(test, http.StatusOK, w.Code)
		require.JSONEq(test, `{"success": "Connection successful"}`, w.Body.String())
	})

	test.Run("User not found", func(test *testing.T) {
		mockUserService.On("LoginAuthentication", "test@test.com", "password", "basic").
			Return("", fmt.Errorf("Could not find requested user")).Once()

		body := `{
			"email": "test@test.com",
			"password": "password"
		}`

		req, _ := http.NewRequest("POST", "/login", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		require.Equal(test, http.StatusUnauthorized, w.Code)
		require.JSONEq(test, `{"error": "Could not find requested user"}`, w.Body.String())
	})

	test.Run("Other error", func(test *testing.T) {
		mockUserService.On("LoginAuthentication", "test@test.com", "password", "basic").
			Return("token", errors.New("Other errors")).Once()

		body := `{
			"email": "test@test.com",
			"password": "password"
		}`

		req, _ := http.NewRequest("POST", "/login", strings.NewReader(body))
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		require.Equal(test, http.StatusInternalServerError, w.Code)
		require.JSONEq(test, `{"error": "Other errors"}`, w.Body.String())
	})

	test.Run("Fail JSON Bind", func(test *testing.T) {
		mockUserService.On("LoginAuthentication", "test@test.com", "password", "basic").
			Return("token", nil).Once()

		req, _ := http.NewRequest("POST", "/login", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		require.Equal(test, http.StatusBadRequest, w.Code)
		require.JSONEq(test, `{"error": "Invalid request body"}`, w.Body.String())
	})
}

func TestLoginCallback(test *testing.T) {
	handler, router, mockUserService := createMockAndRoute(false)
	router.POST("/login-callback", handler.loginCallback)

	test.Run("Successful", func(test *testing.T) {
		mockUserService.On("LoginWithService", "code", "github", "web").
			Return("token", nil)

		body := `{
			"service": "github",
			"apptype": "web"
		}`

		req, _ := http.NewRequest("POST", "/login-callback?code=code", strings.NewReader(body))
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		require.Equal(test, http.StatusOK, w.Code)
		require.JSONEq(test, `{"success": "Connection successful"}`, w.Body.String())
	})

	test.Run("Invalid code authorization", func(t *testing.T) {
		body := `{
			"service": "Github",
			"apptype": "web"
		}`

		req, _ := http.NewRequest("POST", "/login-callback", strings.NewReader(body))
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		require.Equal(test, http.StatusBadRequest, w.Code)
		require.JSONEq(test, `{"error": "Invalid code authorization"}`, w.Body.String())
	})

	test.Run("Invalid app type", func(test *testing.T) {
		body := `{
			"service": "Github",
			"apptype": "unknown"
		}`

		req, _ := http.NewRequest("POST", "/login-callback?code=valid-code", strings.NewReader(body))
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		require.Equal(test, http.StatusBadRequest, w.Code)
		require.JSONEq(test, `{"error": "Invalid app type"}`, w.Body.String())
	})

	test.Run("Failed to connect", func(test *testing.T) {
		mockUserService.On("LoginWithService", "code", "Github", "web").
			Return("", fmt.Errorf("service error"))

		body := `{
			"service": "Github",
			"apptype": "web"
		}`

		req, _ := http.NewRequest("POST", "/login-callback?code=code", strings.NewReader(body))
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		require.Equal(test, http.StatusInternalServerError, w.Code)
		require.JSONEq(test, `{"error": "Failed to connect with requested service"}`, w.Body.String())
	})

	test.Run("Fail JSON Bind", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/login-callback", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		require.Equal(test, http.StatusBadRequest, w.Code)
		require.JSONEq(test, `{"error": "Invalid request body"}`, w.Body.String())
	})
}

func TestGetUser(test *testing.T) {
	handler, router, mockUserService := createMockAndRoute(false)

	token := createToken(test)

	router.Use(func(c *gin.Context) {
		c.Set("email", "email")
		c.Set("connectionType", "service")
	})
	router.GET("/user", handler.getUser)

	test.Run("Successful", func(test *testing.T) {
		mockUserService.On("GetUser", "email", "service").
			Return(entities.UserInfos{}, nil)

		req := requestForProtected("GET", "/user", token, nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		require.Equal(test, http.StatusOK, w.Code)
	})

	test.Run("Error Finding User", func(test *testing.T) {
		mockUserService.ExpectedCalls = nil
		mockUserService.On("GetUser", "email", "service").
			Return(entities.UserInfos{}, errors.New("user not found")).Once()

		req := requestForProtected("GET", "/user", token, nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		require.JSONEq(test, `{"error": "Could not find user"}`, w.Body.String())
	})
}

func TestLogoutUser(test *testing.T) {
	handler, router, mockUserService := createMockAndRoute(false)

	token := createToken(test)

	router.Use(func(c *gin.Context) {
		c.Set("email", "email")
		c.Set("connectionType", "service")
	})
	router.POST("/logout", handler.logoutUser)

	test.Run("Successful", func(test *testing.T) {
		mockUserService.On("GetUser", "email", "service").
			Return(entities.UserInfos{}, nil)

		req := requestForProtected("POST", "/logout", token, nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		require.Equal(test, http.StatusOK, w.Code)
		require.JSONEq(test, `{"success": "Logout successful"}`, w.Body.String())
	})

	test.Run("Error Finding User", func(test *testing.T) {
		mockUserService.ExpectedCalls = nil
		mockUserService.On("GetUser", "email", "service").
			Return(entities.UserInfos{}, errors.New("user not found")).Once()

		req := requestForProtected("POST", "/logout", token, nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		require.JSONEq(test, `{"error": "Could not find user"}`, w.Body.String())
	})
}

func TestModifyPassword(test *testing.T) {
	handler, router, mockUserService := createMockAndRoute(false)

	token := createToken(test)

	router.Use(func(c *gin.Context) {
		c.Set("email", "email")
		c.Set("connectionType", "service")
	})
	router.PUT("/modify-password", handler.modifyPassword)

	expectedPasswordChange := entities.UserModifyPassword{
		OldPassword: "oldpass",
		Password:    "newpass",
	}

	test.Run("Successful", func(test *testing.T) {
		mockUserService.On("ModifyPassword", "email", "service", expectedPasswordChange).Return(nil).Once()

		body := `{
			"oldpassword": "oldpass",
			"password": "newpass"
		}`

		req := requestForProtected("PUT", "/modify-password", token, strings.NewReader(body))
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		require.Equal(test, http.StatusOK, w.Code)
		require.JSONEq(test, `{"success": "Password modified"}`, w.Body.String())
	})

	test.Run("User Not Found", func(test *testing.T) {
		mockUserService.On("ModifyPassword", "email", "service", expectedPasswordChange).
			Return(errors.New("Could not find requested user")).Once()

		body := `{
			"oldpassword": "oldpass",
			"password": "newpass"
		}`

		req := requestForProtected("PUT", "/modify-password", token, strings.NewReader(body))
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		require.Equal(test, http.StatusBadRequest, w.Code)
		require.JSONEq(test, `{"error": "Could not find requested user"}`, w.Body.String())
	})

	test.Run("Old Password Incorrect", func(test *testing.T) {
		mockUserService.On("ModifyPassword", "email", "service", expectedPasswordChange).
			Return(errors.New("Old password is incorrect")).Once()

		body := `{
			"oldpassword": "oldpass",
			"password": "newpass"
		}`

		req := requestForProtected("PUT", "/modify-password", token, strings.NewReader(body))
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		require.Equal(test, http.StatusForbidden, w.Code)
		require.JSONEq(test, `{"error": "Old password is incorrect"}`, w.Body.String())
	})

	test.Run("Other errors", func(test *testing.T) {
		mockUserService.On("ModifyPassword", "email", "service", expectedPasswordChange).
			Return(errors.New("Other errors")).Once()

		body := `{
			"oldpassword": "oldpass",
			"password": "newpass"
		}`

		req := requestForProtected("PUT", "/modify-password", token, strings.NewReader(body))
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		require.Equal(test, http.StatusInternalServerError, w.Code)
		require.JSONEq(test, `{"error": "Other errors"}`, w.Body.String())
	})

	test.Run("Fail JSON Bind", func(test *testing.T) {
		req := requestForProtected("PUT", "/modify-password", token, nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		require.Equal(test, http.StatusBadRequest, w.Code)
		require.JSONEq(test, `{"error": "Invalid request body"}`, w.Body.String())
	})
}

func TestDeleteAccount(test *testing.T) {
	handler, router, mockUserService := createMockAndRoute(false)

	token := createToken(test)

	router.Use(func(c *gin.Context) {
		c.Set("email", "email")
		c.Set("connectionType", "service")
	})
	router.DELETE("/user", handler.deleteAccount)

	test.Run("Successful", func(test *testing.T) {
		mockUserService.On("DeleteAccount", "email", "service").Return(nil).Once()

		req := requestForProtected("DELETE", "/user", token, nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		require.Equal(test, http.StatusOK, w.Code)
		require.JSONEq(test, `{"success": "Account deleted"}`, w.Body.String())
	})

	test.Run("Error Deleting Account", func(test *testing.T) {
		mockUserService.On("DeleteAccount", "email", "service").Return(errors.New("account not deleted")).Once()

		req := requestForProtected("DELETE", "/user", token, nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		require.JSONEq(test, `{"error": "Could not delete account"}`, w.Body.String())
	})
}
