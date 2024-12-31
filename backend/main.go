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
	_, errCronCreationEveryMinute := cronJob.AddFunc("@every 1m", func() {
		services.WorkflowService.CheckTimeAndDateActions()
	})
	if errCronCreationEveryMinute != nil {
		panic(errCronCreationEveryMinute)
	}
	_, errCronCreation15Minute := cronJob.AddFunc("@every 15m", func() {
		services.WorkflowService.CheckGithubActions()
		services.WorkflowService.CheckRedditActions()
	})
	if errCronCreation15Minute != nil {
		panic(errCronCreation15Minute)
	}

	_, errCronCreation12Hour := cronJob.AddFunc("@every 12h", func() {
		services.WorkflowService.CheckWeatherActions()
	})
	if errCronCreation12Hour != nil {
		panic(errCronCreation12Hour)
	}

	cronJob.Start()

	errHandler := handlers.Run()
	if errHandler != nil {
		panic(errHandler)
	}
}
