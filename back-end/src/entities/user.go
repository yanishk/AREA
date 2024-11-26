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

func (databaseContext *Database) LoginAuthentification(context *gin.Context) {
	var loginData struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := context.ShouldBindJSON(&loginData)
	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	foundUser, err := database.FindUserByEmail(databaseContext.Db, loginData.Email)

	if err != nil || len(foundUser.Email) <= 0 {
		context.IndentedJSON(http.StatusUnauthorized, gin.H{
			"error": "Could not find requested user",
		})
		return
	}

	resBcrypt := bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(loginData.Password))

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

func (databaseContext *Database) CreateNewUser(context *gin.Context) {
	var newUser struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

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
