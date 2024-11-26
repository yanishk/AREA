package database

import (
	"database/sql"
	"fmt"
)

type Trigger struct {
	Id          string
	Name        string
	Description string
	ServiceId   string
}

func CreateTrigger(db *sql.DB, name string, description string, serviceName string) error {
	sqlStatement := `INSERT INTO triggers (name, description, serviceid) VALUES ($1, $2, $3);`

	_, err := findTriggerByName(db, name)
	if err == nil {
		return fmt.Errorf("Trigger already exist")
	}

	retrievedService, err := findServiceByName(db, serviceName)
	if err != nil {
		return err
	}

	_, err = db.Exec(sqlStatement, name, description, retrievedService.Id)
	if err != nil {
		return err
	}
	return nil
}

func findTriggerByName(db *sql.DB, name string) (Trigger, error) {
	sqlStatement := `SELECT * FROM triggers where name = ($1)`
	var retrievedTrigger Trigger

	row := db.QueryRow(sqlStatement, name)
	err := row.Scan(&retrievedTrigger.Id, &retrievedTrigger.ServiceId, &retrievedTrigger.Name, &retrievedTrigger.Description)
	if err != nil {
		return retrievedTrigger, err
	}
	return retrievedTrigger, nil
}
