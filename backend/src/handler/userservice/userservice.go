package userservice_handler

import (
	"backend/src/handler/middleware"
	_ "backend/src/handler/userservice/docs"
	"backend/src/service"
	"net/http"

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
	self.PrivateRoutes(router)
}

func (self *UserServiceHandler) PrivateRoutes(router *gin.Engine) {
	private := router.Group("", middleware.VerifyJWTCookie, middleware.VerifyEmailFromContext, middleware.VerifyConnectionTypeFromContext)
	private.POST("/service-callback/:service/:app-type", self.ServiceCallback)
	private.GET("/Github/user/repositories", self.getGithubUserRepositories)
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
			"error": "Cannot retrieve Github user repositories",
		})
		return
	}

	context.IndentedJSON(http.StatusOK, gin.H{
		"repositories": repositories,
	})
}
