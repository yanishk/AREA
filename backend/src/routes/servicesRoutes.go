package routes

import (
	"github.com/gin-gonic/gin"
	"poc-goland.com/poc/src/entities"
)

func PrivateServicesRoutes(router *gin.Engine) {
	router.GET("/auth/spotify-auth", entities.Auth02Spotify)
	router.GET("/spotify-callback", entities.SpotifyCallBack)
}
