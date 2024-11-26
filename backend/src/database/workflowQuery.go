package database

import (
	"database/sql"
	"fmt"
)

type Workflow struct {
	Id          string
	Name        string
	OwnerId     string
	ActionId    string
	TriggerId   string
	IsActivated bool
	CreatedAt   string
}

func createWorkflow(db *sql.DB, name string, ownerEmail string, triggerName string, actionName string) error {
	sqlStatement := `INSERT INTO workflows (name, ownerid, triggerid, actionid, isactivated) VALUES ($1, $2, $3, $4, $5);`

	_, err := findWorkflowByName(db, name)
	if err == nil {
		return fmt.Errorf("Workflow already exist")
	}

	retrievedUser, err := FindUserByEmail(db, ownerEmail)
	if err != nil {
		return err
	}

	retrievedTrigger, err := findTriggerByName(db, triggerName)
	if err != nil {
		return err
	}

	retrievedAction, err := findActionByName(db, actionName)
	if err != nil {
		return err
	}

	_, err = db.Exec(sqlStatement, name, retrievedUser.Id, retrievedTrigger.Id, retrievedAction.Id, true)
	if err != nil {
		return err
	}
	return nil
}

func findWorkflowByName(db *sql.DB, name string) (Workflow, error) {
	sqlStatement := `SELECT * FROM workflows where name = ($1)`
	var retrievedWorkflow Workflow

	row := db.QueryRow(sqlStatement, name)
	err := row.Scan(&retrievedWorkflow.Id, &retrievedWorkflow.Name, &retrievedWorkflow.OwnerId, &retrievedWorkflow.ActionId,
		&retrievedWorkflow.TriggerId, &retrievedWorkflow.IsActivated, &retrievedWorkflow.CreatedAt)
	if err != nil {
		return retrievedWorkflow, err
	}
	return retrievedWorkflow, nil
}
