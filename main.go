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
)

type User struct {
	id        string
	email     string
	password  string
	createdAt string
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

func createUser(db *sql.DB, email string, password string) {
	sqlStatement := `
    INSERT INTO users (email, password)
    VALUES ($1, $2, $3);
    `
	_, err := db.Exec(sqlStatement, email, password)
	if err != nil {
		panic(err)
	}
}

func findUserByEmail(db *sql.DB, email string) (User, error) {
	var retrievedUser User
	sqlStatement := `SELECT * FROM users where email = ($1)`
	rows, err := db.Query(sqlStatement, email)

	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&retrievedUser.id, &retrievedUser.email,
			&retrievedUser.password, &retrievedUser.createdAt)
		if err != nil {
			panic(err)
		}
	}
	err = rows.Err()
	if err != nil {
		panic(err)
	}
	return retrievedUser, err
}

func createService(db *sql.DB, name string) {
	sqlStatement := `
    INSERT INTO services (name)
    VALUES ($1);
    `
	_, err := db.Exec(sqlStatement, name)
	if err != nil {
		panic(err)
	}
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

func main() {
	db, err := connectToDatabase()
	if err != nil {
		panic(err)
	}
	var router *gin.Engine = gin.Default()
	router.GET("/sum/:numberFrst/:numberSnd", sumNumbers)
	router.GET("/hi", helloServer)
	router.GET("/random/:maxRand", randomNumber)
	router.Run(":8080")
	defer db.Close()
}
