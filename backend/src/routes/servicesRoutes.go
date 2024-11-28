package routes

import (
	"github.com/gin-gonic/gin"
	"poc-goland.com/poc/src/entities"
)

func PrivateServicesRoutes(router *gin.Engine, databaseContext entities.Database) {
	router.GET("/auth/spotify-auth", entities.Auth02Spotify)
	router.GET("/spotify-callback", func(context *gin.Context) {
		databaseContext.SpotifyCallBack(context)
	})
	router.GET("spotify/call-api", func(context *gin.Context) {
		databaseContext.SpotifyCallApi(context)
	})
}
