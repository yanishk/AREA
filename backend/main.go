package main

import (
	"fmt"

	_ "github.com/lib/pq"
	"github.com/robfig/cron/v3"

	_ "backend/docs"
	"backend/src/handler"
	"backend/src/service/domain"
	"backend/src/storage/postgres"
)

// @title	Documentation for AREA Rest API
// @host	localhost:8080
func main() {
	repositories := postgres.New()
	services := domain.New(repositories)
	handlers := handler.New(services)

	cronJob := cron.New()

	_, errCronCreation := cronJob.AddFunc("@every 1m", func() {
		err := services.WorkflowService.CallTimeApi()
		if err != nil {
			fmt.Println("Call Time API Failed")
		}
		err = services.WorkflowService.CheckGithubActions()
		if err != nil {
			fmt.Println("Check Github Actions Failed")
		}
	})
	if errCronCreation != nil {
		panic(errCronCreation)
	}
	cronJob.Start()

	errHandler := handlers.Run()
	if errHandler != nil {
		panic(errHandler)
	}
}
