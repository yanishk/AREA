package postgres

import (
	"fmt"
	"os"
	"strconv"
	"database/sql"
	"github.com/joho/godotenv"

	"backend/src/storage"
	"backend/src/storage/postgres/user"
	"backend/src/storage/postgres/service"
	"backend/src/storage/postgres/action"
	"backend/src/storage/postgres/trigger"
	"backend/src/storage/postgres/workflow"
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
		UserRepository: user_repository.NewUserRepository(db),
		ServiceRepository: service_repository.NewServiceRepository(db),
		ActionRepository: action_repository.NewActionRepository(db),
		TriggerRepository: trigger_repository.NewTriggerRepository(db),
		WorkflowRepository: workflow_repository.NewWorkflowRepository(db),
	}
}
