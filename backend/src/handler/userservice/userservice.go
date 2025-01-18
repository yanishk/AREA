package userservice_handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"backend/src/entities"
	"backend/src/handler/middleware"
	_ "backend/src/handler/service/docs"
	_ "backend/src/handler/userservice/docs"
	"backend/src/service"
)

type UserServiceHandler struct {
	UserServiceService service.UserServiceService
}

func NewUserServiceHandler(UserServiceService service.UserServiceService, router *gin.Engine) *UserServiceHandler {
	handler := &UserServiceHandler{UserServiceService: UserServiceService}
	handler.setupRoutes(router)
	return handler
}

func (self *UserServiceHandler) setupRoutes(router *gin.Engine) {
	self.privateRoutes(router)
}

func (self *UserServiceHandler) privateRoutes(router *gin.Engine) {
	private := router.Group("", middleware.VerifyJWTCookie, middleware.VerifyEmailFromContext, middleware.VerifyConnectionTypeFromContext)
	private.POST("/service-callback", self.serviceCallback)
	private.GET("/service-authentication-status", self.getUserServiceAuthenticationStatus)
	private.GET("/github/user/repositories", self.getGithubUserRepositories)
	private.GET("/gitlab/user/projects", self.getGitlabUserProjects)
	private.GET("/discord/user/servers", self.getDiscordUserServers)
	asana := private.Group("/asana")
	{
		asana.GET("/user/workspaces", self.getAsanaUserWorkspaces)
		asana.GET("/workspace/assignees", self.getAsanaWorkspaceAssignees)
		asana.GET("/workspace/projects", self.getAsanaWorkspaceProjects)
		asana.GET("/workspace/tags", self.getAsanaWorkspaceTags)
	}
}

// @Summary      Service Callback
// @Description  Callback for services
// @Tags         Callbacks
// @Produce      json
// @Param        code     query     string  true  "Authorization code given by the service"
// @Param		callback-informations	body		entities.CallbackInformations true	"Callback informations"
// @Failure		200		{object}	docs_userservice.UserServiceServiceCallbackSuccessResponse
// @Failure		400		{object}	docs_userservice.UserServiceServiceCallbackBadRequestResponse
// @Failure		500		{object}	docs_userservice.UserServiceServiceCallbackInternalServerErrorResponse
// @Router       /service-callback [post]
func (self *UserServiceHandler) serviceCallback(context *gin.Context) {
	var callbackInformations entities.CallbackInformations
	code := context.Query("code")
	email := context.GetString("email")
	connectionType := context.GetString("connectionType")

	errorBody := context.ShouldBindJSON(&callbackInformations)
	if errorBody != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	if code == "" {
		context.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid code authorization",
		})
		return
	}

	if callbackInformations.AppType != "web" && callbackInformations.AppType != "mobile" {
		context.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid app type",
		})
		return
	}

	errUpdate := self.UserServiceService.UpdateTokenForService(code, callbackInformations.Service, callbackInformations.AppType, email, connectionType)
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

// @Summary       User Service Authentication Status
// @Description  Checks if connected user is authenticated to the specified service
// @Tags         Authentication
// @Produce      json
// @Param        service     query     string  true  "Service name (Github, Spotify, Discord..)"
// @Failure		200		{object}	docs_userservice.UserServiceGetUserServiceAuthenticationStatusSuccessResponse
// @Failure		400		{object}	docs_userservice.UserServiceGetUserServiceAuthenticationStatusBadRequestResponse
// @Router       /service-authentication-status [get]
func (self *UserServiceHandler) getUserServiceAuthenticationStatus(context *gin.Context) {
	serviceName := context.Query("service")
	email := context.GetString("email")
	connectionType := context.GetString("connectionType")

	isUserAuthenticated, err := self.UserServiceService.RetrieveUserServiceAuthenticationStatus(email, connectionType, serviceName)
	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	context.IndentedJSON(http.StatusOK, gin.H{
		"authenticated": isUserAuthenticated,
	})
}

// @Summary      Retrieve User Repositories
// @Description  Retrieve the repositories of the authenticated user
// @Tags         Github
// @Produce      json
// @Failure		200		{object}	docs_userservice.UserServiceGetGithubUserRepositoriesSuccessResponse
// @Failure		500		{object}	docs_userservice.UserServiceGetGithubUserRepositoriesInternalServerErrorResponse
// @Router       /github/user/repositories [get]
func (self *UserServiceHandler) getGithubUserRepositories(context *gin.Context) {
	email := context.GetString("email")
	connectionType := context.GetString("connectionType")

	repositories, err := self.UserServiceService.RetrieveGithubUserRepositories(email, connectionType)
	if err != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": "Could not retrieve the user's repositories",
		})
		return
	}

	context.IndentedJSON(http.StatusOK, gin.H{
		"repositories": repositories,
	})
}

// @Summary      Retrieve User Projects
// @Description  Retrieve the projects of the authenticated user
// @Tags         Gitlab
// @Produce      json
// @Failure		200		{object}	docs_userservice.UserServiceGetGitlabUserProjectsSuccessResponse
// @Failure		500		{object}	docs_userservice.UserServiceGetGitlabUserProjectsInternalServerErrorResponse
// @Router       /gitlab/user/projects [get]
func (self *UserServiceHandler) getGitlabUserProjects(context *gin.Context) {
	email := context.GetString("email")
	connectionType := context.GetString("connectionType")

	projects, err := self.UserServiceService.RetrieveGitlabUserProjects(email, connectionType)
	if err != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": "Could not retrieve the user's projects",
		})
		return
	}

	context.IndentedJSON(http.StatusOK, gin.H{
		"projects": projects,
	})
}

// @Summary		Retrieve Discord User Servers
// @Description	Retrieve the discord servers from the authenticated user
// @Tags			Discord
// @Produce		json
// @Success		200		{object}	docs_userservice.UserServiceGetDiscordUserServersSuccessResponse
// @Failure		500		{object}	docs_userservice.UserServiceGetDiscordUserServersInternalServerErrorResponse
// @Router			/discord/user/servers [get]
func (self *UserServiceHandler) getDiscordUserServers(context *gin.Context) {
	email := context.GetString("email")
	connectionType := context.GetString("connectionType")

	servers, err := self.UserServiceService.RetrieveDiscordUserServers(email, connectionType)
	if err != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": "Could not retrieve the user's servers",
		})
		return
	}

	context.IndentedJSON(http.StatusOK, gin.H{
		"servers": servers,
	})
}

// @Summary		Retrieve Asana Workspaces
// @Description	Retrieve the workspaces from the authenticated user
// @Tags			Asana
// @Produce		json
// @Success		200		{object}	docs_userservice.UserServiceGetAsanaUserWorkspacesSuccessResponse
// @Failure		500		{object}	docs_userservice.UserServiceGetAsanaUserWorkspacesInternalServerErrorResponse
// @Router /asana/user/workspaces [get]
func (self *UserServiceHandler) getAsanaUserWorkspaces(context *gin.Context) {
	email := context.GetString("email")
	connectionType := context.GetString("connectionType")

	workspaces, err := self.UserServiceService.RetrieveAsanaUserWorkspaces(email, connectionType)
	if err != nil || len(workspaces.Data) == 0 {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": "Could not retrieve the user's workspaces",
		})
		return
	}

	context.IndentedJSON(http.StatusOK, gin.H{
		"workspaces": workspaces,
	})
}

// @Summary		Retrieve Asana Workspace Assignees
// @Description	Retrieve the assignees from the selected workspace
// @Tags			Asana
// @Produce		json
// @Param        id  query      string  true  "workspace id"
// @Success		200		{object}	docs_userservice.UserServiceGetAsanaWorkspaceAssigneesSuccessResponse
// @Failure		500		{object}	docs_userservice.UserServiceGetAsanaWorkspaceAssigneesInternalServerErrorResponse
// @Router /asana/workspace/assignees [get]
func (self *UserServiceHandler) getAsanaWorkspaceAssignees(context *gin.Context) {
	email := context.GetString("email")
	connectionType := context.GetString("connectionType")
	workspaceId := context.Query("id")

	assignees, err := self.UserServiceService.RetrieveAsanaWorkspaceAssignees(email, connectionType, workspaceId)
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
// @Tags			Asana
// @Produce		json
// @Param        id  query      string  true  "workspace id"
// @Success		200		{object}	docs_userservice.UserServiceGetAsanaWorkspaceProjectsSuccessResponse
// @Failure		500		{object}	docs_userservice.UserServiceGetAsanaWorkspaceProjectsInternalServerErrorResponse
// @Router /asana/workspace/projects [get]
func (self *UserServiceHandler) getAsanaWorkspaceProjects(context *gin.Context) {
	email := context.GetString("email")
	connectionType := context.GetString("connectionType")
	workspaceId := context.Query("id")

	projects, err := self.UserServiceService.RetrieveAsanaWorkspaceProjects(email, connectionType, workspaceId)
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
// @Tags			Asana
// @Produce		json
// @Param        id  query      string  true  "workspace id"
// @Success		200		{object}	docs_userservice.UserServiceGetAsanaWorkspaceTagsSuccessResponse
// @Failure		500		{object}	docs_userservice.UserServiceGetAsanaWorkspaceTagsInternalServerErrorResponse
// @Router /asana/workspace/tags [get]
func (self *UserServiceHandler) getAsanaWorkspaceTags(context *gin.Context) {
	email := context.GetString("email")
	connectionType := context.GetString("connectionType")
	workspaceId := context.Query("id")

	tags, err := self.UserServiceService.RetrieveAsanaWorkspaceTags(email, connectionType, workspaceId)
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
