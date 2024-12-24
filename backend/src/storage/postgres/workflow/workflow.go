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
	sqlStatement := `INSERT INTO workflows (name, ownerid, actionid, reactionid, isactivated, actionparam, reactionparam) VALUES ($1, $2, $3, $4, $5, $6, $7)`

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
	sqlStatement := `SELECT * FROM workflows WHERE name = ($1) AND ownerid = ($2)`
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
	sqlStatement := `SELECT * FROM workflows WHERE actionid = ($1)`

	rows, errQuery := self.db.Query(sqlStatement, actionId)
	if errQuery != nil {
		return nil, errQuery
	}
	defer rows.Close()

	retrievedWorkflows, err := self.appendWorkflowsSlices(rows)
	if err != nil {
		return nil, err
	}
	return retrievedWorkflows, nil
}

func (self *WorkflowRepository) FindWorkflowByUserId(userId string) ([]entities.Workflow, error) {
	sqlStatement := `SELECT * FROM workflows WHERE ownerid = ($1)`

	rows, errQuery := self.db.Query(sqlStatement, userId)
	if errQuery != nil {
		return nil, errQuery
	}
	defer rows.Close()

	retrievedWorkflows, err := self.appendWorkflowsSlices(rows)
	if err != nil {
		return nil, err
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
	sqlStatement := `SELECT * FROM workflows WHERE id = ($1)`
	var retrievedWorkflow entities.Workflow
	var actionParamByte []byte
	var reactionParamByte []byte

	row := self.db.QueryRow(sqlStatement, workflowId)

	err := row.Scan(&retrievedWorkflow.Id, &retrievedWorkflow.Name, &retrievedWorkflow.OwnerId, &retrievedWorkflow.ReactionId,
		&retrievedWorkflow.ActionId, &retrievedWorkflow.IsActivated, &retrievedWorkflow.CreatedAt, &actionParamByte, &reactionParamByte)
	if err != nil {
		return retrievedWorkflow, err
	}

	retrievedWorkflow, errUnmarshal := self.unmarshalActionReactionParam(actionParamByte, reactionParamByte, retrievedWorkflow)
	if errUnmarshal != nil {
		return retrievedWorkflow, errUnmarshal
	}
	return retrievedWorkflow, nil
}

func (self *WorkflowRepository) unmarshalActionReactionParam(actionParamByte []byte, reactionParamByte []byte, retrievedWorkflow entities.Workflow) (entities.Workflow, error) {
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
	return retrievedWorkflow, nil
}

func (self *WorkflowRepository) appendWorkflowsSlices(rows *sql.Rows) ([]entities.Workflow, error) {
	var retrievedWorkflows []entities.Workflow

	for rows.Next() {
		var workflow entities.Workflow
		var actionParamByte []byte
		var reactionParamByte []byte

		err := rows.Scan(&workflow.Id, &workflow.Name, &workflow.OwnerId, &workflow.ActionId,
			&workflow.ReactionId, &workflow.IsActivated, &workflow.CreatedAt, &actionParamByte, &reactionParamByte)
		if err != nil {
			return nil, err
		}

		workflow, errUnmarshal := self.unmarshalActionReactionParam(actionParamByte, reactionParamByte, workflow)
		if errUnmarshal != nil {
			return retrievedWorkflows, errUnmarshal
		}

		retrievedWorkflows = append(retrievedWorkflows, workflow)
	}
	return retrievedWorkflows, nil
}

func (self *WorkflowRepository) DeleteWorkflow(workflowId string, ownerId string) error {
	sqlStatement := `DELETE FROM workflows WHERE id = ($1) AND ownerid = ($2)`

	res, err := self.db.Exec(sqlStatement, workflowId, ownerId)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("Could not delete workflow")
	}
	return nil
}

func (self *WorkflowRepository) DeleteWorkflowByOwnerId(ownerId string) error {
	sqlStatement := `DELETE FROM workflows WHERE ownerid = ($1)`

	_, err := self.db.Exec(sqlStatement, ownerId)
	if err != nil {
		return err
	}
	return nil
}
