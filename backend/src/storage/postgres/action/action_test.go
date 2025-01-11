package action_repository

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"backend/src/entities"
)

func createMockDb(test *testing.T) (*sql.DB, sqlmock.Sqlmock, *ActionRepository) {
	db, mock, err := sqlmock.New()
	if err != nil {
		test.Fatalf("Mock DB fail")
	}
	repo := NewActionRepository(db)
	return db, mock, repo
}

func assertAction(test *testing.T, action entities.Action, id, serviceId, name, description string, nbParam int) {
	assert.Equal(test, id, action.Id)
	assert.Equal(test, name, action.Name)
	assert.Equal(test, description, action.Description)
	assert.Equal(test, serviceId, action.ServiceId)
	assert.Equal(test, nbParam, action.NbParam)
}

func TestCreateAction(test *testing.T) {
	db, mock, repo := createMockDb(test)
	defer db.Close()

	test.Run("Successful", func(test *testing.T) {
		sqlStatement := `INSERT INTO actions \(name, description, serviceid, nbparam\) VALUES \(\$1, \$2, \$3, \$4\)`
		mock.ExpectExec(sqlStatement).
			WithArgs("name", "description", "serviceid", 3).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.CreateAction("name", "description", "serviceid", 3)

		assert.NoError(test, err)

		err = mock.ExpectationsWereMet()
		if err != nil {
			test.Errorf("Expectation fail")
		}
	})

	test.Run("Reaction already exist", func(test *testing.T) {
		findSqlStatement := `SELECT \* FROM actions WHERE name = \(\$1\)`
		mockRow := sqlmock.NewRows([]string{"id", "name", "description", "serviceid", "nbparam", "parameters"}).
			AddRow("id", "name", "description", "serviceid", 3, nil)

		mock.ExpectQuery(findSqlStatement).
			WithArgs("name").
			WillReturnRows(mockRow)

		sqlStatement := `INSERT INTO actions \(name, description, serviceid, nbparam\) VALUES \(\$1, \$2, \$3, \$4\)`
		mock.ExpectExec(sqlStatement).
			WithArgs("name", "description", "serviceid", 3).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.CreateAction("name", "description", "serviceid", 3)

		err = mock.ExpectationsWereMet()
		if err == nil {
			test.Errorf("Expectation fail")
		}
	})
}

func TestFindActionById(test *testing.T) {
	db, mock, repo := createMockDb(test)
	defer db.Close()

	test.Run("Successful", func(test *testing.T) {
		sqlStatement := `SELECT \* FROM actions WHERE id = \(\$1\)`
		mockRow := sqlmock.NewRows([]string{"id", "serviceid", "name", "description", "nbparam", "parameters"}).
			AddRow("id", "serviceid", "name", "description", 3, nil)

		mock.ExpectQuery(sqlStatement).
			WithArgs("id").
			WillReturnRows(mockRow)

		action, err := repo.FindActionById("id")

		assert.NoError(test, err)
		assertAction(test, action, "id", "serviceid", "name", "description", 3)

		err = mock.ExpectationsWereMet()
		if err != nil {
			test.Errorf("Expectation fail")
		}
	})

	test.Run("Action not found", func(test *testing.T) {

		sqlStatement := `SELECT \* FROM actions WHERE id = \(\$1\)`
		mock.ExpectQuery(sqlStatement).
			WithArgs("id")

		_, err := repo.FindActionById("id")
		assert.Error(test, err)

		err = mock.ExpectationsWereMet()
		if err != nil {
			test.Errorf("Expectation fail")
		}
	})
}

func TestFindActionByServiceId(test *testing.T) {
	db, mock, repo := createMockDb(test)
	defer db.Close()

	test.Run("Successful", func(test *testing.T) {
		sqlStatement := `SELECT \* FROM actions WHERE serviceid = \(\$1\)`
		mockRow := sqlmock.NewRows([]string{"id", "serviceid", "name", "description", "nbparam", "parameters"}).
			AddRow("id", "serviceid", "name", "description", 3, nil)

		mock.ExpectQuery(sqlStatement).
			WithArgs("serviceid").
			WillReturnRows(mockRow)

		reaction, err := repo.FindActionsByServiceId("serviceid")

		assert.NoError(test, err)
		assertAction(test, reaction[0], "id", "serviceid", "name", "description", 3)

		err = mock.ExpectationsWereMet()
		if err != nil {
			test.Errorf("Expectation fail")
		}
	})

	test.Run("Action not found", func(test *testing.T) {

		sqlStatement := `SELECT \* FROM actions WHERE serviceid = \(\$1\)`
		mock.ExpectQuery(sqlStatement).
			WithArgs("serviceid")

		_, err := repo.FindActionsByServiceId("serviceid")
		assert.Error(test, err)

		err = mock.ExpectationsWereMet()
		if err != nil {
			test.Errorf("Expectation fail")
		}
	})
}
