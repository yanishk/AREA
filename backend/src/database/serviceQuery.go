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
	IdUser       string
	ExpiryToken  string
}

func CreateService(db *sql.DB, name string, accessToken string, refreshToken string, uuidUser string, expiry string) error {
	sqlStatement := `INSERT INTO services (name, accessToken, refreshToken, iduser, expirytoken) VALUES ($1, $2, $3, $4, $5);`

	_, err := FindServiceByName(db, name)
	if err == nil {
		return fmt.Errorf("Service already exist")
	}

	_, err = db.Exec(sqlStatement, name, accessToken, refreshToken, uuidUser, expiry)
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

func FindUserServiceByName(db *sql.DB, name string, userId string) (Service, error) {
	sqlStatement := `SELECT * FROM services where name = ($1) AND iduser = ($2)`
	var retrievedService Service

	row := db.QueryRow(sqlStatement, name, userId)
	err := row.Scan(&retrievedService.Id, &retrievedService.Name, &retrievedService.AccessToken, &retrievedService.RefreshToken, &retrievedService.IdUser, &retrievedService.ExpiryToken)
	if err != nil {
		return retrievedService, err
	}
	return retrievedService, nil
}

func UpdateServiceByUserAndName(db *sql.DB, name string, accessToken string, refreshToken string, uuidUser string, expiryToken string) error {
	sqlStatement := `UPDATE services SET accessToken = ($1), refreshToken = ($2), expirytoken = ($3) WHERE name = ($4);`

	_, err := FindUserServiceByName(db, name, uuidUser)
	if err != nil {
		return fmt.Errorf("Service doesn't exist")
	}

	_, err = db.Exec(sqlStatement, accessToken, refreshToken, expiryToken, name)
	if err != nil {
		return err
	}
	return nil
}
