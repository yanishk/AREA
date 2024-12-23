package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func VerifyEmailFromContext(context *gin.Context) {
	email, ok := context.Get("email")
	if !ok {
		context.IndentedJSON(http.StatusUnauthorized, gin.H{
			"error": "Email not found in token",
		})
		context.Abort()
		return
	}

	_, ok = email.(string)
	if !ok {
		context.IndentedJSON(http.StatusUnauthorized, gin.H{
			"error": "Email is not a valid string",
		})
		context.Abort()
		return
	}
	context.Next()
}
