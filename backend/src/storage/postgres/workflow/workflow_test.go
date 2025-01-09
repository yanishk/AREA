package workflow_repository

import (
	"backend/src/entities"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCreateWorkflow(test *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		test.Fatalf("Mock DB fail")
	}
	defer db.Close()

	repo := NewWorkflowRepository(db)

	sqlStatement := `INSERT INTO workflows \(name, ownerid, actionid, reactionid, isactivated, actionparam, reactionparam, actiondata\) VALUES \(\$1, \$2, \$3, \$4, \$5, \$6, \$7, \$8\)`
	mock.ExpectExec(sqlStatement).
		WithArgs("workflow", "owner", "action", "reaction", true, []byte("{\"key\":\"value\"}"), []byte("{\"key\":\"value\"}"), []byte("{\"key\":\"value\"}")).
		WillReturnResult(sqlmock.NewResult(1, 1))

	actionParam := map[string]interface{}{"key": "value"}
	reactionParam := map[string]interface{}{"key": "value"}
	actionData := map[string]interface{}{"key": "value"}

	err = repo.CreateWorkflow("workflow", "owner", "action", "reaction", actionParam, reactionParam, actionData)

	assert.NoError(test, err)

	err = mock.ExpectationsWereMet()
	if err != nil {
		test.Errorf("Expectation fail")
	}
}

func TestFindWorkflowById(test *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		test.Fatalf("Mock DB fail")
	}
	defer db.Close()

	repo := NewWorkflowRepository(db)

	id := "1234"
	sqlStatement := `SELECT \* FROM workflows WHERE id = \(\$1\)`

	rows := sqlmock.NewRows([]string{
		"id", "name", "ownerid", "actionid", "reactionid", "isactivated", "createdat",
		"actionparam", "reactionparam", "actiondata",
	}).AddRow(id, "workflow", "owner", "action", "reaction", true, "2025-01-01T00:00:00Z",
		[]byte(`{"key":"value"}`), []byte(`{"key":"value"}`), []byte(`{"key":"value"}`),
	)

	mock.ExpectQuery(sqlStatement).
		WithArgs(id).
		WillReturnRows(rows)

	workflow, err := repo.FindWorkflowById(id)

	assert.NoError(test, err)
	assert.Equal(test, id, workflow.Id)
	assert.Equal(test, "workflow", workflow.Name)
	assert.Equal(test, "owner", workflow.OwnerId)
	assert.Equal(test, "action", workflow.ActionId)
	assert.Equal(test, "reaction", workflow.ReactionId)
	assert.Equal(test, true, workflow.IsActivated)

	err = mock.ExpectationsWereMet()
	if err != nil {
		test.Errorf("Expectation fail")
	}
}

func TestFindWorkflowsByActionId(test *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		test.Fatalf("Mock DB fail")
	}
	defer db.Close()

	repo := NewWorkflowRepository(db)

	idAction := "1234"
	sqlStatement := `SELECT \* FROM workflows WHERE actionid = \(\$1\)`

	rows := sqlmock.NewRows([]string{
		"id", "name", "ownerid", "actionid", "reactionid", "isactivated", "createdat",
		"actionparam", "reactionparam", "actiondata",
	}).AddRow(idAction, "workflow", "owner", "action", "reaction", true, "2025-01-01T00:00:00Z",
		[]byte(`{"key":"value"}`), []byte(`{"key":"value"}`), []byte(`{"key":"value"}`),
	)

	mock.ExpectQuery(sqlStatement).
		WithArgs(idAction).
		WillReturnRows(rows)

	workflow, err := repo.FindWorkflowsByActionId(idAction)

	assert.NoError(test, err)
	assert.Equal(test, idAction, workflow[0].Id)
	assert.Equal(test, "workflow", workflow[0].Name)
	assert.Equal(test, "owner", workflow[0].OwnerId)
	assert.Equal(test, "action", workflow[0].ActionId)
	assert.Equal(test, "reaction", workflow[0].ReactionId)
	assert.Equal(test, true, workflow[0].IsActivated)

	err = mock.ExpectationsWereMet()
	if err != nil {
		test.Errorf("Expectation fail")
	}
}

func TestUpdateWorkflow(test *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		test.Fatalf("Mock DB fail")
	}
	defer db.Close()

	repo := NewWorkflowRepository(db)

	sqlStatement := `UPDATE workflows SET name = \(\$1\), actionid = \(\$2\), reactionid = \(\$3\), isactivated = \(\$4\), actionparam = \(\$5\), reactionparam = \(\$6\), actiondata = \(\$7\) WHERE id = \(\$8\)`

	var workflowToUpdate entities.Workflow
	workflowToUpdate.Id = "1234"
	workflowToUpdate.Name = "name"
	workflowToUpdate.IsActivated = false
	workflowToUpdate.OwnerId = "owner"
	workflowToUpdate.ActionId = "action"
	workflowToUpdate.ReactionId = "reaction"
	workflowToUpdate.ActionParam = map[string]interface{}{"key": "value"}
	workflowToUpdate.ReactionParam = map[string]interface{}{"key": "value"}
	workflowToUpdate.ActionData = map[string]interface{}{"key": "value"}

	actionParamJson, reactionParamJson, actionDataJson, err := marshalWorkflowParameters(workflowToUpdate.ActionParam, workflowToUpdate.ReactionParam, workflowToUpdate.ActionData)
	if err != nil {
		test.Fatalf("Failed to marshal parameters: %v", err)
	}

	mock.ExpectExec(sqlStatement).
		WithArgs("name", "action", "reaction", false, actionParamJson, reactionParamJson, actionDataJson, "1234").
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.UpdateWorkflow("1234", workflowToUpdate)

	assert.NoError(test, err)

	err = mock.ExpectationsWereMet()
	if err != nil {
		test.Errorf("Expectation fail")
	}
}
