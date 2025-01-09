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
