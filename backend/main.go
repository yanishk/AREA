package main

import (
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
	_, errCronCreationTime := cronJob.AddFunc("@every 1m", func() {
		services.WorkflowService.CallTimeApi()
	})
	if errCronCreationTime != nil {
		panic(errCronCreationTime)
	}

	_, errCronCreationReddit := cronJob.AddFunc("@every 15m", func() {
		services.WorkflowService.CallRedditApi()
	})
	if errCronCreationReddit != nil {
		panic(errCronCreationReddit)
	}

	cronJob.Start()

	errHandler := handlers.Run()
	if errHandler != nil {
		panic(errHandler)
	}
}
