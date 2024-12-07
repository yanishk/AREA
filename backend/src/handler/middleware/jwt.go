package middleware

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func VerifyJWTCookie(context *gin.Context) {
	cookie, err := context.Cookie("JWToken")

	if err != nil {
		context.IndentedJSON(http.StatusUnauthorized, gin.H{
			"error": "No Authentication Token",
		})
		context.Abort()
		return
	}

	token, errParse := jwt.Parse(cookie, func(token *jwt.Token) (interface{}, error) {
		_, isCorrectMethod := token.Method.(*jwt.SigningMethodHMAC)
		if !isCorrectMethod {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("SECRET_KEY")), nil
	})

	if errParse != nil {
		context.IndentedJSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid Token",
		})
		context.Abort()
		return
	}

	claims, errorClaim := token.Claims.(jwt.MapClaims)
	if !errorClaim || !token.Valid {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
		context.Abort()
		return
	}

	context.Set("email", claims["email"])
	context.Next()
}
