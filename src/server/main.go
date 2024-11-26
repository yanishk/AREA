package main

import (
	"database/sql"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
	"golang.org/x/crypto/bcrypt"
	"poc-goland.com/poc/src/server/middleware"
	ginSwagger "github.com/swaggo/gin-swagger"
    swaggerFiles "github.com/swaggo/files"
	_ "poc-goland.com/poc/src/server/docs"
)

type User struct {
	id        string
	email     string
	password  string
	createdAt string
}

type Service struct {
	id   string
	name string
}

type Trigger struct {
	id          string
	name        string
	description string
	serviceId   string
}

type Action struct {
	id          string
	name        string
	description string
	serviceId   string
}

type Workflow struct {
	id          string
	name        string
	ownerId     string
	actionId    string
	triggerId   string
	isActivated bool
	createdAt   string
}

type Database struct {
	db *sql.DB
}

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


func retrieveDatabaseInfos() (string, int, string, string, string, error) {
	port, err := strconv.Atoi(os.Getenv("PGPORT"))
	if err != nil {
		panic(err)
	}
	return os.Getenv("PGHOST"), port, os.Getenv("PGUSER"), os.Getenv("PGPASSWORD"), os.Getenv("PGDATABASE"), err
}

func connectToDatabase() (*sql.DB, error) {
	godotenv.Load()

	host, port, user, password, dbname, err := retrieveDatabaseInfos()
	if err != nil {
		panic(err)
	}

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=require",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")
	return db, err
}

func createUser(db *sql.DB, email string, password string) error {
	sqlStatement := `INSERT INTO users (email, password) VALUES ($1, $2);`

	_, err := findUserByEmail(db, email)
	if err == nil {
		return fmt.Errorf("User already exist")
	}

	_, err = db.Exec(sqlStatement, email, password)
	if err != nil {
		return err
	}
	return nil
}

func createService(db *sql.DB, name string) error {
	sqlStatement := `INSERT INTO services (name) VALUES ($1);`

	_, err := findServiceByName(db, name)
	if err == nil {
		return fmt.Errorf("Service already exist")
	}

	_, err = db.Exec(sqlStatement, name)
	if err != nil {
		return err
	}
	return nil
}

func createTrigger(db *sql.DB, name string, description string, serviceName string) error {
	sqlStatement := `INSERT INTO triggers (name, description, serviceid) VALUES ($1, $2, $3);`

	_, err := findTriggerByName(db, name)
	if err == nil {
		return fmt.Errorf("Trigger already exist")
	}

	retrievedService, err := findServiceByName(db, serviceName)
	if err != nil {
		return err
	}

	_, err = db.Exec(sqlStatement, name, description, retrievedService.id)
	if err != nil {
		return err
	}
	return nil
}

func createAction(db *sql.DB, name string, description string, serviceName string) error {
	sqlStatement := `INSERT INTO actions (name, description, serviceid) VALUES ($1, $2, $3);`

	_, err := findActionByName(db, name)
	if err == nil {
		return fmt.Errorf("Action already exist")
	}

	retrievedService, err := findServiceByName(db, serviceName)
	if err != nil {
		return err
	}

	_, err = db.Exec(sqlStatement, name, description, retrievedService.id)
	if err != nil {
		return err
	}
	return nil
}

func createWorkflow(db *sql.DB, name string, ownerEmail string, triggerName string, actionName string) error {
	sqlStatement := `INSERT INTO workflows (name, ownerid, triggerid, actionid, isactivated) VALUES ($1, $2, $3, $4, $5);`

	_, err := findWorkflowByName(db, name)
	if err == nil {
		return fmt.Errorf("Workflow already exist")
	}

	retrievedUser, err := findUserByEmail(db, ownerEmail)
	if err != nil {
		return err
	}

	retrievedTrigger, err := findTriggerByName(db, triggerName)
	if err != nil {
		return err
	}

	retrievedAction, err := findActionByName(db, actionName)
	if err != nil {
		return err
	}

	_, err = db.Exec(sqlStatement, name, retrievedUser.id, retrievedTrigger.id, retrievedAction.id, true)
	if err != nil {
		return err
	}
	return nil
}

func findUserByEmail(db *sql.DB, email string) (User, error) {
	sqlStatement := `SELECT * FROM users where email = ($1)`
	var retrievedUser User

	row := db.QueryRow(sqlStatement, email)
	err := row.Scan(&retrievedUser.email, &retrievedUser.password, &retrievedUser.id, &retrievedUser.createdAt)
	if err != nil {
		return retrievedUser, err
	}
	return retrievedUser, nil
}

func findServiceByName(db *sql.DB, name string) (Service, error) {
	sqlStatement := `SELECT * FROM services where name = ($1)`
	var retrievedService Service

	row := db.QueryRow(sqlStatement, name)
	err := row.Scan(&retrievedService.id, &retrievedService.name)
	if err != nil {
		return retrievedService, err
	}
	return retrievedService, nil
}

func findTriggerByName(db *sql.DB, name string) (Trigger, error) {
	sqlStatement := `SELECT * FROM triggers where name = ($1)`
	var retrievedTrigger Trigger

	row := db.QueryRow(sqlStatement, name)
	err := row.Scan(&retrievedTrigger.id, &retrievedTrigger.serviceId, &retrievedTrigger.name, &retrievedTrigger.description)
	if err != nil {
		return retrievedTrigger, err
	}
	return retrievedTrigger, nil
}

func findActionByName(db *sql.DB, name string) (Action, error) {
	sqlStatement := `SELECT * FROM actions where name = ($1)`
	var retrievedAction Action

	row := db.QueryRow(sqlStatement, name)
	err := row.Scan(&retrievedAction.id, &retrievedAction.serviceId, &retrievedAction.name, &retrievedAction.description)
	if err != nil {
		return retrievedAction, err
	}
	return retrievedAction, nil
}

func findWorkflowByName(db *sql.DB, name string) (Workflow, error) {
	sqlStatement := `SELECT * FROM workflows where name = ($1)`
	var retrievedWorkflow Workflow

	row := db.QueryRow(sqlStatement, name)
	err := row.Scan(&retrievedWorkflow.id, &retrievedWorkflow.name, &retrievedWorkflow.ownerId, &retrievedWorkflow.actionId,
		&retrievedWorkflow.triggerId, &retrievedWorkflow.isActivated, &retrievedWorkflow.createdAt)
	if err != nil {
		return retrievedWorkflow, err
	}
	return retrievedWorkflow, nil
}

func updateUser(db *sql.DB, email string, password string) error {
	sqlStatement := `UPDATE users SET password = ($1) WHERE email = ($2);`

	_, err := findUserByEmail(db, email)
	if err != nil {
		return fmt.Errorf("User doesn't exist")
	}

	_, err = db.Exec(sqlStatement, password, email)
	if err != nil {
		return err
	}
	return nil
}

func deleteUser(db *sql.DB, email string) error {
	sqlStatement := `DELETE FROM users WHERE email = ($1)`

	_, err := findUserByEmail(db, email)
	if err != nil {
		return fmt.Errorf("User doesn't exist")
	}

	_, err = db.Exec(sqlStatement, email)
	if err != nil {
		return err
	}
	return nil
}

func sumNumbers(context *gin.Context) {
	firstNbr, err1 := strconv.Atoi(context.Param("numberFrst"))
	sndNbr, err2 := strconv.Atoi(context.Param("numberSnd"))

	if err1 != nil || err2 != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid numbers",
		})
		return
	}

	sum := firstNbr + sndNbr

	context.JSON(http.StatusOK, gin.H{
		"Result of sum": sum,
	})
}

func helloServer(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{
		"Message": "HI YOU!",
	})
}

func randomNumber(context *gin.Context) {
	maximunNbr, err := strconv.Atoi(context.Param("maxRand"))

	if err != nil || maximunNbr <= 0 {
		context.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid number",
		})
		return
	}
	randNbr := rand.Intn(maximunNbr) + 1

	context.JSON(http.StatusOK, gin.H{
		"Random Number": randNbr,
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
func (databaseContext *Database) loginAuthentication(context *gin.Context) {
	var user LoginData
	err := context.ShouldBindJSON(&user)
	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	foundUser, err := findUserByEmail(databaseContext.db, user.Email)

	if err != nil || len(foundUser.email) <= 0 {
		context.IndentedJSON(http.StatusUnauthorized, gin.H{
			"error": "Could not find requested user",
		})
		return
	}

	resBcrypt := bcrypt.CompareHashAndPassword([]byte(foundUser.password), []byte(user.Password))

	if resBcrypt == nil {
		errToken := createToken(foundUser.email, context)
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
func (databaseContext *Database) createNewUser(context *gin.Context) {
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

	errorUser := createUser(databaseContext.db, newUser.Email, string(hash))
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
func (databaseContext *Database) modifyPassword(context *gin.Context) {
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

	errUpdate := updateUser(databaseContext.db, emailStr, string(hash))
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
func (databaseContext *Database) deleteAccount(context *gin.Context) {
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

	errDelete := deleteUser(databaseContext.db, emailStr)
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

//	@title	Documentation for AREA Rest API
//	@host	localhost:8080
func main() {
	var databaseContext Database
	var err error

	databaseContext.db, err = connectToDatabase()
	if err != nil {
		panic(err)
	}

	corsAllow := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:8081"},
		AllowedMethods: []string{"GET", "POST"},
		AllowedHeaders: []string{"Content-Type"},
	})

	var router *gin.Engine = gin.Default()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.GET("/hi", middleware.VerifyJWTCookie, helloServer)
	router.GET("/sum/:numberFrst/:numberSnd", sumNumbers)
	router.GET("/random/:maxRand", randomNumber)

	router.POST("/login", databaseContext.loginAuthentication)
	router.POST("/register", databaseContext.createNewUser)

	router.PUT("/modify-password", middleware.VerifyJWTCookie, databaseContext.modifyPassword)

	router.DELETE("/delete-account", middleware.VerifyJWTCookie, databaseContext.deleteAccount)

	handler := corsAllow.Handler(router)

	err = http.ListenAndServe(":8080", handler)
	if err != nil {
		panic(err)
	}

	defer databaseContext.db.Close()
}
