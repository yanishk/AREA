package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func VerifyEmailFromContext(context *gin.Context) {
	email, errGetEmail := context.Get("email")

	if !errGetEmail {
		context.IndentedJSON(http.StatusUnauthorized, gin.H{
			"error": "Email not found in token",
		})
		context.Abort()
		return
	}

	_, errGetEmailStr := email.(string)
	if !errGetEmailStr {
		context.IndentedJSON(http.StatusUnauthorized, gin.H{
			"error": "Email is not a valid string",
		})
		context.Abort()
		return
	}
	context.Next()
}
