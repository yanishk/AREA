package database

import (
	"database/sql"
	"fmt"
)

type Action struct {
	Id          string
	Name        string
	Description string
	ServiceId   string
}

func createAction(db *sql.DB, name string, description string, serviceName string) error {
	sqlStatement := `INSERT INTO actions (name, description, serviceid) VALUES ($1, $2, $3);`

	_, err := findActionByName(db, name)
	if err == nil {
		return fmt.Errorf("Action already exist")
	}

	retrievedService, err := FindServiceByName(db, serviceName)
	if err != nil {
		return err
	}

	_, err = db.Exec(sqlStatement, name, description, retrievedService.Id)
	if err != nil {
		return err
	}
	return nil
}

func findActionByName(db *sql.DB, name string) (Action, error) {
	sqlStatement := `SELECT * FROM actions where name = ($1)`
	var retrievedAction Action

	row := db.QueryRow(sqlStatement, name)
	err := row.Scan(&retrievedAction.Id, &retrievedAction.ServiceId, &retrievedAction.Name, &retrievedAction.Description)
	if err != nil {
		return retrievedAction, err
	}
	return retrievedAction, nil
}
