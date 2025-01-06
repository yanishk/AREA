package service_handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"backend/src/handler/middleware"
	_ "backend/src/handler/service/docs"
	"backend/src/service"
)

type ServiceHandler struct {
	ServiceService service.ServiceService
	UserService    service.UserService
}

const unknownServiceMessage = "Unknown service"
const cannotRetrieveServiceMessage = "Could not retrieve requested service"
const internalServerErrorMessage = "Internal server error"

func NewServiceHandler(ServiceService service.ServiceService, UserService service.UserService, router *gin.Engine) *ServiceHandler {
	handler := &ServiceHandler{
		ServiceService: ServiceService,
		UserService:    UserService,
	}
	handler.setupRoutes(router)
	return handler
}

func (self *ServiceHandler) setupRoutes(router *gin.Engine) {
	self.publicRoutes(router)
	self.privateRoutes(router)
}

func (self *ServiceHandler) publicRoutes(router *gin.Engine) {
	authentication := router.Group("/authentication")
	{
		authentication.GET("/:service/:callback-type/:app-type", self.oauth2Service)
	}
}

func (self *ServiceHandler) privateRoutes(router *gin.Engine) {
	private := router.Group("", middleware.VerifyJWTCookie, middleware.VerifyEmailFromContext, middleware.VerifyConnectionTypeFromContext)
	service := private.Group("/:service")
	{
		service.GET("/actions", self.retrieveActionsFromService)
		service.GET("/reactions", self.retrieveReactionsFromService)
	}
	services := private.Group("/services")
	{
		services.GET("/:id", self.retrieveServiceById)
		services.GET("/action/:actionid", self.retrieveServiceByActionId)
		services.GET("/reaction/:reactionid", self.retrieveServiceByReactionId)
		services.GET("/actions", self.retrieveActionsServices)
		services.GET("/reactions", self.retrieveReactionsServices)
	}
	private.GET("Discord/server/channels", self.getDiscordGuildChannels)
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
func (self *ServiceHandler) oauth2Service(context *gin.Context) {
	serviceName := context.Param("service")
	callbackType := context.Param("callback-type")
	appType := context.Param("app-type")

	if appType != "web" && appType != "mobile" {
		context.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid app type",
		})
		return
	}
	if callbackType != "login" && callbackType != "service" {
		context.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid callback type",
		})
		return
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
func (self *ServiceHandler) retrieveActionsFromService(context *gin.Context) {
	serviceName := context.Param("service")

	actions, errActions := self.ServiceService.RetrieveActionsFromService(serviceName)
	if errActions != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": unknownServiceMessage,
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
func (self *ServiceHandler) retrieveReactionsFromService(context *gin.Context) {
	serviceName := context.Param("service")

	reactions, errReactions := self.ServiceService.RetrieveReactionsFromService(serviceName)
	if errReactions != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": unknownServiceMessage,
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
func (self *ServiceHandler) retrieveServiceById(context *gin.Context) {
	serviceId := context.Param("id")

	service, err := self.ServiceService.FindServiceById(serviceId)
	if err != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": cannotRetrieveServiceMessage,
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
func (self *ServiceHandler) retrieveServiceByActionId(context *gin.Context) {
	id := context.Param("actionid")

	service, err := self.ServiceService.FindServiceByActionId(id)
	if err != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": cannotRetrieveServiceMessage,
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
func (self *ServiceHandler) retrieveServiceByReactionId(context *gin.Context) {
	id := context.Param("reactionid")

	service, err := self.ServiceService.FindServiceByReactionId(id)
	if err != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": cannotRetrieveServiceMessage,
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
func (self *ServiceHandler) retrieveActionsServices(context *gin.Context) {
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
func (self *ServiceHandler) retrieveReactionsServices(context *gin.Context) {
	reactionsServices, err := self.ServiceService.RetrieveReactionsServices()
	if err != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": internalServerErrorMessage,
		})
		return
	}
	context.IndentedJSON(http.StatusOK, reactionsServices)
}

// @Summary		Retrieve Discord Server Channels
// @Description	Retrieve the channels from the selected server
// @Tags			Discord
// @Produce		json
// @Param        id  query      string  true  "server id"
// @Success		200		{object}	docs_service.ServiceGetDiscordGuildChannelsSuccessResponse
// @Failure		500		{object}	docs_service.ServiceGetDiscordGuildChannelsInternalServerErrorResponse
// @Router /Discord/server/channels [get]
func (self *ServiceHandler) getDiscordGuildChannels(context *gin.Context) {
	guildId := context.Query("id")

	channels, err := self.ServiceService.RetrieveDiscordGuildChannels(guildId)
	if err != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": "Could not retrieve the server's channels",
		})
		return
	}

	context.IndentedJSON(http.StatusOK, gin.H{
		"channels": channels,
	})
}
