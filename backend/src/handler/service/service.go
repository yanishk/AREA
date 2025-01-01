package service_handler

import (
	"backend/src/entities"
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

	services := private.Group("/services")
	{
		services.GET("/:id", self.RetrieveServiceById)
		services.GET("/actions", self.RetrieveActionsServices)
		services.GET("/reactions", self.RetrieveReactionsServices)
	}

	private.GET("/Discord/user/servers", self.getDiscordUserServers)
	private.GET("/Discord/server/channels", self.getDiscordGuildChannels)
	private.GET("/Asana/user/workspaces", self.getAsanaUserWorkspaces)
	private.GET("/Asana/workspace/assignees", self.getAsanaWorkspaceAssignees)
	private.GET("/Asana/workspace/projects", self.getAsanaWorkspaceProjects)
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
			"error": "Unknown service",
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

// @Summary		Retrieve Reactions From Service
// @Description	Retrieve reactions from requested service
// @Tags			Service
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
// @Param        id     path     string  true  "Service id"
// @Success		200		{object}	docs_service.ServiceRetrieveServiceByIdSuccessResponse
// @Failure		400		{object}	docs_service.ServiceRetrieveServiceByIdBadRequestResponse
// @Failure		500		{object}	docs_service.ServiceRetrieveServiceByIdInternalServerErrorResponse
// @Router			/services/{id} [get]
func (self *ServiceHandler) RetrieveServiceById(context *gin.Context) {
	serviceId := context.Param("id")
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
			"error": InternalServerErrorMessage,
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
			"error": InternalServerErrorMessage,
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

func (self *ServiceHandler) getAsanaUserWorkspaces(context *gin.Context) {
	var workspaces entities.AsanaUserWorkspaces
	url := "https://app.asana.com/api/1.0/workspaces"
	email := context.GetString("email")
	connectionType := context.GetString("connectionType")

	accessToken, errRefreshing := self.UserServiceService.CallApiAndRefresh(email, connectionType, "Asana")
	if errRefreshing != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": "Could not retrieve or refresh the token",
		})
		return
	}

	resp, errRequest := self.ServiceService.ExecuteApiRequest(url, "GET", "Bearer ", accessToken, nil)
	if errRequest != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": "Could not retrieve the user's workspaces",
		})
		return
	}
	defer resp.Body.Close()

	errDecodeWorkspace := json.NewDecoder(resp.Body).Decode(&workspaces)
	if errDecodeWorkspace != nil || len(workspaces.Data) == 0 {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": "Could not decode the user's workspaces",
		})
		return
	}

	context.IndentedJSON(http.StatusOK, gin.H{
		"workspaces": workspaces,
	})
}

func (self *ServiceHandler) getAsanaWorkspaceAssignees(context *gin.Context) {
	var assignees entities.AsanaUserGID
	workspaceId := context.Query("id")
	url := "https://app.asana.com/api/1.0/workspaces/" + workspaceId + "/users"
	email := context.GetString("email")
	connectionType := context.GetString("connectionType")

	accessToken, errRefreshing := self.UserServiceService.CallApiAndRefresh(email, connectionType, "Asana")
	if errRefreshing != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": "Could not retrieve or refresh the token",
		})
		return
	}

	resp, errRequest := self.ServiceService.ExecuteApiRequest(url, "GET", "Bearer ", accessToken, nil)
	if errRequest != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": "Could not retrieve the workspace's assignees",
		})
		return
	}
	defer resp.Body.Close()

	errDecodeAssignee := json.NewDecoder(resp.Body).Decode(&assignees)
	if errDecodeAssignee != nil || len(assignees.Data) == 0 {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": "Could not decode the workspace's assignees",
		})
		return
	}

	context.IndentedJSON(http.StatusOK, gin.H{
		"assignees": assignees,
	})
}

func (self *ServiceHandler) getAsanaWorkspaceProjects(context *gin.Context) {
	var projects entities.AsanaUserGID
	workspaceId := context.Query("id")
	url := "https://app.asana.com/api/1.0/workspaces/" + workspaceId + "/projects"
	email := context.GetString("email")
	connectionType := context.GetString("connectionType")

	accessToken, errRefreshing := self.UserServiceService.CallApiAndRefresh(email, connectionType, "Asana")
	if errRefreshing != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": "Could not retrieve or refresh the token",
		})
		return
	}

	resp, errRequest := self.ServiceService.ExecuteApiRequest(url, "GET", "Bearer ", accessToken, nil)
	if errRequest != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": "Could not retrieve the workspace's assignees",
		})
		return
	}
	defer resp.Body.Close()

	errDecodeAssignee := json.NewDecoder(resp.Body).Decode(&projects)
	if errDecodeAssignee != nil || len(projects.Data) == 0 {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": "Could not decode the workspace's assignees",
		})
		return
	}

	context.IndentedJSON(http.StatusOK, gin.H{
		"projects": projects,
	})
}
