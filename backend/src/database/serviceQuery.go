package database

import (
	"database/sql"
	"fmt"
)

type Service struct {
	Id           string
	Name         string
	AccessToken  string
	RefreshToken string
}

func CreateService(db *sql.DB, name string, accessToken string, refreshToken string) error {
	sqlStatement := `INSERT INTO services (name, accessToken, refreshToken) VALUES ($1, $2, $3);`

	_, err := FindServiceByName(db, name)
	if err == nil {
		return fmt.Errorf("Service already exist")
	}

	_, err = db.Exec(sqlStatement, name, accessToken, refreshToken)
	if err != nil {
		return err
	}
	return nil
}

func FindServiceByName(db *sql.DB, name string) (Service, error) {
	sqlStatement := `SELECT * FROM services where name = ($1)`
	var retrievedService Service

	row := db.QueryRow(sqlStatement, name)
	err := row.Scan(&retrievedService.Id, &retrievedService.Name, &retrievedService.AccessToken, &retrievedService.RefreshToken)
	if err != nil {
		return retrievedService, err
	}
	return retrievedService, nil
}
