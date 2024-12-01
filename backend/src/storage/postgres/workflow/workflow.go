package workflow_repository

import (
	"database/sql"
	"fmt"

	"backend/src/entities"
)

type WorkflowRepository struct {
	db *sql.DB
}

func NewWorkflowRepository(db *sql.DB) *WorkflowRepository {
	return &WorkflowRepository{db: db}
}

func (self *WorkflowRepository) CreateWorkflow(name string, ownerId string, triggerId string, actionId string) error {
	sqlStatement := `INSERT INTO workflows (name, ownerid, triggerid, actionid, isactivated) VALUES ($1, $2, $3, $4, $5);`

	_, err := self.FindWorkflowByName(name)
	if err == nil {
		return fmt.Errorf("Workflow already exist")
	}

	_, err = self.db.Exec(sqlStatement, name, ownerId, triggerId, actionId, true)
	if err != nil {
		return err
	}
	return nil
}

func (self *WorkflowRepository) FindWorkflowByName(name string) (entities.Workflow, error) {
	sqlStatement := `SELECT * FROM workflows where name = ($1)`
	var retrievedWorkflow entities.Workflow

	row := self.db.QueryRow(sqlStatement, name)
	err := row.Scan(&retrievedWorkflow.Id, &retrievedWorkflow.Name, &retrievedWorkflow.OwnerId, &retrievedWorkflow.ActionId,
		&retrievedWorkflow.TriggerId, &retrievedWorkflow.IsActivated, &retrievedWorkflow.CreatedAt)
	if err != nil {
		return retrievedWorkflow, err
	}
	return retrievedWorkflow, nil
}
