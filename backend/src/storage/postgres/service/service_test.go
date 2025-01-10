package service_repository

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCreateService(test *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		test.Fatalf("Mock DB fail")
	}
	defer db.Close()

	repo := NewServiceRepository(db)

	test.Run("Successful", func(test *testing.T) {
		sqlStatement := `INSERT INTO services \(name, color, logo\) VALUES \(\$1, \$2, \$3\)`
		mock.ExpectExec(sqlStatement).
			WithArgs("name", "color", "logo").
			WillReturnResult(sqlmock.NewResult(1, 1))

		err = repo.CreateService("name", "color", "logo")

		assert.NoError(test, err)

		err = mock.ExpectationsWereMet()
		if err != nil {
			test.Errorf("Expectation fail")
		}
	})

	test.Run("User already exist", func(test *testing.T) {
		findSqlStatement := `SELECT \* FROM services WHERE name = \(\$1\)`
		mockRow := sqlmock.NewRows([]string{"id", "name", "color", "logo", "hasactions", "hasreactions", "isauthneeded"}).
			AddRow("id", "name", "color", "logo", false, false, false)

		mock.ExpectQuery(findSqlStatement).
			WithArgs("name").
			WillReturnRows(mockRow)

		sqlStatement := `INSERT INTO services \(name, color, logo\) VALUES \(\$1, \$2, \$3\)`
		mock.ExpectExec(sqlStatement).
			WithArgs("name", "color", "logo").
			WillReturnResult(sqlmock.NewResult(1, 1))

		err = repo.CreateService("name", "color", "logo")

		err = mock.ExpectationsWereMet()
		if err == nil {
			test.Errorf("Expectation fail")
		}
	})
}
