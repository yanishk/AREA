package service_repository

import (
	"backend/src/entities"
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

	test.Run("Service already exist", func(test *testing.T) {
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

func assertService(test *testing.T, service entities.Service, id, name, color, logo string, hasActions, hasReactions, isAuthNeeded bool) {
	assert.Equal(test, id, service.Id)
	assert.Equal(test, name, service.Name)
	assert.Equal(test, color, service.Color)
	assert.Equal(test, logo, service.Logo)
	assert.Equal(test, hasActions, service.HasActions)
	assert.Equal(test, hasReactions, service.HasReactions)
	assert.Equal(test, isAuthNeeded, service.IsAuthNeeded)
}

func TestFindServiceById(test *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		test.Fatalf("Mock DB fail")
	}
	defer db.Close()

	repo := NewServiceRepository(db)

	test.Run("Successful", func(test *testing.T) {

		sqlStatement := `SELECT \* FROM services WHERE id = \(\$1\)`
		mockRow := sqlmock.NewRows([]string{"id", "name", "color", "logo", "hasactions", "hasreactions", "isauthneeded"}).
			AddRow("id", "name", "color", "logo", false, false, false)

		mock.ExpectQuery(sqlStatement).
			WithArgs("id").
			WillReturnRows(mockRow)

		service, err := repo.FindServiceById("id")

		assert.NoError(test, err)
		assertService(test, service, "id", "name", "color", "logo", false, false, false)

		err = mock.ExpectationsWereMet()
		if err != nil {
			test.Errorf("Expectation fail")
		}
	})

	test.Run("Service not found", func(test *testing.T) {

		sqlStatement := `SELECT \* FROM services WHERE id = \(\$1\)`
		mock.ExpectQuery(sqlStatement).
			WithArgs("id")

		_, err := repo.FindServiceById("id")
		assert.Error(test, err)

		err = mock.ExpectationsWereMet()
		if err != nil {
			test.Errorf("Expectation fail")
		}
	})
}
