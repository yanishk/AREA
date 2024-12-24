package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func VerifyConnectionTypeFromContext(context *gin.Context) {
	connectionType, ok := context.Get("connectionType")
	if !ok {
		context.IndentedJSON(http.StatusUnauthorized, gin.H{
			"error": "Connection type not found in token",
		})
		context.Abort()
		return
	}

	_, ok = connectionType.(string)
	if !ok {
		context.IndentedJSON(http.StatusUnauthorized, gin.H{
			"error": "Connection type is not a valid string",
		})
		context.Abort()
		return
	}
	context.Next()
}
