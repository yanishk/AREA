package postgres

import (
	"fmt"
	"os"
	"strconv"
	"database/sql"
	"github.com/joho/godotenv"

	"backend/src/storage"
	"backend/src/storage/postgres/user"
	"backend/src/storage/postgres/action"
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
		ActionRepository: action_repository.NewActionRepository(db),
	}
}
