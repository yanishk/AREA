package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func VerifyConnectionTypeFromContext(context *gin.Context) {
	connectionType, errGetConnectionType := context.Get("connectionType")
	if !errGetConnectionType {
		context.IndentedJSON(http.StatusUnauthorized, gin.H{
			"error": "Connection type not found in token",
		})
		context.Abort()
		return
	}

	_, errGetConnectionTypeStr := connectionType.(string)
	if !errGetConnectionTypeStr {
		context.IndentedJSON(http.StatusUnauthorized, gin.H{
			"error": "Connection type is not a valid string",
		})
		context.Abort()
		return
	}
	context.Next()
}
