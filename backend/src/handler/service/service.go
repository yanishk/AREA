package service_handler

import (
	"backend/src/handler/middleware"
	_ "backend/src/handler/service/docs"
	"backend/src/service"
	"encoding/json"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type ServiceHandler struct {
	ServiceService     service.ServiceService
	UserServiceService service.UserServiceService
	UserService        service.UserService
}

const callbackTypeParamName = "callback-type"
const appTypeParamName = "app-type"

const botBearer = "Bot "
const accessBearer = "Bearer "

const InternalServerErrorMessage = "Internal server error"
const InvalidCallbackTypeMessage = "Invalid callback type"
const InvalidAppTypeMessage = "Invalid app type"
const getUserServersFailedMessage = "Could not retrieve the user's servers"
const getServerChannelsFailedMessage = "Could not retrieve the server's channels"

const authURL = "auth-url"

func NewServiceHandler(ServiceService service.ServiceService, UserServiceService service.UserServiceService, UserService service.UserService, router *gin.Engine) *ServiceHandler {
	handler := &ServiceHandler{
		ServiceService:     ServiceService,
		UserServiceService: UserServiceService,
		UserService:        UserService,
	}
	handler.SetupRoutes(router)
	return handler
}

func (self *ServiceHandler) SetupRoutes(router *gin.Engine) {
	self.PublicRoutes(router)
	self.PrivateRoutes(router)
}

func (self *ServiceHandler) PublicRoutes(router *gin.Engine) {
	authentication := router.Group("/authentication")
	{
		authentication.GET("/Spotify/:callback-type/:app-type", self.OAuth2Spotify)
		authentication.GET("/Google/:callback-type/:app-type", self.OAuth2Google)
		authentication.GET("/Discord/:callback-type/:app-type", self.OAuth2Discord)
	}
}

func (self *ServiceHandler) PrivateRoutes(router *gin.Engine) {
	private := router.Group("", middleware.VerifyJWTCookie, middleware.VerifyEmailFromContext, middleware.VerifyConnectionTypeFromContext)
	service := private.Group("/:service")
	{
		service.GET("/actions", self.RetrieveActionsFromService)
		service.GET("/reactions", self.RetrieveReactionsFromService)
	}
	private.GET("/actions-services", self.RetrieveActionsServices)
	private.GET("/reactions-services", self.RetrieveReactionsServices)
	private.GET("/Discord/user/servers", self.getDiscordUserServers)
	private.GET("/Discord/server/channels", self.getDiscordGuildChannels)
}

// @Summary		Discord User Servers
// @Description	Retrieve the discord servers from the authenticated user
// @Tags			Services
// @Produce		json
// @Success		200		{object}	docs_service.ServiceGetDiscordUserServersSuccessResponse
// @Failure		500		{object}	docs_service.ServiceGetDiscordUserServersInternalServerErrorResponse
// @Router			/Discord/user/servers [get]
func (self *ServiceHandler) getDiscordUserServers(context *gin.Context) {
	email := context.GetString("email")
	connectionType := context.GetString("connectionType")

	accessToken, errRefreshing := self.UserServiceService.CallApiAndRefresh(email, connectionType, "Discord")
	if errRefreshing != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": "Could not retrieve or refresh the token",
		})
		return
	}

	resp, errRequest := self.ServiceService.ExecuteApiRequest("https://discord.com/api/v10/users/@me/guilds", "GET", accessBearer, accessToken, nil)
	if errRequest != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": getUserServersFailedMessage,
		})
		return
	}
	defer resp.Body.Close()

	var guilds []map[string]interface{}
	var ownedServers []map[string]interface{}
	errDecoder := json.NewDecoder(resp.Body).Decode(&guilds)
	if errDecoder != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": getUserServersFailedMessage,
		})
		return
	}

	for _, guild := range guilds {
		if guild["owner"] == true {
			ownedServers = append(ownedServers, guild)
		}
	}

	context.IndentedJSON(http.StatusOK, gin.H{
		"servers": ownedServers,
	})
}

// @Summary		Discord Server Channels
// @Description	Retrieve the channels from the selected server
// @Tags			Services
// @Produce		json
// @Param        id  query      string  true  "server id"
// @Success		200		{object}	docs_service.ServiceGetDiscordGuildChannelsSuccessResponse
// @Failure		500		{object}	docs_service.ServiceGetDiscordGuildChannelsInternalServerErrorResponse
// @Router /Discord/server/channels [get]
func (self *ServiceHandler) getDiscordGuildChannels(context *gin.Context) {
	tokenBot := os.Getenv("DISCORD_BOT_TOKEN")
	guildId := context.Query("id")
	getChannelUrl := "https://discord.com/api/v10/guilds/" + guildId + "/channels"

	resp, errRequest := self.ServiceService.ExecuteApiRequest(getChannelUrl, "GET", botBearer, tokenBot, nil)
	if errRequest != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": getServerChannelsFailedMessage,
		})
		return
	}
	defer resp.Body.Close()

	var channels []map[string]interface{}
	errDecoder := json.NewDecoder(resp.Body).Decode(&channels)
	if errDecoder != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": getServerChannelsFailedMessage,
		})
		return
	}
	context.IndentedJSON(http.StatusOK, gin.H{
		"channels": channels,
	})
}

// @Summary		Google OAuth2
// @Description	Get Google OAuth2 URL
// @Tags			Services Authentication
// @Produce		json
// @Param        callback-type  path      string  true  "Callback type (login or service)"
// @Param        app-type  path      string  true  "App type (web or mobile)"
// @Success		200		{object}	docs_service.ServiceOAuth2SuccessResponse
// @Failure		400		{object}	docs_service.ServiceOAuth2BadRequestResponse
// @Router			/authentication/Google/{callback-type}/{app-type} [get]
func (self *ServiceHandler) OAuth2Google(context *gin.Context) {
	var callbackLink string
	var appTypeLink string
	callbackType := context.Param(callbackTypeParamName)
	appType := context.Param(appTypeParamName)

	if callbackType == "login" {
		callbackLink = os.Getenv("GOOGLE_LOGIN_CALLBACK")
	} else if callbackType == "service" {
		callbackLink = os.Getenv("GOOGLE_SERVICE_CALLBACK")
	} else {
		context.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": InvalidCallbackTypeMessage,
		})
		return
	}

	if appType == "web" {
		appTypeLink = os.Getenv("GOOGLE_WEB_CLIENT_ID")
	} else if appType == "mobile" {
		appTypeLink = os.Getenv("GOOGLE_MOBILE_CLIENT_ID")
	} else {
		context.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": InvalidCallbackTypeMessage,
		})
		return
	}

	googleAuthUrl := "https://accounts.google.com/o/oauth2/auth?client_id=" +
		appTypeLink + "&response_type=code&access_type=offline&redirect_uri=" +
		callbackLink + "&scope=email"

	context.IndentedJSON(http.StatusOK, gin.H{
		authURL: googleAuthUrl,
	})
}

// @Summary		Spotify OAuth2
// @Description	Get Spotify OAuth2 URL
// @Tags			Services Authentication
// @Produce		json
// @Param        callback-type  path      string  true  "Callback type (login or service)"
// @Param        app-type  path      string  true  "App type (web or mobile)"
// @Success		200		{object}	docs_service.ServiceOAuth2SuccessResponse
// @Failure		400		{object}	docs_service.ServiceOAuth2BadRequestResponse
// @Router			/authentication/Spotify/{callback-type}/{app-type} [get]
func (self *ServiceHandler) OAuth2Spotify(context *gin.Context) {
	var callbackLink string
	callbackType := context.Param(callbackTypeParamName)
	appType := context.Param(appTypeParamName)

	if callbackType == "login" {
		callbackLink = os.Getenv("SPOTIFY_LOGIN_CALLBACK")
	} else if callbackType == "service" {
		callbackLink = os.Getenv("SPOTIFY_SERVICE_CALLBACK")
	} else {
		context.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": InvalidCallbackTypeMessage,
		})
		return
	}

	if appType != "web" && appType != "mobile" {
		context.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": InvalidCallbackTypeMessage,
		})
		return
	}

	spotifyAuthUrl := "https://accounts.spotify.com/authorize?client_id=" +
		os.Getenv("SPOTIFY_CLIENT_ID") + "&response_type=code&redirect_uri=" +
		callbackLink + "&scope=user-read-private user-read-email user-modify-playback-state user-read-playback-state user-library-modify"

	context.IndentedJSON(http.StatusOK, gin.H{
		authURL: spotifyAuthUrl,
	})
}

// @Summary		Discord OAuth2
// @Description	Get Discord OAuth2 URL
// @Tags			Services Authentication
// @Produce		json
// @Param        callback-type  path      string  true  "Callback type (login or service)"
// @Param        app-type  path      string  true  "App type (web or mobile)"
// @Success		200		{object}	docs_service.ServiceOAuth2SuccessResponse
// @Failure		400		{object}	docs_service.ServiceOAuth2BadRequestResponse
// @Router			/authentication/Discord/{callback-type}/{app-type} [get]
func (self *ServiceHandler) OAuth2Discord(context *gin.Context) {
	var callbackLink string
	callbackType := context.Param(callbackTypeParamName)
	appType := context.Param(appTypeParamName)

	if callbackType == "login" {
		callbackLink = os.Getenv("DISCORD_LOGIN_CALLBACK")
	} else if callbackType == "service" {
		callbackLink = os.Getenv("DISCORD_SERVICE_CALLBACK")
	} else {
		context.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": InvalidCallbackTypeMessage,
		})
		return
	}

	if appType != "web" && appType != "mobile" {
		context.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": InvalidCallbackTypeMessage,
		})
		return
	}
	discordAuthUrl := "https://discord.com/oauth2/authorize?client_id=" +
		os.Getenv("DISCORD_CLIENT_ID") +
		"&permissions=141376&response_type=code&redirect_uri=" +
		callbackLink +
		"&scope=identify+email+guilds+bot"

	context.IndentedJSON(http.StatusOK, gin.H{
		authURL: discordAuthUrl,
	})
}

// @Summary		Actions Services
// @Description	Retrieve services with actions available
// @Tags			Services
// @Produce		json
// @Success		200		{object}	docs_service.ServiceRetrieveActionsReactionsServicesSuccessResponse
// @Failure		500		{object}	docs_service.ServiceInternalServerErrorResponse
// @Router			/actions-services [get]
func (self *ServiceHandler) RetrieveActionsServices(context *gin.Context) {
	actionsServices, err := self.ServiceService.RetrieveActionsServices()
	if err != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": InternalServerErrorMessage,
		})
		return
	}
	context.IndentedJSON(http.StatusOK, actionsServices)
}

// @Summary		Reactions Services
// @Description	Retrieve services with reactions available
// @Tags			Services
// @Produce		json
// @Success		200		{object}	docs_service.ServiceRetrieveActionsReactionsServicesSuccessResponse
// @Failure		500		{object}	docs_service.ServiceInternalServerErrorResponse
// @Router			/reactions-services [get]
func (self *ServiceHandler) RetrieveReactionsServices(context *gin.Context) {
	reactionsServices, err := self.ServiceService.RetrieveReactionsServices()
	if err != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": InternalServerErrorMessage,
		})
		return
	}
	context.IndentedJSON(http.StatusOK, reactionsServices)
}

// @Summary		Actions From Service
// @Description	Retrieve actions from requested service
// @Tags			Services
// @Produce		json
// @Param        service  path      string  true  "Service name (Google, Spotify, Time & Date..)"
// @Success		200		{object}	docs_service.ServiceRetrieveActionsFromServiceSuccessResponse
// @Failure		500		{object}	docs_service.ServiceInternalServerErrorResponse
// @Router			/{service}/actions [get]
func (self *ServiceHandler) RetrieveActionsFromService(context *gin.Context) {
	serviceName := context.Param("service")

	actions, errActions := self.ServiceService.RetrieveActionsFromService(serviceName)
	if errActions != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": InternalServerErrorMessage,
		})
		return
	}

	context.IndentedJSON(http.StatusOK, gin.H{
		"actions": actions,
	})
}

// @Summary		Reactions From Service
// @Description	Retrieve reactions from requested service
// @Tags			Services
// @Produce		json
// @Param        service  path      string  true  "Service name (Google, Spotify, Time & Date..)"
// @Success		200		{object}	docs_service.ServiceRetrieveReactionsFromServiceSuccessResponse
// @Failure		500		{object}	docs_service.ServiceInternalServerErrorResponse
// @Router			/{service}/reactions [get]
func (self *ServiceHandler) RetrieveReactionsFromService(context *gin.Context) {
	serviceName := context.Param("service")

	reactions, errReactions := self.ServiceService.RetrieveReactionsFromService(serviceName)
	if errReactions != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": InternalServerErrorMessage,
		})
		return
	}

	context.IndentedJSON(http.StatusOK, gin.H{
		"reactions": reactions,
	})
}
