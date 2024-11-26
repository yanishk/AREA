package database

import (
	"database/sql"
	"fmt"
)

type Service struct {
	Id   string
	Name string
}

func createService(db *sql.DB, name string) error {
	sqlStatement := `INSERT INTO services (name) VALUES ($1);`

	_, err := findServiceByName(db, name)
	if err == nil {
		return fmt.Errorf("Service already exist")
	}

	_, err = db.Exec(sqlStatement, name)
	if err != nil {
		return err
	}
	return nil
}

func findServiceByName(db *sql.DB, name string) (Service, error) {
	sqlStatement := `SELECT * FROM services where name = ($1)`
	var retrievedService Service

	row := db.QueryRow(sqlStatement, name)
	err := row.Scan(&retrievedService.Id, &retrievedService.Name)
	if err != nil {
		return retrievedService, err
	}
	return retrievedService, nil
}
