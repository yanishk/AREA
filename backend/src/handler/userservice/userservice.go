package userservice_handler

import (
	"backend/src/service"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

type UserServiceHandler struct {
	UserServiceService service.UserServiceService
}

func NewUserServiceHandler(UserServiceService service.UserServiceService, router *gin.Engine) *UserServiceHandler {
	handler := &UserServiceHandler{UserServiceService: UserServiceService}
	handler.SetupRoutes(router)
	return handler
}

func (self *UserServiceHandler) SetupRoutes(router *gin.Engine) {
	self.PrivateRoutes(router)
}

func (self *UserServiceHandler) PrivateRoutes(router *gin.Engine) {
	spotify := router.Group("/spotify")
	{
		spotify.GET("/authentication", self.Auth02Spotify)
		spotify.GET("/call-api", self.SpotifyCallApi)
		spotify.GET("/spotify-callback", self.SpotifyCallBack)
	}
}

func (self *UserServiceHandler) Auth02Spotify(context *gin.Context) {
	spotifyAuthUrl := "https://accounts.spotify.com/authorize?client_id=" +
		os.Getenv("SPOTIFY_ID_CLIENT") + "&response_type=code&redirect_uri=" +
		os.Getenv("CALLBACK_SPOTIFY") + "&scope=user-read-private user-read-email"
	context.Redirect(http.StatusTemporaryRedirect, spotifyAuthUrl)
}

func (self *UserServiceHandler) getEmailFromContext(context *gin.Context) (string, bool) {
	email, err := context.Get("email")

	if !err {
		return "", err
	}

	emailStr, errEmail := email.(string)
	if !errEmail {
		return "", errEmail
	}
	return emailStr, true
}

func (self *UserServiceHandler) SpotifyCallBack(context *gin.Context) {
	codeSpotify := context.Query("code")

	if codeSpotify == "" {
		context.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid code authorization",
		})
		return
	}

	tokenSpotify, errToken := self.UserServiceService.GetAccesTokenFromCode(codeSpotify)
	if errToken != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to exchange token",
		})
		return
	}

	emailStr, err := self.getEmailFromContext(context)
	if !err {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get email",
		})
		return
	}

	foundUser, errGetUser := self.UserServiceService.GetUser(emailStr)
	if errGetUser != nil || len(foundUser.Email) <= 0 {
		context.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": "Could not find requested user",
		})
		return
	}

	foundService, errFoundService := self.UserServiceService.FindServiceByName("Spotify")
	if errFoundService != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": "Could not find requested service",
		})
		return
	}

	self.UserServiceService.CreateOrUpdateUserService(foundService.Id, foundUser.Id, tokenSpotify)

	context.IndentedJSON(http.StatusOK, gin.H{
		"access_token": tokenSpotify,
	})
}

func (self *UserServiceHandler) SpotifyCallApi(context *gin.Context) {
	emailStr, err := self.getEmailFromContext(context)
	if !err {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get email",
		})
		return
	}

	foundUser, errGetUser := self.UserServiceService.GetUser(emailStr)
	if errGetUser != nil || len(foundUser.Email) <= 0 {
		context.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": "Could not find requested user",
		})
		return
	}

	foundService, errServiceFound := self.UserServiceService.FindServiceByName("Spotify")
	if errServiceFound != nil {
		context.IndentedJSON(http.StatusUnauthorized, gin.H{
			"error": "Service does not exist",
		})
		return
	}

	foundUserService, errUserService := self.UserServiceService.FindUserServiceByIdServiceandIdUser(foundService.Id, foundUser.Id)
	if errUserService != nil {
		context.IndentedJSON(http.StatusUnauthorized, gin.H{
			"error": "User is not registered to the service",
		})
		return
	}

	expiryTime, errTimestamp := time.Parse("2006-01-02 15:04:05-07:00", foundUserService.ExpiryToken)
	if errTimestamp != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": "Could not check timestamp",
		})
		return
	}

	if expiryTime.Before(time.Now()) {
		_, errRefresh := self.UserServiceService.RefreshToken(foundUserService.RefreshToken, foundUser.Id)
		if errRefresh != nil {
			context.IndentedJSON(http.StatusInternalServerError, gin.H{
				"error": "Could not refresh token",
			})
			return
		}
		context.IndentedJSON(http.StatusOK, gin.H{
			"Success": "API successfully called",
		})
		return
	}

	context.IndentedJSON(http.StatusOK, gin.H{
		"Success": "no refresh need and API work",
	})

	// call API instead of this stooooopid print
}
