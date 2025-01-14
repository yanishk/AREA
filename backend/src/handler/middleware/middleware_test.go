package middleware

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/require"
)

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

func TestVerifyConnectionTypeFromContext(test *testing.T) {
	test.Run("Successful", func(test *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Set("connectionType", "basic")

		VerifyConnectionTypeFromContext(c)

		require.Equal(test, http.StatusOK, w.Code)
	})

	test.Run("Connection type not found", func(test *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		VerifyConnectionTypeFromContext(c)

		require.Equal(test, http.StatusUnauthorized, w.Code)
		require.JSONEq(test, `{"error": "Connection type not found in token"}`, w.Body.String())
	})

	test.Run("Connection Invalid String", func(test *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Set("connectionType", 1)

		VerifyConnectionTypeFromContext(c)

		require.Equal(test, http.StatusUnauthorized, w.Code)
		require.JSONEq(test, `{"error": "Connection type is not a valid string"}`, w.Body.String())
	})
}

func TestVerifyEmailFromContext(test *testing.T) {
	test.Run("Successful", func(test *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Set("email", "test@test.com")

		VerifyEmailFromContext(c)

		require.Equal(test, http.StatusOK, w.Code)
	})

	test.Run("Email not found", func(test *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		VerifyEmailFromContext(c)

		require.Equal(test, http.StatusUnauthorized, w.Code)
		require.JSONEq(test, `{"error": "Email not found in token"}`, w.Body.String())
	})

	test.Run("Connection Invalid String", func(test *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Set("email", 1)

		VerifyEmailFromContext(c)

		require.Equal(test, http.StatusUnauthorized, w.Code)
		require.JSONEq(test, `{"error": "Email is not a valid string"}`, w.Body.String())
	})
}

func TestVerifyJWTCookie(test *testing.T) {
	test.Run("Successful", func(test *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		validToken := createToken(test)
		c.Request = &http.Request{
			Header: http.Header{
				"Cookie": {
					"JWToken=" + validToken,
				},
			},
		}

		VerifyJWTCookie(c)

		require.Equal(test, http.StatusOK, w.Code)

		_, exists := c.Get("email")
		require.True(test, exists)

		_, exists = c.Get("connectionType")
		require.True(test, exists)
	})

	test.Run("No Token", func(test *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request = &http.Request{
			Header: http.Header{},
		}

		VerifyJWTCookie(c)

		require.Equal(test, http.StatusUnauthorized, w.Code)
		require.JSONEq(test, `{"error": "No authentication token"}`, w.Body.String())
	})

	test.Run("Invalid Token", func(test *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request = &http.Request{Header: http.Header{"Cookie": {"JWToken=invalids"}}}

		VerifyJWTCookie(c)

		require.Equal(test, http.StatusUnauthorized, w.Code)
		require.JSONEq(test, `{"error": "Invalid token"}`, w.Body.String())
	})
}
