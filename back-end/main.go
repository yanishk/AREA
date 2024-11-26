package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/rs/cors"

	"poc-goland.com/poc/back-end/database"
	"poc-goland.com/poc/back-end/entities"
	"poc-goland.com/poc/back-end/routes"
)

func main() {
	var databaseContext entities.Database
	var err error

	databaseContext.Db, err = database.ConnectToDatabase()
	if err != nil {
		panic(err)
	}

	corsAllow := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:8081"},
		AllowedMethods: []string{"GET", "POST"},
		AllowedHeaders: []string{"Content-Type"},
	})

	var router *gin.Engine = gin.Default()

	routes.PublicRoutes(router, databaseContext)
	routes.PrivateRoutes(router, databaseContext)

	handler := corsAllow.Handler(router)

	err = http.ListenAndServe(":8080", handler)
	if err != nil {
		panic(err)
	}

	defer databaseContext.Db.Close()
}
