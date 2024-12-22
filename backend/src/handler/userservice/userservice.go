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
	router.POST("/service-callback/:service/:app-type", middleware.VerifyJWTCookie, self.ServiceCallback)
}

// @Summary      Service Callback
// @Description  Callback for services
// @Tags         Callbacks
// @Produce      json
// @Param        service  path      string  true  "Service name (Google, Spotify..)"
// @Param        app-type     path     string  true  "App type (web or mobile)"
// @Param        code     query     string  true  "Authorization code given by the service"
// @Failure		200		{object}	docs_userservice.ServiceCallbackSuccessResponse
// @Failure		400		{object}	docs_userservice.ServiceCallbackBadRequestResponse
// @Failure		500		{object}	docs_userservice.ServiceCallbackInternalServerErrorResponse
// @Router       /service-callback/{service}/{app-type} [post]
func (self *UserServiceHandler) ServiceCallback(context *gin.Context) {
	appType := context.Param("app-type")
	serviceName := context.Param("service")
	code := context.Query("code")

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

	emailStr, err := self.UserServiceService.GetEmailFromContext(context)
	if !err {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get email",
		})
		return
	}

	connectionTypeStr, err := self.UserServiceService.GetConnectionTypeFromContext(context)
	if !err {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get connection type",
		})
		return
	}

	errUpdate := self.UserServiceService.UpdateTokenForService(code, serviceName, appType, emailStr, connectionTypeStr)
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
