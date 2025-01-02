package userservice_handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"backend/src/handler/middleware"
	_ "backend/src/handler/userservice/docs"
	"backend/src/service"
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
	private := router.Group("", middleware.VerifyJWTCookie, middleware.VerifyEmailFromContext, middleware.VerifyConnectionTypeFromContext)
	private.POST("/service-callback/:service/:app-type", self.ServiceCallback)
	private.GET("/Github/user/repositories", self.getGithubUserRepositories)
	private.GET("/Discord/user/servers", self.getDiscordUserServers)
	asana := private.Group("/Asana")
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
// @Param        service  path      string  true  "Service name (Google, Spotify..)"
// @Param        app-type     path     string  true  "App type (web or mobile)"
// @Param        code     query     string  true  "Authorization code given by the service"
// @Failure		200		{object}	docs_userservice.UserServiceServiceCallbackSuccessResponse
// @Failure		400		{object}	docs_userservice.UserServiceServiceCallbackBadRequestResponse
// @Failure		500		{object}	docs_userservice.UserServiceServiceCallbackInternalServerErrorResponse
// @Router       /service-callback/{service}/{app-type} [post]
func (self *UserServiceHandler) ServiceCallback(context *gin.Context) {
	appType := context.Param("app-type")
	serviceName := context.Param("service")
	code := context.Query("code")
	email := context.GetString("email")
	connectionType := context.GetString("connectionType")

	if code == "" {
		context.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid code authorization",
		})
		return
	}

	if appType != "web" && appType != "mobile" {
		context.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid app type",
		})
		return
	}

	errUpdate := self.UserServiceService.UpdateTokenForService(code, serviceName, appType, email, connectionType)
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

// @Summary      Retrieve User Repositories
// @Description  Retrieve the repositories of the authenticated user
// @Tags         Github
// @Produce      json
// @Failure		200		{object}	docs_userservice.UserServiceGetGithubUserRepositoriesSuccessResponse
// @Failure		500		{object}	docs_userservice.UserServiceGetGithubUserRepositoriesInternalServerErrorResponse
// @Router       /Github/user/repositories [get]
func (self *UserServiceHandler) getGithubUserRepositories(context *gin.Context) {
	email := context.GetString("email")
	connectionType := context.GetString("connectionType")

	repositories, err := self.UserServiceService.RetrieveGithubUserRepositories(email, connectionType)
	if err != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": "Could not retrieve the user's servers",
		})
		return
	}

	context.IndentedJSON(http.StatusOK, gin.H{
		"repositories": repositories,
	})
}

// @Summary		Retrieve Discord User Servers
// @Description	Retrieve the discord servers from the authenticated user
// @Tags			Discord
// @Produce		json
// @Success		200		{object}	docs_userservice.UserServiceGetDiscordUserServersSuccessResponse
// @Failure		500		{object}	docs_userservice.UserServiceGetDiscordUserServersInternalServerErrorResponse
// @Router			/Discord/user/servers [get]
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
// @Router /Asana/user/workspaces [get]
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
// @Router /Asana/workspace/assignees [get]
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
// @Router /Asana/workspace/projects [get]
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
// @Router /Asana/workspace/tags [get]
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
