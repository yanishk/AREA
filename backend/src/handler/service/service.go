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

const unknownServiceMessage = "Unknown service"
const cantRetrieveServiceMessage = "Could not retrieve requested service"
const internalServerErrorMessage = "Internal server error"
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

	services := private.Group("/services")
	{
		services.GET("/:id", self.RetrieveServiceById)
		services.GET("/action/:actionid", self.RetrieveServiceByActionId)
		services.GET("/reaction/:reactionid", self.RetrieveServiceByReactionId)
		services.GET("/actions", self.RetrieveActionsServices)
		services.GET("/reactions", self.RetrieveReactionsServices)
	}

	private.GET("/Discord/user/servers", self.getDiscordUserServers)
	private.GET("/Discord/server/channels", self.getDiscordGuildChannels)
}

// @Summary		Service OAuth2
// @Description	Get Service OAuth2 URL
// @Tags			Authentication
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

	authUrl, err := self.ServiceService.OAuth2Service(serviceName, callbackType, appType)
	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": unknownServiceMessage,
		})
		return
	}

	context.IndentedJSON(http.StatusOK, gin.H{
		"auth-url": authUrl,
	})
}

// @Summary		Retrieve Actions From Service
// @Description	Retrieve actions from requested service
// @Tags			Service
// @Produce		json
// @Param        service  path      string  true  "Service name (Google, Spotify, Time & Date..)"
// @Success		200		{object}	docs_service.ServiceRetrieveActionsFromServiceSuccessResponse
// @Failure		400		{object}	docs_service.ServiceRetrieveActionsFromServiceBadRequestResponse
// @Router			/{service}/actions [get]
func (self *ServiceHandler) RetrieveActionsFromService(context *gin.Context) {
	serviceName := context.Param("service")

	actions, errActions := self.ServiceService.RetrieveActionsFromService(serviceName)
	if errActions != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": UnknownServiceMessage,
		})
		return
	}

	context.IndentedJSON(http.StatusOK, gin.H{
		"actions": actions,
	})
}

// @Summary		Retrieve Reactions From Service
// @Description	Retrieve reactions from requested service
// @Tags			Service
// @Produce		json
// @Param        service  path      string  true  "Service name (Google, Spotify, Time & Date..)"
// @Success		200		{object}	docs_service.ServiceRetrieveReactionsFromServiceSuccessResponse
// @Failure		400		{object}	docs_service.ServiceRetrieveReactionsFromServiceBadRequestResponse
// @Router			/{service}/reactions [get]
func (self *ServiceHandler) RetrieveReactionsFromService(context *gin.Context) {
	serviceName := context.Param("service")

	reactions, errReactions := self.ServiceService.RetrieveReactionsFromService(serviceName)
	if errReactions != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": UnknownServiceMessage,
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
// @Param        id     path     string  true  "Service id"
// @Success		200		{object}	docs_service.ServiceRetrieveServiceByIdSuccessResponse
// @Failure		500		{object}	docs_service.ServiceRetrieveServiceByIdInternalServerErrorResponse
// @Router			/services/{id} [get]
func (self *ServiceHandler) RetrieveServiceById(context *gin.Context) {
	serviceId := context.Param("id")

	service, err := self.ServiceService.FindServiceById(serviceId)
	if err != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": cantRetrieveServiceMessage,
		})
		return
	}

	context.IndentedJSON(http.StatusOK, gin.H{
		"service": service,
	})
}

// @Summary		Retrieve Service
// @Description	Retrieve service by action id
// @Tags			Services
// @Produce		json
// @Param        actionid     path     string  true  "Action id"
// @Success		200		{object}	docs_service.ServiceRetrieveServiceByIdSuccessResponse
// @Failure		500		{object}	docs_service.ServiceRetrieveServiceByIdInternalServerErrorResponse
// @Router			/services/action/{actionid} [get]
func (self *ServiceHandler) RetrieveServiceByActionId(context *gin.Context) {
	id := context.Param("actionid")

	service, err := self.ServiceService.FindServiceByActionId(id)
	if err != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": cantRetrieveServiceMessage,
		})
		return
	}

	context.IndentedJSON(http.StatusOK, gin.H{
		"service": service,
	})
}

// @Summary		Retrieve Service
// @Description	Retrieve service by reaction id
// @Tags			Services
// @Produce		json
// @Param        reactionid     path     string  true  "Reaction id"
// @Success		200		{object}	docs_service.ServiceRetrieveServiceByIdSuccessResponse
// @Failure		500		{object}	docs_service.ServiceRetrieveServiceByIdInternalServerErrorResponse
// @Router			/services/reaction/{reactionid} [get]
func (self *ServiceHandler) RetrieveServiceByReactionId(context *gin.Context) {
	id := context.Param("reactionid")

	service, err := self.ServiceService.FindServiceByReactionId(id)
	if err != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": cantRetrieveServiceMessage,
		})
		return
	}

	context.IndentedJSON(http.StatusOK, gin.H{
		"service": service,
	})
}

// @Summary		Retrieve Services With Actions
// @Description	Retrieve services with actions available
// @Tags			Services
// @Produce		json
// @Success		200		{object}	docs_service.ServiceRetrieveActionsReactionsServicesSuccessResponse
// @Failure		500		{object}	docs_service.ServiceInternalServerErrorResponse
// @Router			/services/actions [get]
func (self *ServiceHandler) RetrieveActionsServices(context *gin.Context) {
	actionsServices, err := self.ServiceService.RetrieveActionsServices()
	if err != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": internalServerErrorMessage,
		})
		return
	}
	context.IndentedJSON(http.StatusOK, actionsServices)
}

// @Summary		Retrieve Services With Reactions
// @Description	Retrieve services with reactions available
// @Tags			Services
// @Produce		json
// @Success		200		{object}	docs_service.ServiceRetrieveActionsReactionsServicesSuccessResponse
// @Failure		500		{object}	docs_service.ServiceInternalServerErrorResponse
// @Router			/services/reactions [get]
func (self *ServiceHandler) RetrieveReactionsServices(context *gin.Context) {
	reactionsServices, err := self.ServiceService.RetrieveReactionsServices()
	if err != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": internalServerErrorMessage,
		})
		return
	}
	context.IndentedJSON(http.StatusOK, reactionsServices)
}

// @Summary		Retrieve Discord User Servers
// @Description	Retrieve the discord servers from the authenticated user
// @Tags			Service
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

// @Summary		Retrieve Discord Server Channels
// @Description	Retrieve the channels from the selected server
// @Tags			Service
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
