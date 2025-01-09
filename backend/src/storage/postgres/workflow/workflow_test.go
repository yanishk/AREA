package workflow_repository

import (
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
