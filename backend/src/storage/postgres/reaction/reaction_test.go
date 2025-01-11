package reaction_repository

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"backend/src/entities"
)

func createMockDb(test *testing.T) (*sql.DB, sqlmock.Sqlmock, *ReactionRepository) {
	db, mock, err := sqlmock.New()
	if err != nil {
		test.Fatalf("Mock DB fail")
	}
	repo := NewReactionRepository(db)
	return db, mock, repo
}

func assertReaction(test *testing.T, reaction entities.Reaction, id, serviceId, name, description string, nbParam int) {
	assert.Equal(test, id, reaction.Id)
	assert.Equal(test, name, reaction.Name)
	assert.Equal(test, description, reaction.Description)
	assert.Equal(test, serviceId, reaction.ServiceId)
	assert.Equal(test, nbParam, reaction.NbParam)
}

func TestCreateReaction(test *testing.T) {
	db, mock, repo := createMockDb(test)
	defer db.Close()

	test.Run("Successful", func(test *testing.T) {
		sqlStatement := `INSERT INTO reactions \(name, description, serviceid, nbparam\) VALUES \(\$1, \$2, \$3, \$4\)`
		mock.ExpectExec(sqlStatement).
			WithArgs("name", "description", "serviceid", 3).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.CreateReaction("name", "description", "serviceid", 3)

		assert.NoError(test, err)

		err = mock.ExpectationsWereMet()
		if err != nil {
			test.Errorf("Expectation fail")
		}
	})

	test.Run("Reaction already exist", func(test *testing.T) {
		findSqlStatement := `SELECT \* FROM reactions WHERE name = \(\$1\)`
		mockRow := sqlmock.NewRows([]string{"id", "name", "description", "serviceid", "nbparam", "parameters"}).
			AddRow("id", "name", "description", "serviceid", 3, nil)

		mock.ExpectQuery(findSqlStatement).
			WithArgs("name").
			WillReturnRows(mockRow)

		sqlStatement := `INSERT INTO reactions \(name, description, serviceid, nbparam\) VALUES \(\$1, \$2, \$3, \$4\)`
		mock.ExpectExec(sqlStatement).
			WithArgs("name", "description", "serviceid", 3).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.CreateReaction("name", "description", "serviceid", 3)

		err = mock.ExpectationsWereMet()
		if err == nil {
			test.Errorf("Expectation fail")
		}
	})
}

func TestFindReactionById(test *testing.T) {
	db, mock, repo := createMockDb(test)
	defer db.Close()

	test.Run("Successful", func(test *testing.T) {
		sqlStatement := `SELECT \* FROM reactions WHERE id = \(\$1\)`
		mockRow := sqlmock.NewRows([]string{"id", "serviceid", "name", "description", "nbparam", "parameters"}).
			AddRow("id", "serviceid", "name", "description", 3, nil)

		mock.ExpectQuery(sqlStatement).
			WithArgs("id").
			WillReturnRows(mockRow)

		reaction, err := repo.FindReactionById("id")

		assert.NoError(test, err)
		assertReaction(test, reaction, "id", "serviceid", "name", "description", 3)

		err = mock.ExpectationsWereMet()
		if err != nil {
			test.Errorf("Expectation fail")
		}
	})

	test.Run("Reaction not found", func(test *testing.T) {

		sqlStatement := `SELECT \* FROM reactions WHERE id = \(\$1\)`
		mock.ExpectQuery(sqlStatement).
			WithArgs("id")

		_, err := repo.FindReactionById("id")
		assert.Error(test, err)

		err = mock.ExpectationsWereMet()
		if err != nil {
			test.Errorf("Expectation fail")
		}
	})
}

func TestFindReactionByServiceId(test *testing.T) {
	db, mock, repo := createMockDb(test)
	defer db.Close()

	test.Run("Successful", func(test *testing.T) {
		sqlStatement := `SELECT \* FROM reactions WHERE serviceid = \(\$1\)`
		mockRow := sqlmock.NewRows([]string{"id", "serviceid", "name", "description", "nbparam", "parameters"}).
			AddRow("id", "serviceid", "name", "description", 3, nil)

		mock.ExpectQuery(sqlStatement).
			WithArgs("serviceid").
			WillReturnRows(mockRow)

		reaction, err := repo.FindReactionsByServiceId("serviceid")

		assert.NoError(test, err)
		assertReaction(test, reaction[0], "id", "serviceid", "name", "description", 3)

		err = mock.ExpectationsWereMet()
		if err != nil {
			test.Errorf("Expectation fail")
		}
	})

	test.Run("Reaction not found", func(test *testing.T) {

		sqlStatement := `SELECT \* FROM reactions WHERE serviceid = \(\$1\)`
		mock.ExpectQuery(sqlStatement).
			WithArgs("serviceid")

		_, err := repo.FindReactionsByServiceId("serviceid")
		assert.Error(test, err)

		err = mock.ExpectationsWereMet()
		if err != nil {
			test.Errorf("Expectation fail")
		}
	})
}
