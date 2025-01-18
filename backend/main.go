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
		services.WorkflowService.CheckNewGitlabWorkflows()
		services.WorkflowService.CheckNewGithubWorkflows()
	})
	if errCronCreationEveryMinute != nil {
		panic(errCronCreationEveryMinute)
	}

	_, errCronCreationEvery15Minutes := cronJob.AddFunc("@every 15m", func() {
		services.WorkflowService.CheckRedditActions()
		services.WorkflowService.CheckGithubActions()
	})
	if errCronCreationEvery15Minutes != nil {
		panic(errCronCreationEvery15Minutes)
	}

	_, errCronCreationEvery12Hours := cronJob.AddFunc("@every 12h", func() {
		services.WorkflowService.CheckWeatherActions()
	})
	if errCronCreationEvery12Hours != nil {
		panic(errCronCreationEvery12Hours)
	}

	cronJob.Start()

	errHandler := handlers.Run()
	if errHandler != nil {
		panic(errHandler)
	}
}
