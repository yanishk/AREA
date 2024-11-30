package entities

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

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

	email, err := context.Get("email")

	if !err {
		context.IndentedJSON(http.StatusUnauthorized, gin.H{
			"error": "Email not found in token",
		})
		return
	}

	emailStr, ok := email.(string)
	if !ok {
		context.IndentedJSON(http.StatusUnauthorized, gin.H{
			"error": "Email is not a valid string",
		})
		return
	}
	foundUser, errSeekMail := database.FindUserByEmail(databaseContext.Db, emailStr)
	if errSeekMail != nil || len(foundUser.Email) <= 0 {
		context.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": "Could not find requested user",
		})
		return
	}

	currentTime := time.Now()
	expiryTime := currentTime.Add(time.Duration(tokenSpotify.ExpiresIn) * time.Second).Format("2006-01-02 15:04:05-07:00")

	_, errFindingService := database.FindUserServiceByName(databaseContext.Db, "spotify", foundUser.Id)
	if errFindingService == nil {
		database.UpdateServiceByUserAndName(databaseContext.Db, "spotify", tokenSpotify.AccessToken, tokenSpotify.RefreshToken, foundUser.Id, expiryTime)
	} else {
		database.CreateService(databaseContext.Db, "spotify", tokenSpotify.AccessToken, tokenSpotify.RefreshToken, foundUser.Id, expiryTime)
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

func (databaseContext *Database) SpotifyCallApi(context *gin.Context) {
	email, err := context.Get("email")

	if !err {
		context.IndentedJSON(http.StatusUnauthorized, gin.H{
			"error": "Email not found in token",
		})
		return
	}

	emailStr, ok := email.(string)
	if !ok {
		context.IndentedJSON(http.StatusUnauthorized, gin.H{
			"error": "Email is not a valid string",
		})
		return
	}
	foundUser, errSeekMail := database.FindUserByEmail(databaseContext.Db, emailStr)
	if errSeekMail != nil || len(foundUser.Email) <= 0 {
		context.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": "Could not find requested user",
		})
		return
	}

	foundService, errServiceFound := database.FindUserServiceByName(databaseContext.Db, "spotify", foundUser.Id)
	if errServiceFound != nil {
		fmt.Println(errServiceFound)
		context.IndentedJSON(http.StatusUnauthorized, gin.H{
			"error": "User is not connected to the service",
		})
		return
	}

	expiryTime, errTimestamp := time.Parse("2006-01-02T15:04:05Z", foundService.ExpiryToken)
	if errTimestamp != nil {
		fmt.Println(errTimestamp)
		context.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": "Could not check timestamp",
		})
		return
	}

	if expiryTime.Before(time.Now()) {
		refreshedToken, errRefresh := databaseContext.refreshToken(foundService.RefreshToken, foundUser.Id)
		if errRefresh != nil {
			context.IndentedJSON(http.StatusInternalServerError, gin.H{
				"error": "Could not refresh token",
			})
			return
		}
		fmt.Println(refreshedToken)
		context.IndentedJSON(http.StatusOK, gin.H{
			"Success": "API successfully called",
		})
	}

	context.IndentedJSON(http.StatusOK, gin.H{
		"Success": "no refresh need and API work",
	})

	//call API instead of this stooooopid print

}

func (databaseContext *Database) refreshToken(refreshToken string, userUuid string) (ResultToken, error) {
	tokenUrl := "https://accounts.spotify.com/api/token"
	jsonBody := fmt.Sprintf(
		"grant_type=refresh_token&refresh_token=%s",
		refreshToken,
	)
	var tokensRefreshRes ResultToken

	request, errRequest := http.NewRequest("POST", tokenUrl, bytes.NewBuffer([]byte(jsonBody)))
	if errRequest != nil {
		return tokensRefreshRes, errRequest
	}

	request.SetBasicAuth(os.Getenv("SPOTIFY_ID_CLIENT"), os.Getenv("SPOTIFY_CLIENT_SECRET"))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	res, err := client.Do(request)
	if err != nil {
		return tokensRefreshRes, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return tokensRefreshRes, err
	}

	errDecoder := json.NewDecoder(res.Body).Decode(&tokensRefreshRes)
	if errDecoder != nil {
		return tokensRefreshRes, errDecoder
	}

	currentTime := time.Now()
	expiryTime := currentTime.Add(time.Duration(tokensRefreshRes.ExpiresIn) * time.Second).Format("2006-01-02 15:04:05-07:00")

	fmt.Println("TOKEN REFRESH -> " + tokensRefreshRes.RefreshToken)
	database.UpdateServiceByUserAndName(databaseContext.Db, "spotify", tokensRefreshRes.AccessToken, tokensRefreshRes.RefreshToken, userUuid, expiryTime)
	return tokensRefreshRes, err
}
