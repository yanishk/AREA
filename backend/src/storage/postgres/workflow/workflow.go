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

func marshalWorkflowParameters(actionParam, reactionParam, actionData map[string]interface{}) ([]byte, []byte, []byte, error) {
	actionParamJson, err := json.Marshal(actionParam)
	if err != nil {
		return nil, nil, nil, err
	}

	reactionParamJson, err := json.Marshal(reactionParam)
	if err != nil {
		return nil, nil, nil, err
	}

	actionDataJson, err := json.Marshal(actionData)
	if err != nil {
		return nil, nil, nil, err
	}
	return actionParamJson, reactionParamJson, actionDataJson, nil
}

func unmarshalWorkflowParameters(actionParamBytes, reactionParamBytes, actionDataBytes []byte, workflow entities.Workflow) (entities.Workflow, error) {
	if len(actionParamBytes) > 0 {
		errActionBytes := json.Unmarshal(actionParamBytes, &workflow.ActionParam)
		if errActionBytes != nil {
			return workflow, errActionBytes
		}
	}
	if len(reactionParamBytes) > 0 {
		errReactionBytes := json.Unmarshal(reactionParamBytes, &workflow.ReactionParam)
		if errReactionBytes != nil {
			return workflow, errReactionBytes
		}
	}
	if len(actionDataBytes) > 0 {
		errDataBytes := json.Unmarshal(actionDataBytes, &workflow.ActionData)
		if errDataBytes != nil {
			return workflow, errDataBytes
		}
	}
	return workflow, nil
}

func appendWorkflowsSlices(rows *sql.Rows) ([]entities.Workflow, error) {
	var workflows []entities.Workflow

	for rows.Next() {
		var workflow entities.Workflow
		var actionParamBytes []byte
		var reactionParamBytes []byte
		var actionDataBytes []byte

		err := rows.Scan(&workflow.Id, &workflow.Name, &workflow.OwnerId, &workflow.ActionId,
			&workflow.ReactionId, &workflow.IsActivated, &workflow.CreatedAt, &actionParamBytes, &reactionParamBytes, &actionDataBytes)
		if err != nil {
			return nil, err
		}

		workflow, errUnmarshal := unmarshalWorkflowParameters(actionParamBytes, reactionParamBytes, actionDataBytes, workflow)
		if errUnmarshal != nil {
			return workflows, errUnmarshal
		}

		workflows = append(workflows, workflow)
	}
	return workflows, nil
}

func (self *WorkflowRepository) CreateWorkflow(name, ownerId, actionId, reactionId string, actionParam, reactionParam, actionData map[string]interface{}) error {
	sqlStatement := `INSERT INTO workflows (name, ownerid, actionid, reactionid, isactivated, actionparam, reactionparam, actiondata) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	actionParamJson, reactionParamJson, actionDataJson, err := marshalWorkflowParameters(actionParam, reactionParam, actionData)
	if err != nil {
		return err
	}

	_, err = self.db.Exec(sqlStatement, name, ownerId, actionId, reactionId, true, actionParamJson, reactionParamJson, actionDataJson)
	if err != nil {
		return err
	}
	return nil
}

func (self *WorkflowRepository) FindWorkflowById(id string) (entities.Workflow, error) {
	sqlStatement := `SELECT * FROM workflows WHERE id = ($1)`
	var workflow entities.Workflow
	var actionParamBytes, reactionParamBytes, actionDataBytes []byte

	row := self.db.QueryRow(sqlStatement, id)

	err := row.Scan(&workflow.Id, &workflow.Name, &workflow.OwnerId, &workflow.ActionId,
		&workflow.ReactionId, &workflow.IsActivated, &workflow.CreatedAt, &actionParamBytes, &reactionParamBytes, &actionDataBytes)
	if err != nil {
		return workflow, err
	}

	workflow, errUnmarshal := unmarshalWorkflowParameters(actionParamBytes, reactionParamBytes, actionDataBytes, workflow)
	if errUnmarshal != nil {
		return workflow, errUnmarshal
	}
	return workflow, nil
}

func (self *WorkflowRepository) FindWorkflowsByActionId(actionId string) ([]entities.Workflow, error) {
	sqlStatement := `SELECT * FROM workflows WHERE actionid = ($1)`

	rows, errQuery := self.db.Query(sqlStatement, actionId)
	if errQuery != nil {
		return nil, errQuery
	}
	defer rows.Close()

	workflows, err := appendWorkflowsSlices(rows)
	if err != nil {
		return nil, err
	}
	return workflows, nil
}

func (self *WorkflowRepository) FindWorkflowsByOwnerId(userId string) ([]entities.Workflow, error) {
	sqlStatement := `SELECT * FROM workflows WHERE ownerid = ($1)`

	rows, errQuery := self.db.Query(sqlStatement, userId)
	if errQuery != nil {
		return nil, errQuery
	}
	defer rows.Close()

	workflows, err := appendWorkflowsSlices(rows)
	if err != nil {
		return nil, err
	}

	return workflows, nil
}

func (self *WorkflowRepository) UpdateWorkflow(id string, updatedWorkflow entities.Workflow) error {
	sqlStatement := `UPDATE workflows SET name = ($1), actionid = ($2), reactionid = ($3), isactivated = ($4), actionparam = ($5), reactionparam = ($6), actiondata = ($7) WHERE id = ($8)`

	actionParamJson, reactionParamJson, actionDataJson, err := marshalWorkflowParameters(updatedWorkflow.ActionParam, updatedWorkflow.ReactionParam, updatedWorkflow.ActionData)
	if err != nil {
		return err
	}

	_, err = self.db.Exec(sqlStatement, updatedWorkflow.Name, updatedWorkflow.ActionId, updatedWorkflow.ReactionId, updatedWorkflow.IsActivated, actionParamJson, reactionParamJson, actionDataJson, id)
	if err != nil {
		return err
	}
	return nil
}

func (self *WorkflowRepository) DeleteWorkflow(id, ownerId string) error {
	sqlStatement := `DELETE FROM workflows WHERE id = ($1) AND ownerid = ($2)`

	res, err := self.db.Exec(sqlStatement, id, ownerId)
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
