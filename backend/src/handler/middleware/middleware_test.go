package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

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
