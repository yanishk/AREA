package workflow_repository

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"backend/src/entities"
)

type WorkflowRepository struct {
	db *sql.DB
}

func NewWorkflowRepository(db *sql.DB) *WorkflowRepository {
	return &WorkflowRepository{db: db}
}

func (self *WorkflowRepository) CreateWorkflow(name string, ownerId string, actionId string, reactionId string,
	actionParam map[string]interface{}, reactionParam map[string]interface{}) error {
	sqlStatement := `INSERT INTO workflows (name, ownerid, actionid, reactionid, isactivated, actionparam, reactionparam) VALUES ($1, $2, $3, $4, $5, $6, $7);`

	actionParamJson, errMarshallActionParam := json.Marshal(actionParam)
	if errMarshallActionParam != nil {
		return errMarshallActionParam
	}
	reactionParamJson, errMarshallReactionParam := json.Marshal(reactionParam)
	if errMarshallReactionParam != nil {
		return errMarshallReactionParam
	}

	_, err := self.FindWorkflowByNameAndUserId(name, ownerId)
	if err == nil {
		return fmt.Errorf("Workflow already exist")
	}

	_, err = self.db.Exec(sqlStatement, name, ownerId, actionId, reactionId, true, actionParamJson, reactionParamJson)
	if err != nil {
		return err
	}
	return nil
}

func (self *WorkflowRepository) FindWorkflowByNameAndUserId(name string, ownerId string) (entities.Workflow, error) {
	sqlStatement := `SELECT * FROM workflows where name = ($1) AND ownerid = ($2)`
	var retrievedWorkflow entities.Workflow

	row := self.db.QueryRow(sqlStatement, name, ownerId)
	err := row.Scan(&retrievedWorkflow.Id, &retrievedWorkflow.Name, &retrievedWorkflow.OwnerId, &retrievedWorkflow.ReactionId,
		&retrievedWorkflow.ActionId, &retrievedWorkflow.IsActivated, &retrievedWorkflow.CreatedAt, &retrievedWorkflow.ActionParam, &retrievedWorkflow.ReactionParam)
	if err != nil {
		return retrievedWorkflow, err
	}
	return retrievedWorkflow, nil
}

func (self *WorkflowRepository) FindWorkflowByActionId(actionId string) ([]entities.Workflow, error) {
	sqlStatement := `SELECT * FROM workflows where actionid = ($1)`
	var retrievedWorkflows []entities.Workflow

	rows, errQuery := self.db.Query(sqlStatement, actionId)
	if errQuery != nil {
		return nil, errQuery
	}
	defer rows.Close()

	for rows.Next() {
		var workflow entities.Workflow
		var actionParamByte []byte
		var reactionParamByte []byte

		err := rows.Scan(&workflow.Id, &workflow.Name, &workflow.OwnerId, &workflow.ActionId,
			&workflow.ReactionId, &workflow.IsActivated, &workflow.CreatedAt, &actionParamByte, &reactionParamByte)
		if err != nil {
			return nil, err
		}

		if len(actionParamByte) > 0 {
			errActionByte := json.Unmarshal(actionParamByte, &workflow.ActionParam)
			if errActionByte != nil {
				return nil, errActionByte
			}
		}
		if len(reactionParamByte) > 0 {
			errReactionByte := json.Unmarshal(reactionParamByte, &workflow.ReactionParam)
			if errReactionByte != nil {
				return nil, errReactionByte
			}
		}

		retrievedWorkflows = append(retrievedWorkflows, workflow)
	}
	return retrievedWorkflows, nil
}

func (self *WorkflowRepository) FindWorkflowByUserId(userId string) ([]entities.Workflow, error) {
	sqlStatement := `SELECT * FROM workflows where ownerid = ($1)`
	var retrievedWorkflows []entities.Workflow

	rows, errQuery := self.db.Query(sqlStatement, userId)
	if errQuery != nil {
		return nil, errQuery
	}
	defer rows.Close()

	for rows.Next() {
		var workflow entities.Workflow
		var actionParamByte []byte
		var reactionParamByte []byte

		err := rows.Scan(&workflow.Id, &workflow.Name, &workflow.OwnerId, &workflow.ActionId,
			&workflow.ReactionId, &workflow.IsActivated, &workflow.CreatedAt, &actionParamByte, &reactionParamByte)
		if err != nil {
			return nil, err
		}

		if len(actionParamByte) > 0 {
			errActionByte := json.Unmarshal(actionParamByte, &workflow.ActionParam)
			if errActionByte != nil {
				return nil, errActionByte
			}
		}
		if len(reactionParamByte) > 0 {
			errReactionByte := json.Unmarshal(reactionParamByte, &workflow.ReactionParam)
			if errReactionByte != nil {
				return nil, errReactionByte
			}
		}

		retrievedWorkflows = append(retrievedWorkflows, workflow)
	}

	return retrievedWorkflows, nil
}

func (self *WorkflowRepository) ActivateDeactivateWorkflow(workflowId string, isActivated bool) error {
	sqlStatement := `UPDATE workflows SET isactivated = ($1) WHERE id = ($2)`

	_, err := self.db.Exec(sqlStatement, !isActivated, workflowId)
	if err != nil {
		return err
	}
	return nil
}

func (self *WorkflowRepository) FindWorkflowById(workflowId string) (entities.Workflow, error) {
	sqlStatement := `SELECT * FROM workflows where id = ($1)`
	var retrievedWorkflow entities.Workflow
	var actionParamByte []byte
	var reactionParamByte []byte

	row := self.db.QueryRow(sqlStatement, workflowId)
	err := row.Scan(&retrievedWorkflow.Id, &retrievedWorkflow.Name, &retrievedWorkflow.OwnerId, &retrievedWorkflow.ReactionId,
		&retrievedWorkflow.ActionId, &retrievedWorkflow.IsActivated, &retrievedWorkflow.CreatedAt, &actionParamByte, &reactionParamByte)

	if len(actionParamByte) > 0 {
		errActionByte := json.Unmarshal(actionParamByte, &retrievedWorkflow.ActionParam)
		if errActionByte != nil {
			return retrievedWorkflow, errActionByte
		}
	}
	if len(reactionParamByte) > 0 {
		errReactionByte := json.Unmarshal(reactionParamByte, &retrievedWorkflow.ReactionParam)
		if errReactionByte != nil {
			return retrievedWorkflow, errReactionByte
		}
	}

	if err != nil {
		return retrievedWorkflow, err
	}
	return retrievedWorkflow, nil
}
