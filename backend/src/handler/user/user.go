package user_handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"backend/src/entities"
	"backend/src/handler/middleware"
	_ "backend/src/handler/user/docs"
	"backend/src/service"
)

type UserHandler struct {
	UserService service.UserService
}

const InvalidRequestBodyMessage = "Invalid request body"

func NewUserHandler(UserService service.UserService, router *gin.Engine) *UserHandler {
	handler := &UserHandler{UserService: UserService}
	handler.SetupRoutes(router)
	return handler
}

func (self *UserHandler) SetupRoutes(router *gin.Engine) {
	self.PublicRoutes(router)
	self.PrivateRoutes(router)
}

func (self *UserHandler) PublicRoutes(router *gin.Engine) {
	router.POST("/register", self.CreateUser)
	router.POST("/login", self.LoginAuthentication)
	router.POST("/login-callback/:service/:app-type", self.LoginCallback)
}

func (self *UserHandler) PrivateRoutes(router *gin.Engine) {
	private := router.Group("", middleware.VerifyJWTCookie, middleware.VerifyEmailFromContext, middleware.VerifyConnectionTypeFromContext)
	user := private.Group("/user")
	{
		user.PUT("/modify-password", self.ModifyPassword)
		user.DELETE("/delete-account", self.DeleteAccount)
		user.GET("/email", self.GetUserEmail)
	}
}

// @Summary		Register
// @Description	Create user account
// @Tags			Users
// @Accept			json
// @Produce		json
// @Param			user	body		entities.UserCredentials	true	"User credentials"
// @Success		200		{object}	docs_user.UserRegisterSuccessResponse
// @Failure		400		{object}	docs_user.UserRegisterBadRequestResponse
// @Failure		409		{object}	docs_user.UserRegisterConflictResponse
// @Failure		500		{object}	docs_user.UserInternalServerErrorResponse
// @Router			/register [post]
func (self *UserHandler) CreateUser(context *gin.Context) {
	var newUser entities.UserCredentials

	err := context.ShouldBindJSON(&newUser)
	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": InvalidRequestBodyMessage,
		})
		return
	}

	err = self.UserService.CreateUser(newUser.Email, newUser.Password, "basic")

	if err != nil {
		if err.Error() == "Connection type doesn't exist" {
			context.IndentedJSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		} else if err.Error() == "Email address already used" {
			context.IndentedJSON(http.StatusConflict, gin.H{
				"error": err.Error(),
			})
		} else {
			context.IndentedJSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		}
	} else {
		context.IndentedJSON(http.StatusOK, gin.H{
			"success": "New user created",
		})
	}
}

// @Summary		Login
// @Description	User Authentication
// @Tags			Users
// @Accept			json
// @Produce		json
// @Param			user	body		entities.UserCredentials	true	"User credentials"
// @Success		200		{object}	docs_user.UserLoginSuccessResponse
// @Failure		400		{object}	docs_user.UserInvalidBodyResponse
// @Failure		401		{object}	docs_user.UserLoginUnauthorizedResponse
// @Failure		500		{object}	docs_user.UserLoginTokenErrorResponse
// @Router			/login [post]
func (self *UserHandler) LoginAuthentication(context *gin.Context) {
	var user entities.UserCredentials

	err := context.ShouldBindJSON(&user)
	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": InvalidRequestBodyMessage,
		})
		return
	}

	token, err := self.UserService.LoginAuthentication(user.Email, user.Password, "basic")
	if err != nil {
		if err.Error() == "Could not find requested user" || err.Error() == "Wrong password" {
			context.IndentedJSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
		} else {
			context.IndentedJSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		}
		return
	}
	context.SetSameSite(http.SameSiteNoneMode)
	context.SetCookie("JWToken", token, 3600, "", "", true, true)
	context.IndentedJSON(http.StatusOK, gin.H{
		"success": "Connection successful",
	})
}

// @Summary      Login Callback
// @Description  Callback for login
// @Tags         Callbacks
// @Produce      json
// @Param        service  path      string  true  "Service name (Google, Spotify..)"
// @Param        app-type     path     string  true  "App type (web or mobile)"
// @Param        code     query     string  true  "Authorization code given by the service"
// @Failure		200		{object}	docs_user.UserLoginCallbackSuccessResponse
// @Failure		400		{object}	docs_user.UserLoginCallbackBadRequestResponse
// @Failure		500		{object}	docs_user.UserLoginCallbackInternalServerErrorResponse
// @Router       /login-callback/{service}/{app-type} [post]
func (self *UserHandler) LoginCallback(context *gin.Context) {
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

	token, err := self.UserService.LoginWithService(code, serviceName, appType)
	if err != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to connect with requested service",
		})
		return
	}
	context.SetSameSite(http.SameSiteNoneMode)
	context.SetCookie("JWToken", token, 3600, "", "", true, true)
	context.IndentedJSON(http.StatusOK, gin.H{
		"success": "Connection successful",
	})
}

// @Summary		Modify Password
// @Description	Modify the user password
// @Tags			Users
// @Accept			json
// @Produce		json
// @Param			user	body		entities.UserModifyPassword true	"New password of the user"
// @Success		200		{object}	docs_user.UserModifyPasswordSuccessResponse
// @Failure		400		{object}	docs_user.UserModifyPasswordBadRequestResponse
// @Failure		401		{object}	docs_user.UserModifyPasswordUnauthorizedResponse
// @Failure		403		{object}	docs_user.UserModifyPasswordOldPasswordNotValidResponse
// @Failure		500		{object}	docs_user.UserModifyPasswordInternalServerErrorResponse
// @Router			/user/modify-password [put]
func (self *UserHandler) ModifyPassword(context *gin.Context) {
	var newPassword entities.UserModifyPassword
	email := context.GetString("email")
	connectionType := context.GetString("connectionType")

	errorBody := context.ShouldBindJSON(&newPassword)
	if errorBody != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": InvalidRequestBodyMessage,
		})
		return
	}

	errorModifyPassword := self.UserService.ModifyPassword(email, connectionType, newPassword)
	if errorModifyPassword != nil {
		if errorModifyPassword.Error() == "Could not find requested user" {
			context.IndentedJSON(http.StatusBadRequest, gin.H{
				"error": errorModifyPassword.Error(),
			})
		} else if errorModifyPassword.Error() == "Old password is incorrect" {
			context.IndentedJSON(http.StatusForbidden, gin.H{
				"error": errorModifyPassword.Error(),
			})
		} else {
			context.IndentedJSON(http.StatusInternalServerError, gin.H{
				"error": errorModifyPassword.Error(),
			})
		}
	} else {
		context.IndentedJSON(http.StatusOK, gin.H{
			"success": "Password modified",
		})
	}
}

// @Summary		Delete Account
// @Description	Delete the user account
// @Tags			Users
// @Accept			json
// @Produce		json
// @Success		200		{object}	docs_user.UserDeleteAccountSuccessResponse
// @Failure		401		{object}	docs_user.UserDeleteAccountUnauthorizedResponse
// @Failure		500		{object}	docs_user.UserDeleteAccountInternalServerErrorResponse
// @Router			/user/delete-account [delete]
func (self *UserHandler) DeleteAccount(context *gin.Context) {
	email := context.GetString("email")
	connectionType := context.GetString("connectionType")

	errDelete := self.UserService.DeleteAccount(email, connectionType)
	if errDelete != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": "Could not delete account",
		})
		return
	}
	context.IndentedJSON(http.StatusOK, gin.H{
		"success": "Account deleted",
	})
}

// @Summary		Get User Email
// @Description	Get the email of user connected
// @Tags			Users
// @Produce		json
// @Success		200		{object}	docs_user.UserGetUserEmailSuccessResponse
// @Failure		401		{object}	docs_user.UserGetUserEmailUnauthorizedResponse
// @Router			/user/email [get]
func (self *UserHandler) GetUserEmail(context *gin.Context) {
	email := context.GetString("email")

	context.IndentedJSON(http.StatusOK, gin.H{
		"email": email,
	})
}
