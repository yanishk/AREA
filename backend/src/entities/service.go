package entities

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func Auth02Spotify(context *gin.Context) {
	spotifyAuthUrl := "https://accounts.spotify.com/authorize?client_id=" +
		os.Getenv("SPOTIFY_ID_CLIENT") + "&response_type=code&redirect_uri=" +
		os.Getenv("CALLBACK_SPOTIFY") + "&scope=user-read-private user-read-email"
	context.Redirect(http.StatusTemporaryRedirect, spotifyAuthUrl)
}

func SpotifyCallBack(context *gin.Context) {
	codeSpotify := context.Query("code")

	if codeSpotify == "" {
		context.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid code authorization",
		})
		return
	}
	tokenSpotify, errToken := getAccesTokenFromCode(codeSpotify)

	if errToken != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to exchange token",
			"details": errToken.Error(),
		})
		return
	}

	context.IndentedJSON(http.StatusOK, gin.H{
		"access_token": tokenSpotify,
	})
}
