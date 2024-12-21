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

const authURL = "auth-url"

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
		authentication.GET("/Spotify", self.Auth02Spotify)
		authentication.GET("/Google", self.Auth02Google)
		authentication.GET("/Discord", self.Auth02Discord)
	}
}

func (self *UserServiceHandler) PrivateRoutes(router *gin.Engine) {
	callback := router.Group("/callback", middleware.VerifyJWTCookie)
	{
		callback.POST("/:service", self.ServiceCallBack)
	}
}

// @Summary		Google OAuth2
// @Description	Get Google OAuth2 URL
// @Tags			Authentication
// @Produce		json
// @Success		200		{object}	docs_userservice.OAuth2Response
// @Router			/authentication/Google [get]
func (self *UserServiceHandler) Auth02Google(context *gin.Context) {
	googleAuthUrl := "https://accounts.google.com/o/oauth2/auth?client_id=" +
		os.Getenv("GOOGLE_CLIENT_ID") + "&response_type=code&access_type=offline&redirect_uri=" +
		os.Getenv("GOOGLE_CALLBACK") + "&scope=email"

	context.IndentedJSON(http.StatusOK, gin.H{
		authURL: googleAuthUrl,
	})
}

// @Summary		Spotify OAuth2
// @Description	Get Spotify OAuth2 URL
// @Tags			Authentication
// @Produce		json
// @Success		200		{object}	docs_userservice.OAuth2Response
// @Router			/authentication/Spotify [get]
func (self *UserServiceHandler) Auth02Spotify(context *gin.Context) {
	spotifyAuthUrl := "https://accounts.spotify.com/authorize?client_id=" +
		os.Getenv("SPOTIFY_CLIENT_ID") + "&response_type=code&redirect_uri=" +
		os.Getenv("SPOTIFY_CALLBACK") + "&scope=user-read-private user-read-email user-modify-playback-state user-read-playback-state user-library-modify"

	context.IndentedJSON(http.StatusOK, gin.H{
		authURL: spotifyAuthUrl,
	})
}

// @Summary		Discord OAuth2
// @Description	Get Discord OAuth2 URL
// @Tags			Authentication
// @Produce		json
// @Success		200		{object}	docs_userservice.OAuth2Response
// @Router			/authentication/Discord [get]
func (self *UserServiceHandler) Auth02Discord(context *gin.Context) {
	discordAuthUrl := "https://discord.com/oauth2/authorize?client_id=" +
		os.Getenv("DISCORD_CLIENT_ID") +
		"&permissions=141376&response_type=code&redirect_uri=" +
		os.Getenv("DISCORD_CALLBACK") +
		"&scope=identify+email+guilds+bot"

	context.IndentedJSON(http.StatusOK, gin.H{
		authURL: discordAuthUrl,
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
// @Router       /callback/{service} [post]
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

	errUpdate := self.UserServiceService.UpdateTokenForService(code, serviceName, emailStr)
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
