package userservice_handler

import (
	"backend/src/service"
	"net/http"
	"os"

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
	}
	router.GET("/callback/Spotify", self.SpotifyCallBack)
}

func (self *UserServiceHandler) Auth02Spotify(context *gin.Context) {
	spotifyAuthUrl := "https://accounts.spotify.com/authorize" +
		"?response_type=code" +
		"&client_id=" + os.Getenv("SPOTIFY_ID_CLIENT") +
		"&scope=user-modify-playback-state user-read-playback-state" +
		"&redirect_uri=" + os.Getenv("SPOTIFY_CALLBACK")
	context.Redirect(http.StatusTemporaryRedirect, spotifyAuthUrl)
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

	emailStr, err := self.UserServiceService.GetEmailFromContext(context)
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
