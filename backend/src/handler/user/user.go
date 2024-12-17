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
const EmailNotFoundInTokenMessage = "Email not found in token"
const EmailNotValidMessage = "Email is not a valid string"

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
}

func (self *UserHandler) PrivateRoutes(router *gin.Engine) {
	user := router.Group("/user", middleware.VerifyJWTCookie)
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
// @Success		200		{object}	docs_user.RegisterSuccessResponse
// @Failure		400		{object}	docs_user.FailureInvalidBodyResponse
// @Failure		409		{object}	docs_user.RegisterFailureEmailAlreadyUsedResponse
// @Failure		500		{object}	docs_user.FailureInternalServerErrorResponse
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

	err = self.UserService.CreateUser(newUser.Email, newUser.Password, newUser.ConnectionType)

	if err != nil {
		if err.Error() == "Email address already used" {
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
// @Success		200		{object}	docs_user.LoginSuccessResponse
// @Failure		400		{object}	docs_user.FailureInvalidBodyResponse
// @Failure		401		{object}	docs_user.LoginFailureUnauthorizedResponse
// @Failure		500		{object}	docs_user.LoginFailureTokenErrorResponse
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

	token, err := self.UserService.LoginAuthentication(user.Email, user.Password)
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
	} else {
		context.SetSameSite(http.SameSiteNoneMode)
		context.SetCookie("JWToken", token, 3600, "", "", true, true)
		context.IndentedJSON(http.StatusOK, gin.H{
			"success": "Connection successful",
		})
	}
}

// @Summary		Modify Password
// @Description	Modify the user password
// @Tags			Users
// @Accept			json
// @Produce		json
// @Param			user	body		entities.UserModifyPassword true	"New password of the user"
// @Success		200		{object}	docs_user.ModifyPasswordSuccessResponse
// @Failure		400		{object}	docs_user.ModifyPasswordBadRequestResponse
// @Failure		401		{object}	docs_user.ModifyPasswordFailureEmailNotValidResponse
// @Failure		403		{object}	docs_user.ModifyPasswordFailureOldPasswordNotValidResponse
// @Failure		500		{object}	docs_user.ModifyPasswordFailureInternalServerErrorResponse
// @Router			/user/modify-password [put]
func (self *UserHandler) ModifyPassword(context *gin.Context) {
	email, err := context.Get("email")

	if !err {
		context.IndentedJSON(http.StatusUnauthorized, gin.H{
			"error": EmailNotFoundInTokenMessage,
		})
		return
	}

	emailStr, ok := email.(string)
	if !ok {
		context.IndentedJSON(http.StatusUnauthorized, gin.H{
			"error": EmailNotValidMessage,
		})
		return
	}

	var newPassword entities.UserModifyPassword

	errorBody := context.ShouldBindJSON(&newPassword)
	if errorBody != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": InvalidRequestBodyMessage,
		})
		return
	}

	errorModifyPassword := self.UserService.ModifyPassword(emailStr, newPassword)
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
// @Success		200		{object}	docs_user.DeleteAccountSuccessResponse
// @Failure		401		{object}	docs_user.DeleteAccountFailureUnauthorizedResponse
// @Failure		500		{object}	docs_user.DeleteAccountFailureInternalServerErrorResponse
// @Router			/user/delete-account [delete]
func (self *UserHandler) DeleteAccount(context *gin.Context) {
	email, err := context.Get("email")

	if !err {
		context.IndentedJSON(http.StatusUnauthorized, gin.H{
			"error": EmailNotFoundInTokenMessage,
		})
		return
	}

	emailStr, ok := email.(string)
	if !ok {
		context.IndentedJSON(http.StatusUnauthorized, gin.H{
			"error": EmailNotValidMessage,
		})
		return
	}

	errDelete := self.UserService.DeleteAccount(emailStr)
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
// @Success		200		{object}	docs_user.GetUserEmailSuccessResponse
// @Failure		401		{object}	docs_user.GetUserEmailUnauthorizedResponse
// @Router			/user/email [get]
func (self *UserHandler) GetUserEmail(context *gin.Context) {
	email, err := context.Get("email")

	if !err {
		context.IndentedJSON(http.StatusUnauthorized, gin.H{
			"error": EmailNotFoundInTokenMessage,
		})
		return
	}

	emailStr, ok := email.(string)
	if !ok {
		context.IndentedJSON(http.StatusUnauthorized, gin.H{
			"error": EmailNotValidMessage,
		})
		return
	}
	context.IndentedJSON(http.StatusOK, gin.H{
		"email": emailStr,
	})
}
