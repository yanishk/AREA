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

const unknownServiceMessage = "Unknown service"
const cannotRetrieveServiceMessage = "Could not retrieve requested service"
const internalServerErrorMessage = "Internal server error"
const getUserServersFailedMessage = "Could not retrieve the user's servers"
const getServerChannelsFailedMessage = "Could not retrieve the server's channels"
const failedRetrievingTokenMessage = "Could not retrieve or refresh the token"

const bearerType = "Bearer "

const asanaWorkspaceRoute = "https://app.asana.com/api/1.0/workspaces/"

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
	discord := private.Group("/Discord")
	{
		discord.GET("/user/servers", self.getDiscordUserServers)
		discord.GET("/server/channels", self.getDiscordGuildChannels)
	}
	asana := private.Group("/Asana")
	{
		asana.GET("/user/workspaces", self.getAsanaUserWorkspaces)
		asana.GET("/workspace/assignees", self.getAsanaWorkspaceAssignees)
		asana.GET("/workspace/projects", self.getAsanaWorkspaceProjects)
		asana.GET("/workspace/tags", self.getAsanaWorkspaceTags)
	}
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
func (self *ServiceHandler) RetrieveReactionsFromService(context *gin.Context) {
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
func (self *ServiceHandler) RetrieveServiceById(context *gin.Context) {
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
func (self *ServiceHandler) RetrieveServiceByActionId(context *gin.Context) {
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
func (self *ServiceHandler) RetrieveServiceByReactionId(context *gin.Context) {
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
			"error": failedRetrievingTokenMessage,
		})
		return
	}

	resp, errRequest := self.ServiceService.ExecuteApiRequest("https://discord.com/api/v10/users/@me/guilds", "GET", bearerType, accessToken, nil)
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

// @Summary		Retrieve Asana Workspaces
// @Description	Retrieve the workspaces from the authenticated user
// @Tags			Service
// @Produce		json
// @Success		200		{object}	docs_service.ServiceGetAsanaUserWorkspacesSuccessResponse
// @Failure		500		{object}	docs_service.ServiceGetAsanaUserWorkspacesInternalServerErrorResponse
// @Router /Asana/user/workspaces [get]
func (self *ServiceHandler) getAsanaUserWorkspaces(context *gin.Context) {
	var workspaces entities.AsanaWorkspacesInfo
	url := "https://app.asana.com/api/1.0/workspaces"
	email := context.GetString("email")
	connectionType := context.GetString("connectionType")

	accessToken, errRefreshing := self.UserServiceService.CallApiAndRefresh(email, connectionType, "Asana")
	if errRefreshing != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": failedRetrievingTokenMessage,
		})
		return
	}

	resp, errRequest := self.ServiceService.ExecuteApiRequest(url, "GET", bearerType, accessToken, nil)
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

func (self *ServiceHandler) decodeRequiredWorkspaceInfo(url string, context *gin.Context) (entities.AsanaWorkspacesInfo, error) {
	var info entities.AsanaWorkspacesInfo
	email := context.GetString("email")
	connectionType := context.GetString("connectionType")

	accessToken, errRefreshing := self.UserServiceService.CallApiAndRefresh(email, connectionType, "Asana")
	if errRefreshing != nil {
		return info, errRefreshing
	}

	resp, errRequest := self.ServiceService.ExecuteApiRequest(url, "GET", bearerType, accessToken, nil)
	if errRequest != nil {
		return info, errRequest
	}
	defer resp.Body.Close()

	errDecode := json.NewDecoder(resp.Body).Decode(&info)
	if errDecode != nil || len(info.Data) == 0 {
		return info, errDecode
	}
	return info, nil
}

// @Summary		Retrieve Asana Workspace Assignees
// @Description	Retrieve the assignees from the selected workspace
// @Tags			Service
// @Produce		json
// @Param        id  query      string  true  "workspace id"
// @Success		200		{object}	docs_service.ServiceGetAsanaWorkspaceAssigneesSuccessResponse
// @Failure		500		{object}	docs_service.ServiceGetAsanaWorkspaceAssigneesInternalServerErrorResponse
// @Router /Asana/workspace/assignees [get]
func (self *ServiceHandler) getAsanaWorkspaceAssignees(context *gin.Context) {
	workspaceId := context.Query("id")
	url := asanaWorkspaceRoute + workspaceId + "/users"

	assignees, err := self.decodeRequiredWorkspaceInfo(url, context)
	if err != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": "Could not retrieve workspace's assignees",
		})
		return
	}

	context.IndentedJSON(http.StatusOK, gin.H{
		"assignees": assignees,
	})
}

// @Summary		Retrieve Asana Workspace Projects
// @Description	Retrieve the projects from the selected workspace
// @Tags			Service
// @Produce		json
// @Param        id  query      string  true  "workspace id"
// @Success		200		{object}	docs_service.ServiceGetAsanaWorkspaceProjectsSuccessResponse
// @Failure		500		{object}	docs_service.ServiceGetAsanaWorkspaceProjectsInternalServerErrorResponse
// @Router /Asana/workspace/projects [get]
func (self *ServiceHandler) getAsanaWorkspaceProjects(context *gin.Context) {
	workspaceId := context.Query("id")
	url := asanaWorkspaceRoute + workspaceId + "/projects"

	projects, err := self.decodeRequiredWorkspaceInfo(url, context)
	if err != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": "Could not retrieve workspace's projects",
		})
		return
	}

	context.IndentedJSON(http.StatusOK, gin.H{
		"projects": projects,
	})
}

// @Summary		Retrieve Asana Workspace Tags
// @Description	Retrieve the tags from the selected workspace
// @Tags			Service
// @Produce		json
// @Param        id  query      string  true  "workspace id"
// @Success		200		{object}	docs_service.ServiceGetAsanaWorkspaceTagsSuccessResponse
// @Failure		500		{object}	docs_service.ServiceGetAsanaWorkspaceTagsInternalServerErrorResponse
// @Router /Asana/workspace/tags [get]
func (self *ServiceHandler) getAsanaWorkspaceTags(context *gin.Context) {
	workspaceId := context.Query("id")
	url := asanaWorkspaceRoute + workspaceId + "/tags"

	tags, err := self.decodeRequiredWorkspaceInfo(url, context)
	if err != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": "Could not retrieve workspace's tags",
		})
		return
	}

	context.IndentedJSON(http.StatusOK, gin.H{
		"tags": tags,
	})
}
