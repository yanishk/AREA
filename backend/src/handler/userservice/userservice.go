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
	authentication := router.Group("/authentication")
	{
		authentication.GET("/Spotify", self.Auth02Spotify)
		authentication.GET("/Google", self.Auth02Google)
	}

	callback := router.Group("/callback")
	{
		callback.GET("/:service", self.ServiceCallBack)
	}

	callApi := router.Group("/call-api")
	{
		callApi.GET("/Spotify", self.SpotifyCallApi)
	}
}

func (self *UserServiceHandler) Auth02Google(context *gin.Context) {
	googleAuthUrl := "https://accounts.google.com/o/oauth2/auth?client_id=" +
		os.Getenv("GOOGLE_CLIENT_ID") + "&response_type=code&access_type=offline&redirect_uri=" +
		os.Getenv("GOOGLE_CALLBACK") + "&scope=email"
	context.Redirect(http.StatusTemporaryRedirect, googleAuthUrl)
}

func (self *UserServiceHandler) Auth02Spotify(context *gin.Context) {
	spotifyAuthUrl := "https://accounts.spotify.com/authorize?client_id=" +
		os.Getenv("SPOTIFY_CLIENT_ID") + "&response_type=code&redirect_uri=" +
		os.Getenv("SPOTIFY_CALLBACK") + "&scope=user-read-private user-read-email"
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

func (self *UserServiceHandler) ServiceCallBack(context *gin.Context) {
	code := context.Query("code")
	serviceName := context.Param("service")

	if code == "" {
		context.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid code authorization",
		})
		return
	}

	token, errToken := self.UserServiceService.GetAccessTokenFromCode(code, serviceName)
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

	foundService, errFoundService := self.UserServiceService.FindServiceByName(serviceName)
	if errFoundService != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": "Could not find requested service",
		})
		return
	}

	self.UserServiceService.CreateOrUpdateUserService(foundUser.Id, foundService.Id, token)

	context.IndentedJSON(http.StatusOK, gin.H{
		"access_token": token,
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

	foundUserService, errUserService := self.UserServiceService.FindUserServiceByIdServiceandIdUser(foundUser.Id, foundService.Id)
	if errUserService != nil {
		context.IndentedJSON(http.StatusUnauthorized, gin.H{
			"error": "User is not registered to the service",
		})
		return
	}

	expiryTime, errTimestamp := time.Parse("2006-01-02 15:04:05-07:00", foundUserService.ExpiryDate)
	if errTimestamp != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": "Could not check timestamp",
		})
		return
	}

	if expiryTime.Before(time.Now()) {
		_, errRefresh := self.UserServiceService.RefreshToken(foundUserService.RefreshToken, foundUser.Id, "Spotify")
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
