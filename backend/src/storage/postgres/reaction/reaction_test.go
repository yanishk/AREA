package reaction_repository

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func createMockDb(test *testing.T) (*sql.DB, sqlmock.Sqlmock, *ReactionRepository) {
	db, mock, err := sqlmock.New()
	if err != nil {
		test.Fatalf("Mock DB fail")
	}
	repo := NewReactionRepository(db)
	return db, mock, repo
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
