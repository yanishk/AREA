package main

import (
	"database/sql"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
	"golang.org/x/crypto/bcrypt"
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
	sqlStatement := `
    INSERT INTO users (email, password)
    VALUES ($1, $2);
    `

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
	sqlStatement := `
    INSERT INTO services (name)
    VALUES ($1);
    `
	_, err := db.Exec(sqlStatement, name)
	if err != nil {
		return err
	}
	return nil
}

func createTrigger(db *sql.DB, name string, description string, serviceName string) error {
	retrievedService, err := findServiceByName(db, serviceName)
	if err != nil {
		return err
	}

	sqlStatement := `INSERT INTO triggers (name, description, serviceid) VALUES ($1, $2, $3);`
	_, err = db.Exec(sqlStatement, name, description, retrievedService.id)
	if err != nil {
		return err
	}
	return nil
}

func createAction(db *sql.DB, name string, description string, serviceName string) error {
	retrievedService, err := findServiceByName(db, serviceName)
	if err != nil {
		return err
	}

	sqlStatement := `INSERT INTO actions (name, description, serviceid) VALUES ($1, $2, $3);`
	_, err = db.Exec(sqlStatement, name, description, retrievedService.id)
	if err != nil {
		return err
	}
	return nil
}

func createWorkflow(db *sql.DB, name string, ownerEmail string, triggerName string, actionName string) error {
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

	sqlStatement := `INSERT INTO workflows (name, ownerid, triggerid, actionid, isactivated) VALUES ($1, $2, $3, $4, $5);`
	_, err = db.Exec(sqlStatement, name, retrievedUser.id, retrievedTrigger.id, retrievedAction.id, true)
	if err != nil {
		return err
	}
	return nil
}

func findUserByEmail(db *sql.DB, email string) (User, error) {
	var retrievedUser User
	sqlStatement := `SELECT * FROM users where email = ($1)`

	row := db.QueryRow(sqlStatement, email)
	err := row.Scan(&retrievedUser.email, &retrievedUser.password, &retrievedUser.id, &retrievedUser.createdAt)
	if err != nil {
		return retrievedUser, err
	}
	return retrievedUser, nil
}

func findServiceByName(db *sql.DB, name string) (Service, error) {
	var retrievedService Service
	sqlStatement := `SELECT * FROM services where name = ($1)`

	row := db.QueryRow(sqlStatement, name)
	err := row.Scan(&retrievedService.id, &retrievedService.name)
	if err != nil {
		return retrievedService, err
	}
	return retrievedService, nil
}

func findTriggerByName(db *sql.DB, name string) (Trigger, error) {
	var retrievedTrigger Trigger
	sqlStatement := `SELECT * FROM triggers where name = ($1)`

	row := db.QueryRow(sqlStatement, name)
	err := row.Scan(&retrievedTrigger.id, &retrievedTrigger.serviceId, &retrievedTrigger.name, &retrievedTrigger.description)
	if err != nil {
		return retrievedTrigger, err
	}
	return retrievedTrigger, nil
}

func findActionByName(db *sql.DB, name string) (Action, error) {
	var retrievedAction Action
	sqlStatement := `SELECT * FROM actions where name = ($1)`

	row := db.QueryRow(sqlStatement, name)
	err := row.Scan(&retrievedAction.id, &retrievedAction.serviceId, &retrievedAction.name, &retrievedAction.description)
	if err != nil {
		return retrievedAction, err
	}
	return retrievedAction, nil
}

func findWorkflowByName(db *sql.DB, name string) (Workflow, error) {
	var retrievedWorkflow Workflow
	sqlStatement := `SELECT * FROM workflows where name = ($1)`

	row := db.QueryRow(sqlStatement, name)
	err := row.Scan(&retrievedWorkflow.id, &retrievedWorkflow.name, &retrievedWorkflow.ownerId, &retrievedWorkflow.actionId,
		&retrievedWorkflow.triggerId, &retrievedWorkflow.isActivated, &retrievedWorkflow.createdAt)
	if err != nil {
		return retrievedWorkflow, err
	}
	return retrievedWorkflow, nil
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

func (databaseContext *Database) loginAuthentification(context *gin.Context) {
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

	foundUser, err := findUserByEmail(databaseContext.db, loginData.Email)

	if err != nil || len(foundUser.email) <= 0 {
		context.IndentedJSON(http.StatusUnauthorized, gin.H{
			"error": "Could not find requested user",
		})
		return
	}

	resBcrypt := bcrypt.CompareHashAndPassword([]byte(foundUser.password), []byte(loginData.Password))

	if resBcrypt == nil {
		context.IndentedJSON(http.StatusOK, gin.H{
			"success": "Connection successfull",
		})
		return
	}
	context.IndentedJSON(http.StatusUnauthorized, gin.H{
		"error": "Wrong password",
	})
}

func (databaseContext *Database) createNewUser(context *gin.Context) {
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

	router.GET("/sum/:numberFrst/:numberSnd", sumNumbers)
	router.GET("/hi", helloServer)
	router.GET("/random/:maxRand", randomNumber)
	router.POST("/login", databaseContext.loginAuthentification)
	router.POST("/register", databaseContext.createNewUser)

	handler := corsAllow.Handler(router)

	err = http.ListenAndServe(":8080", handler)
	if err != nil {
		panic(err)
	}

	defer databaseContext.db.Close()
}
