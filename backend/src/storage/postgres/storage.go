package postgres

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"

	"backend/src/storage"
	action_repository "backend/src/storage/postgres/action"
	service_repository "backend/src/storage/postgres/service"
	trigger_repository "backend/src/storage/postgres/trigger"
	user_repository "backend/src/storage/postgres/user"
	user_service_repository "backend/src/storage/postgres/userservice"
	workflow_repository "backend/src/storage/postgres/workflow"
)

func retrieveDatabaseInfos() (string, int, string, string, string, error) {
	port, err := strconv.Atoi(os.Getenv("PGPORT"))
	if err != nil {
		panic(err)
	}
	return os.Getenv("PGHOST"), port, os.Getenv("PGUSER"), os.Getenv("PGPASSWORD"), os.Getenv("PGDATABASE"), err
}

func New() *storage.Repository {
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

	return &storage.Repository{
		UserRepository:        user_repository.NewUserRepository(db),
		ServiceRepository:     service_repository.NewServiceRepository(db),
		UserServiceRepository: user_service_repository.NewUserServiceRepository(db),
		ReactionRepository:    action_repository.NewReactionRepository(db),
		TriggerRepository:     trigger_repository.NewTriggerRepository(db),
		WorkflowRepository:    workflow_repository.NewWorkflowRepository(db),
	}
}
