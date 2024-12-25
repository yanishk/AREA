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

const InternalServerErrorMessage = "Internal server error"
const getUserServersFailedMessage = "Could not retrieve the user's servers"
const getServerChannelsFailedMessage = "Could not retrieve the server's channels"

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
		authentication.GET("/:service/:callback-type/:app-type", self.OAuth2Service)
	}
}

func (self *ServiceHandler) PrivateRoutes(router *gin.Engine) {
	private := router.Group("", middleware.VerifyJWTCookie, middleware.VerifyEmailFromContext, middleware.VerifyConnectionTypeFromContext)
	service := private.Group("/:service")
	{
		service.GET("/actions", self.RetrieveActionsFromService)
		service.GET("/reactions", self.RetrieveReactionsFromService)
	}
	private.GET("/service/retrieve", self.RetrieveServiceById)
	private.GET("/actions-services", self.RetrieveActionsServices)
	private.GET("/reactions-services", self.RetrieveReactionsServices)
	private.GET("/Discord/user/servers", self.getDiscordUserServers)
	private.GET("/Discord/server/channels", self.getDiscordGuildChannels)
}

// @Summary		Service OAuth2
// @Description	Get Service OAuth2 URL
// @Tags			Services Authentication
// @Produce		json
// @Param        service  path      string  true  "Service name (Github, Spotify, Discord..)"
// @Param        callback-type  path      string  true  "Callback type (login or service)"
// @Param        app-type  path      string  true  "App type (web or mobile)"
// @Success		200		{object}	docs_service.ServiceOAuth2SuccessResponse
// @Failure		400		{object}	docs_service.ServiceOAuth2BadRequestResponse
// @Router			/authentication/{service}/{callback-type}/{app-type} [get]
func (self *ServiceHandler) OAuth2Service(context *gin.Context) {
	serviceName := context.Param("service")
	callbackType := context.Param("callback-type")
	appType := context.Param("app-type")
	var authUrl string

	if appType != "web" && appType != "mobile" {
		context.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid app type",
		})
	}
	if callbackType != "login" && callbackType != "service" {
		context.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid callback type",
		})
	}

	if (serviceName == "Google") {
		authUrl = self.OAuth2Google(callbackType, appType)
	} else if (serviceName == "Spotify") {
		authUrl = self.OAuth2Spotify(callbackType)
	} else if (serviceName == "Discord") {
		authUrl = self.OAuth2Discord(callbackType)
	} else if (serviceName == "Github") {
		authUrl = self.OAuth2Github(callbackType)
	} else if (serviceName == "Reddit") {
		authUrl = self.OAuth2Reddit(callbackType)
	} else {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": InternalServerErrorMessage,
		})
		return
	}

	context.IndentedJSON(http.StatusOK, gin.H{
		"auth-url": authUrl,
	})
}

func (self *ServiceHandler) getCallbackandClientId(callbackType string, serviceName string, isIdNecessary bool) (string, string) {
	var callbackLink string
	var clientID string

	if callbackType == "login" {
		if isIdNecessary {
			clientID = os.Getenv(serviceName + "_LOGIN_CLIENT_ID")
		}
		callbackLink = os.Getenv(serviceName + "_LOGIN_CALLBACK")
	} else if callbackType == "service" {
		if isIdNecessary {
			clientID = os.Getenv(serviceName + "_SERVICE_CLIENT_ID")
		}
		callbackLink = os.Getenv(serviceName + "_SERVICE_CALLBACK")
	}
	return callbackLink, clientID
}

func (self *ServiceHandler) OAuth2Google(callbackType string, appType string) string {
	var callbackLink string
	var appTypeLink string

	if callbackType == "login" {
		callbackLink = os.Getenv("GOOGLE_LOGIN_CALLBACK")
	} else if callbackType == "service" {
		callbackLink = os.Getenv("GOOGLE_SERVICE_CALLBACK")
	}

	if appType == "web" {
		appTypeLink = os.Getenv("GOOGLE_WEB_CLIENT_ID")
	} else if appType == "mobile" {
		appTypeLink = os.Getenv("GOOGLE_MOBILE_CLIENT_ID")
	}

	return ("https://accounts.google.com/o/oauth2/auth?client_id=" +
		appTypeLink + "&response_type=code&access_type=offline&redirect_uri=" +
		callbackLink + "&scope=email")
}

func (self *ServiceHandler) OAuth2Spotify(callbackType string) string {
	callbackLink, _ := self.getCallbackandClientId(callbackType, "SPOTIFY", false)

	return ("https://accounts.spotify.com/authorize?client_id=" +
		os.Getenv("SPOTIFY_CLIENT_ID") + "&response_type=code&redirect_uri=" +
		callbackLink + "&scope=user-read-private user-read-email user-modify-playback-state user-read-playback-state user-library-modify")
}

func (self *ServiceHandler) OAuth2Discord(callbackType string) string {
	callbackLink, _ := self.getCallbackandClientId(callbackType, "DISCORD", false)

	return ("https://discord.com/oauth2/authorize?client_id=" +
		os.Getenv("DISCORD_CLIENT_ID") +
		"&permissions=141376&response_type=code&redirect_uri=" +
		callbackLink +
		"&scope=identify+email+guilds+bot")
}

func (self *ServiceHandler) OAuth2Github(callbackType string) string {
	callbackLink, clientID := self.getCallbackandClientId(callbackType, "GITHUB", true)

	return ("https://github.com/login/oauth/authorize?" +
		"client_id=" + clientID +
		"&redirect_uri=" + callbackLink +
		"&scope=repo admin:repo_hook write:repo_hook read:repo_hook read:user user:email user:follow read:org" +
		"&state=" + os.Getenv("GITHUB_STATE"))
}

func (self *ServiceHandler) OAuth2Reddit(callbackType string) string {
	callbackLink, clientID := self.getCallbackandClientId(callbackType, "REDDIT", true)

	return ("https://www.reddit.com/api/v1/authorize?" +
		"client_id=" + clientID +
		"&response_type=code" +
		"&state=" + os.Getenv("REDDIT_STATE") +
		"&redirect_uri=" + callbackLink +
		"&duration=permanent" +
		"&scope=identity read submit vote mysubreddits history account edit privatemessages")
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

// @Summary		Retrieve Service
// @Description	Retrieve service by id
// @Tags			Services
// @Produce		json
// @Param        id     query     string  true  "Service id"
// @Success		200		{object}	docs_service.ServiceRetrieveServiceByIdSuccessResponse
// @Failure		400		{object}	docs_service.ServiceRetrieveServiceByIdBadRequestResponse
// @Failure		500		{object}	docs_service.ServiceRetrieveServiceByIdInternalServerErrorResponse
// @Router			/service/retrieve [get]
func (self *ServiceHandler) RetrieveServiceById(context *gin.Context) {
	serviceId := context.Query("id")
	if serviceId == "" {
		context.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid service id",
		})
		return
	}

	service, err := self.ServiceService.FindServiceById(serviceId)
	if err != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": "Could not retrieve requested service",
		})
		return
	}

	context.IndentedJSON(http.StatusOK, gin.H{
		"service": service,
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

	resp, errRequest := self.ServiceService.ExecuteApiRequest("https://discord.com/api/v10/users/@me/guilds", "GET", "Bearer ", accessToken, nil)
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

	resp, errRequest := self.ServiceService.ExecuteApiRequest(getChannelUrl, "GET", "Bot ", tokenBot, nil)
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
