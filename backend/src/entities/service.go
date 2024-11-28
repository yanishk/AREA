package entities

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"poc-goland.com/poc/src/database"
)

type ResultToken struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

func Auth02Spotify(context *gin.Context) {
	spotifyAuthUrl := "https://accounts.spotify.com/authorize?client_id=" +
		os.Getenv("SPOTIFY_ID_CLIENT") + "&response_type=code&redirect_uri=" +
		os.Getenv("CALLBACK_SPOTIFY") + "&scope=user-read-private user-read-email"
	context.Redirect(http.StatusTemporaryRedirect, spotifyAuthUrl)
}

func (databaseContext *Database) SpotifyCallBack(context *gin.Context) {
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

	_, errFindingService := database.FindServiceByName(databaseContext.Db, "spotify")
	if errFindingService == nil {
		database.UpdateService(databaseContext.Db, "spotify", tokenSpotify.AccessToken, tokenSpotify.RefreshToken)
	} else {
		database.CreateService(databaseContext.Db, "spotify", tokenSpotify.AccessToken, tokenSpotify.RefreshToken)
	}

	context.IndentedJSON(http.StatusOK, gin.H{
		"access_token": tokenSpotify,
	})
}

func getAccesTokenFromCode(codeSpotify string) (ResultToken, error) {
	tokenUrl := "https://accounts.spotify.com/api/token"
	jsonBody := fmt.Sprintf(
		"grant_type=authorization_code&code=%s&redirect_uri=%s",
		codeSpotify,
		os.Getenv("CALLBACK_SPOTIFY"),
	)
	var tokensRes ResultToken

	request, errRequest := http.NewRequest("POST", tokenUrl, bytes.NewBuffer([]byte(jsonBody)))
	if errRequest != nil {
		return tokensRes, errRequest
	}

	request.SetBasicAuth(os.Getenv("SPOTIFY_ID_CLIENT"), os.Getenv("SPOTIFY_CLIENT_SECRET"))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	res, err := client.Do(request)
	if err != nil {
		return tokensRes, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return tokensRes, err
	}

	errDecoder := json.NewDecoder(res.Body).Decode(&tokensRes)
	if errDecoder != nil {
		return tokensRes, errDecoder
	}

	return tokensRes, nil
}
