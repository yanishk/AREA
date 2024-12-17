package userservice_handler

import (
	"backend/src/handler/middleware"
	_ "backend/src/handler/userservice/docs"
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
	self.PublicRoutes(router)
	self.PrivateRoutes(router)
}

func (self *UserServiceHandler) PublicRoutes(router *gin.Engine) {
	authentication := router.Group("/authentication")
	{
		authentication.GET("/Spotify/:callback-type", self.Auth02Spotify) // faire pareil ici
		authentication.GET("/Google/:callback-type", self.Auth02Google)
	}
}

func (self *UserServiceHandler) PrivateRoutes(router *gin.Engine) {
	callback := router.Group("/callback", middleware.VerifyJWTCookie)
	{
		callback.GET("/:service", self.ServiceCallBack)
	}
}

// @Summary		Google OAuth2
// @Description	Get Google OAuth2 URL
// @Tags			Authentication
// @Produce		json
// @Param        callback-type  path      string  true  "Callback type (login or classic)"
// @Success		200		{object}	docs_userservice.OAuth2SuccessResponse
// @Failure		400		{object}	docs_userservice.OAuth2BadRequestResponse
// @Router			/authentication/Google/{callback-type} [get]
func (self *UserServiceHandler) Auth02Google(context *gin.Context) {
	var callbackLink string
	callbackType := context.Param("callback-type")

	if (callbackType == "login") {
		callbackLink = os.Getenv("GOOGLE_LOGIN_CALLBACK")
	} else if (callbackType == "classic") {
		callbackLink = os.Getenv("GOOGLE_CALLBACK")
	} else {
		context.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid callback type",
		})
		return
	}

	googleAuthUrl := "https://accounts.google.com/o/oauth2/auth?client_id=" +
		os.Getenv("GOOGLE_CLIENT_ID") + "&response_type=code&access_type=offline&redirect_uri=" +
		callbackLink + "&scope=email"

	context.IndentedJSON(http.StatusOK, gin.H{
		"auth-url": googleAuthUrl,
	})
}

// @Summary		Spotify OAuth2
// @Description	Get Spotify OAuth2 URL
// @Tags			Authentication
// @Produce		json
// @Param        callback-type  path      string  true  "Callback type (login or classic)"
// @Success		200		{object}	docs_userservice.OAuth2SuccessResponse
// @Failure		400		{object}	docs_userservice.OAuth2BadRequestResponse
// @Router			/authentication/Spotify/callback-type} [get]
func (self *UserServiceHandler) Auth02Spotify(context *gin.Context) {
	var callbackLink string
	callbackType := context.Param("callback-type")

	if (callbackType == "login") {
		callbackLink = os.Getenv("SPOTIFY_LOGIN_CALLBACK")
	} else if (callbackType == "classic") {
		callbackLink = os.Getenv("SPOTIFY_CALLBACK")
	} else {
		context.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid callback type",
		})
		return
	}

	spotifyAuthUrl := "https://accounts.spotify.com/authorize?client_id=" +
		os.Getenv("SPOTIFY_CLIENT_ID") + "&response_type=code&redirect_uri=" +
		callbackLink + "&scope=user-read-private user-read-email user-modify-playback-state user-read-playback-state"

	context.IndentedJSON(http.StatusOK, gin.H{
		"auth-url": spotifyAuthUrl,
	})
}

// @Summary      Service Callback
// @Description  Callback for services
// @Tags         Authentication
// @Produce      json
// @Param        service  path      string  true  "Service name (Google, Spotify..)"
// @Param        code     query     string  true  "Authorization code given by the service"
// @Failure		200		{object}	docs_userservice.ServiceCallBackSuccessResponse
// @Failure		400		{object}	docs_userservice.ServiceCallBackBadRequestResponse
// @Failure		500		{object}	docs_userservice.ServiceCallBackInternalServerErrorResponse
// @Router       /callback/{service} [get]
func (self *UserServiceHandler) ServiceCallBack(context *gin.Context) {
	code := context.Query("code")
	serviceName := context.Param("service")

	if code == "" {
		context.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid code authorization",
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

	connectionTypeStr, err := self.UserServiceService.GetConnectionTypeFromContext(context)
	if !err {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get connection type",
		})
		return
	}

	errUpdate := self.UserServiceService.UpdateTokenForService(code, serviceName, emailStr, connectionTypeStr)
	if errUpdate != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update token",
		})
		return
	}

	context.IndentedJSON(http.StatusOK, gin.H{
		"success": "Token generated",
	})
}
