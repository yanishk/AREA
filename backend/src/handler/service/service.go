package service_handler

import (
	"backend/src/service"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

type ServiceHandler struct {
	ServiceService service.ServiceService
}

func NewServiceHandler(ServiceService service.ServiceService, router *gin.Engine) *ServiceHandler {
	handler := &ServiceHandler{ServiceService: ServiceService}
	handler.SetupRoutes(router)
	return handler
}

func (self *ServiceHandler) SetupRoutes(router *gin.Engine) {
	self.PrivateRoutes(router)
}

func (self *ServiceHandler) PrivateRoutes(router *gin.Engine) {
	spotify := router.Group("/spotify")
	{
		spotify.GET("/authentication", self.Auth02Spotify)
		spotify.GET("/call-api", self.SpotifyCallApi)
		spotify.GET("/spotify-callback", self.SpotifyCallBack)
	}
}

func (self *ServiceHandler) Auth02Spotify(context *gin.Context) {
	spotifyAuthUrl := "https://accounts.spotify.com/authorize?client_id=" +
		os.Getenv("SPOTIFY_ID_CLIENT") + "&response_type=code&redirect_uri=" +
		os.Getenv("CALLBACK_SPOTIFY") + "&scope=user-read-private user-read-email"
	context.Redirect(http.StatusTemporaryRedirect, spotifyAuthUrl)
}

func (self *ServiceHandler) getEmailFromContext(context *gin.Context) (string, bool) {
	email, err := context.Get("email")

	if !err {
		context.IndentedJSON(http.StatusUnauthorized, gin.H{
			"error": "Email not found in token",
		})
		return "", err
	}

	emailStr, errEmail := email.(string)
	if !errEmail {
		context.IndentedJSON(http.StatusUnauthorized, gin.H{
			"error": "Email is not a valid string",
		})
		return "", errEmail
	}
	return emailStr, true
}

func (self *ServiceHandler) SpotifyCallBack(context *gin.Context) {
	codeSpotify := context.Query("code")

	if codeSpotify == "" {
		context.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid code authorization",
		})
		return
	}
	tokenSpotify, errToken := self.ServiceService.GetAccesTokenFromCode(codeSpotify)

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

	foundUser, errGetUser := self.ServiceService.GetUser(emailStr)
	if errGetUser != nil || len(foundUser.Email) <= 0 {
		context.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": "Could not find requested user",
		})
		return
	}

	foundService, errFoundService := self.ServiceService.FindServiceByName("Spotify")
	if errFoundService != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": "Could not find requested service",
		})
		return
	}

	self.ServiceService.CreateOrUpdateUserService(foundService.Id, foundUser.Id, tokenSpotify)

	context.IndentedJSON(http.StatusOK, gin.H{
		"access_token": tokenSpotify,
	})
}

func (self *ServiceHandler) SpotifyCallApi(context *gin.Context) {
	emailStr, err := self.getEmailFromContext(context)
	if !err {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get email",
		})
		return
	}

	foundUser, errGetUser := self.ServiceService.GetUser(emailStr)
	if errGetUser != nil || len(foundUser.Email) <= 0 {
		context.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": "Could not find requested user",
		})
		return
	}

	foundService, errServiceFound := self.ServiceService.FindServiceByName("Spotify")
	if errServiceFound != nil {
		context.IndentedJSON(http.StatusUnauthorized, gin.H{
			"error": "Service does not exist",
		})
		return
	}

	foundUserService, errUserService := self.ServiceService.FindUserServiceByIdServiceandIdUser(foundService.Id, foundUser.Id)
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
		_, errRefresh := self.ServiceService.RefreshToken(foundUserService.RefreshToken, foundUser.Id)
		if errRefresh != nil {
			context.IndentedJSON(http.StatusInternalServerError, gin.H{
				"error": "Could not refresh token",
			})
			return
		}
		context.IndentedJSON(http.StatusOK, gin.H{
			"Success": "API successfully called",
		})
	}

	context.IndentedJSON(http.StatusOK, gin.H{
		"Success": "no refresh need and API work",
	})

	// call API instead of this stooooopid print
}
