package workflow_repository

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"backend/src/entities"
)

func createMockDb(test *testing.T) (*sql.DB, sqlmock.Sqlmock, *WorkflowRepository) {
	db, mock, err := sqlmock.New()
	if err != nil {
		test.Fatalf("Mock DB fail")
	}
	repo := NewWorkflowRepository(db)
	return db, mock, repo
}

func assertWorkflow(test *testing.T, workflow entities.Workflow, id, name, ownerId, actionId, reactionId string, isActivated bool) {
	assert.Equal(test, id, workflow.Id)
	assert.Equal(test, name, workflow.Name)
	assert.Equal(test, ownerId, workflow.OwnerId)
	assert.Equal(test, actionId, workflow.ActionId)
	assert.Equal(test, reactionId, workflow.ReactionId)
	assert.Equal(test, isActivated, workflow.IsActivated)
}

func TestCreateWorkflow(test *testing.T) {
	db, mock, repo := createMockDb(test)
	defer db.Close()

	sqlStatement := `INSERT INTO workflows \(name, ownerid, actionid, reactionid, isactivated, actionparam, reactionparam, actiondata\) VALUES \(\$1, \$2, \$3, \$4, \$5, \$6, \$7, \$8\)`
	mock.ExpectExec(sqlStatement).
		WithArgs("workflow", "owner", "action", "reaction", true, []byte("{\"key\":\"value\"}"), []byte("{\"key\":\"value\"}"), []byte("{\"key\":\"value\"}")).
		WillReturnResult(sqlmock.NewResult(1, 1))

	actionParam := map[string]interface{}{"key": "value"}
	reactionParam := map[string]interface{}{"key": "value"}
	actionData := map[string]interface{}{"key": "value"}

	err := repo.CreateWorkflow("workflow", "owner", "action", "reaction", actionParam, reactionParam, actionData)

	assert.NoError(test, err)

	err = mock.ExpectationsWereMet()
	if err != nil {
		test.Errorf("Expectation fail")
	}
}

func TestFindWorkflowById(test *testing.T) {
	db, mock, repo := createMockDb(test)
	defer db.Close()

	id := "1234"
	sqlStatement := `SELECT \* FROM workflows WHERE id = \(\$1\)`

	rows := sqlmock.NewRows([]string{
		"id", "name", "ownerid", "actionid", "reactionid", "isactivated", "createdat",
		"actionparam", "reactionparam", "actiondata",
	}).AddRow(id, "workflow", "owner", "action", "reaction", true, "createdat",
		[]byte(`{"key":"value"}`), []byte(`{"key":"value"}`), []byte(`{"key":"value"}`),
	)

	mock.ExpectQuery(sqlStatement).
		WithArgs(id).
		WillReturnRows(rows)

	workflow, err := repo.FindWorkflowById(id)

	assert.NoError(test, err)
	assertWorkflow(test, workflow, id, "workflow", "owner", "action", "reaction", true)

	err = mock.ExpectationsWereMet()
	if err != nil {
		test.Errorf("Expectation fail")
	}
}

func TestFindWorkflowsByActionId(test *testing.T) {
	db, mock, repo := createMockDb(test)
	defer db.Close()

	idAction := "1234"
	sqlStatement := `SELECT \* FROM workflows WHERE actionid = \(\$1\)`

	rows := sqlmock.NewRows([]string{
		"id", "name", "ownerid", "actionid", "reactionid", "isactivated", "createdat",
		"actionparam", "reactionparam", "actiondata",
	}).AddRow(idAction, "workflow", "owner", "action", "reaction", true, "createdat",
		[]byte(`{"key":"value"}`), []byte(`{"key":"value"}`), []byte(`{"key":"value"}`),
	)

	mock.ExpectQuery(sqlStatement).
		WithArgs(idAction).
		WillReturnRows(rows)

	workflow, err := repo.FindWorkflowsByActionId(idAction)

	assert.NoError(test, err)
	assertWorkflow(test, workflow[0], idAction, "workflow", "owner", "action", "reaction", true)

	err = mock.ExpectationsWereMet()
	if err != nil {
		test.Errorf("Expectation fail")
	}
}

func TestUpdateWorkflow(test *testing.T) {
	db, mock, repo := createMockDb(test)
	defer db.Close()

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

func TestDeleteWorkflow(test *testing.T) {
	db, mock, repo := createMockDb(test)
	defer db.Close()

	sqlStatement := `DELETE FROM workflows WHERE id = \(\$1\) AND ownerid = \(\$2\)`
	mock.ExpectExec(sqlStatement).
		WithArgs("123", "owner").
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.DeleteWorkflow("123", "owner")

	assert.NoError(test, err)

	err = mock.ExpectationsWereMet()
	if err != nil {
		test.Errorf("Expectation fail")
	}
}

func TestDeleteWorkflowByOwnerId(test *testing.T) {
	db, mock, repo := createMockDb(test)
	defer db.Close()

	sqlStatement := `DELETE FROM workflows WHERE ownerid = \(\$1\)`
	mock.ExpectExec(sqlStatement).
		WithArgs("owner").
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.DeleteWorkflowByOwnerId("owner")

	assert.NoError(test, err)

	err = mock.ExpectationsWereMet()
	if err != nil {
		test.Errorf("Expectation fail")
	}
}
