package user_handler

import (
	"backend/src/entities"
	"backend/src/handler/middleware"
	"io"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
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
		"email":          "eva.legrand@epitech.eu",
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
