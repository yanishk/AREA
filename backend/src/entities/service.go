package entities

import (
	"bytes"
	"encoding/json"
	"fmt"
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

func getAccesTokenFromCode(codeSpotify string) (string, error) {
	tokenUrl := "https://accounts.spotify.com/api/token"
	var result struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		ExpiresIn   int    `json:"expires_in"`
	}
	jsonBody := fmt.Sprintf(
		"grant_type=authorization_code&code=%s&redirect_uri=%s",
		codeSpotify,
		os.Getenv("CALLBACK_SPOTIFY"),
	)

	request, errRequest := http.NewRequest("POST", tokenUrl, bytes.NewBuffer([]byte(jsonBody)))
	if errRequest != nil {
		return "", errRequest
	}

	request.SetBasicAuth(os.Getenv("SPOTIFY_ID_CLIENT"), os.Getenv("SPOTIFY_CLIENT_SECRET"))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	res, err := client.Do(request)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return "", err
	}

	errDecoder := json.NewDecoder(res.Body).Decode(&result)
	if errDecoder != nil {
		return "", errDecoder
	}

	return result.AccessToken, nil
}
