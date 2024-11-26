package entities

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"

	"poc-goland.com/poc/src/database"
)

type LoginData struct {
	Email string `json:"email"`
	Password string `json:"password"`
}

type FailureInvalidBodyResponse struct {
	Msg string `json:"error"example:"Invalid request body"`
}

type FailureInternalServerErrorResponse struct {
	Msg string `json:"error"example:"Internal server error"`
}

func createToken(email string, context *gin.Context) error {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"email": email,
			"exp":   time.Now().Add(time.Hour * 24).Unix(),
		})

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		return err
	}

	context.SetSameSite(http.SameSiteLaxMode)
	context.SetCookie("JWToken", tokenString, 3600, "", "", false, true)
	return nil
}

type ModifyPasswordData struct {
	Password string `json:"password"`
}

type ModifyPasswordSuccessResponse struct {
	Msg string `json:"success"example:"Password modified"`
}

type ModifyPasswordFailureEmailNotValidResponse struct {
	Msg string `json:"errorr"example:"Email is not a valid string"`
}

type ModifyPasswordFailureInternalServerErrorResponse struct {
	Msg string `json:"error"example:"Failed to hash password-Could not modify the password"`
}

//	@Summary		Modify Password
//	@Description	Modify the user password
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			user	body		ModifyPasswordData	true	"New password of the user"
//	@Success		200		{object}	ModifyPasswordSuccessResponse
//	@Failure		400		{object}	FailureInvalidBodyResponse
//	@Failure		401		{object}	ModifyPasswordFailureEmailNotValidResponse
//	@Failure		500		{object}	ModifyPasswordFailureInternalServerErrorResponse
//	@Router			/modify-password [put]
func (databaseContext *Database) ModifyPassword(context *gin.Context) {
	email, err := context.Get("email")

	if !err {
		context.IndentedJSON(http.StatusUnauthorized, gin.H{
			"error": "Email not found in token",
		})
		return
	}

	emailStr, ok := email.(string)
	if !ok {
		context.IndentedJSON(http.StatusUnauthorized, gin.H{
			"error": "Email is not a valid string",
		})
		return
	}

	var newPassword struct {
		Password string `json:"password"`
	}

	errorBody := context.ShouldBindJSON(&newPassword)
	if errorBody != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	hash, errHash := bcrypt.GenerateFromPassword([]byte(newPassword.Password), bcrypt.DefaultCost)

	if errHash != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to hash password",
		})
		return
	}

	errUpdate := database.UpdateUser(databaseContext.Db, emailStr, string(hash))
	if errUpdate != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": "Could not modify the password",
		})
		return
	}
	context.IndentedJSON(http.StatusOK, gin.H{
		"success": "Password modified",
	})
}

type DeleteAccountSuccessResponse struct {
	Msg string `json:"success"example:"Account deleted"`
}

type DeleteAccountFailureUnauthorizedResponse struct {
	Msg string `json:"error"example:"Email not found in token-Email is not a valid string"`
}

type DeleteAccountFailureInternalServerErrorResponse struct {
	Msg string `json:"error"example:"Could not delete account"`
}

//	@Summary		Delete Account
//	@Description	Delete the user account
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Success		200		{object}	DeleteAccountSuccessResponse
//	@Failure		401		{object}	DeleteAccountFailureUnauthorizedResponse
//	@Failure		500		{object}	DeleteAccountFailureInternalServerErrorResponse
//	@Router			/delete-account [delete]
func (databaseContext *Database) DeleteAccount(context *gin.Context) {
	email, err := context.Get("email")

	if !err {
		context.IndentedJSON(http.StatusUnauthorized, gin.H{
			"error": "Email not found in token",
		})
		return
	}

	emailStr, ok := email.(string)
	if !ok {
		context.IndentedJSON(http.StatusUnauthorized, gin.H{
			"error": "Email is not a valid string",
		})
		return
	}

	errDelete := database.DeleteUser(databaseContext.Db, emailStr)
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

type LoginSuccessResponse struct {
	Msg string `json:"success"example:"Connection successfull"`
}

type LoginFailureUnauthorizedResponse struct {
	Msg string `json:"error"example:"Could not find requested user-Wrong password"`
}
type LoginFailureTokenErrorResponse struct {
	Msg string `json:"error"example:"Error creating token"`
}

//	@Summary		Login
//	@Description	User Authentication
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			user	body		LoginData	true	"User credentials"
//	@Success		200		{object}	LoginSuccessResponse
//	@Failure		400		{object}	FailureInvalidBodyResponse
//	@Failure		401		{object}	LoginFailureUnauthorizedResponse
//	@Failure		500		{object}	LoginFailureTokenErrorResponse
//	@Router			/login [post]
func (databaseContext *Database) LoginAuthentication(context *gin.Context) {
	var user LoginData
	err := context.ShouldBindJSON(&user)
	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	foundUser, err := database.FindUserByEmail(databaseContext.Db, user.Email)

	if err != nil || len(foundUser.Email) <= 0 {
		context.IndentedJSON(http.StatusUnauthorized, gin.H{
			"error": "Could not find requested user",
		})
		return
	}

	resBcrypt := bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(user.Password))

	if resBcrypt == nil {
		errToken := createToken(foundUser.Email, context)
		if errToken != nil {
			context.IndentedJSON(http.StatusInternalServerError, gin.H{
				"error": "Error creating token",
			})
			return
		}
		context.IndentedJSON(http.StatusOK, gin.H{
			"success": "Connection successfull",
		})
		return
	}
	context.IndentedJSON(http.StatusUnauthorized, gin.H{
		"error": "Wrong password",
	})
}

type RegisterSuccessResponse struct {
	Msg string `json:"success"example:"New user created"`
}

type RegisterFailureEmailAlreadyUsedResponse struct {
	Msg string `json:"error"example:"Email address already used"`
}

//	@Summary		Register
//	@Description	Create user account
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			user	body		LoginData	true	"User credentials"
//	@Success		200		{object}	RegisterSuccessResponse
//	@Failure		400		{object}	FailureInvalidBodyResponse
//	@Failure		409		{object}	RegisterFailureEmailAlreadyUsedResponse
//	@Failure		500		{object}	FailureInternalServerErrorResponse
//	@Router			/register [post]
func (databaseContext *Database) CreateNewUser(context *gin.Context) {
	var newUser LoginData

	err := context.ShouldBindJSON(&newUser)
	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
		})
		return
	}

	errorUser := database.CreateUser(databaseContext.Db, newUser.Email, string(hash))
	if errorUser != nil {
		context.IndentedJSON(http.StatusConflict, gin.H{
			"error": "Email address already used",
		})
		return
	}

	context.IndentedJSON(http.StatusOK, gin.H{
		"success": "New user created",
	})
}
