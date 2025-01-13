package service_repository

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"backend/src/entities"
)

func createMockDb(test *testing.T) (*sql.DB, sqlmock.Sqlmock, *ServiceRepository) {
	db, mock, err := sqlmock.New()
	if err != nil {
		test.Fatalf("Mock DB fail")
	}
	repo := NewServiceRepository(db)
	return db, mock, repo
}

func assertService(test *testing.T, service entities.Service, id, name, color, logo, description string, hasActions, hasReactions, isAuthNeeded bool) {
	assert.Equal(test, id, service.Id)
	assert.Equal(test, name, service.Name)
	assert.Equal(test, color, service.Color)
	assert.Equal(test, logo, service.Logo)
	assert.Equal(test, hasActions, service.HasActions)
	assert.Equal(test, hasReactions, service.HasReactions)
	assert.Equal(test, isAuthNeeded, service.IsAuthNeeded)
	assert.Equal(test, description, service.Description)
}

func TestCreateService(test *testing.T) {
	db, mock, repo := createMockDb(test)
	defer db.Close()

	test.Run("Successful", func(test *testing.T) {
		sqlStatement := `INSERT INTO services \(name, color, logo\) VALUES \(\$1, \$2, \$3\)`
		mock.ExpectExec(sqlStatement).
			WithArgs("name", "color", "logo").
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.CreateService("name", "color", "logo")

		assert.NoError(test, err)

		err = mock.ExpectationsWereMet()
		if err != nil {
			test.Errorf("Expectation fail")
		}
	})

	test.Run("Service already exist", func(test *testing.T) {
		findSqlStatement := `SELECT \* FROM services WHERE name = \(\$1\)`
		mockRow := sqlmock.NewRows([]string{"id", "name", "color", "logo", "hasactions", "hasreactions", "isauthneeded", "description"}).
			AddRow("id", "name", "color", "logo", false, false, false, "description")

		mock.ExpectQuery(findSqlStatement).
			WithArgs("name").
			WillReturnRows(mockRow)

		sqlStatement := `INSERT INTO services \(name, color, logo\) VALUES \(\$1, \$2, \$3\)`
		mock.ExpectExec(sqlStatement).
			WithArgs("name", "color", "logo").
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.CreateService("name", "color", "logo")

		err = mock.ExpectationsWereMet()
		if err == nil {
			test.Errorf("Expectation fail")
		}
	})
}

func TestFindServiceById(test *testing.T) {
	db, mock, repo := createMockDb(test)
	defer db.Close()

	test.Run("Successful", func(test *testing.T) {

		sqlStatement := `SELECT \* FROM services WHERE id = \(\$1\)`
		mockRow := sqlmock.NewRows([]string{"id", "name", "color", "logo", "hasactions", "hasreactions", "isauthneeded", "description"}).
			AddRow("id", "name", "color", "logo", false, false, false, "description")

		mock.ExpectQuery(sqlStatement).
			WithArgs("id").
			WillReturnRows(mockRow)

		service, err := repo.FindServiceById("id")

		assert.NoError(test, err)
		assertService(test, service, "id", "name", "color", "logo", "description", false, false, false)

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

func TestFindAllServices(test *testing.T) {
	db, mock, repo := createMockDb(test)
	defer db.Close()

	sqlStatement := `SELECT \* FROM services`
	mockRow := sqlmock.NewRows([]string{"id", "name", "color", "logo", "hasactions", "hasreactions", "isauthneeded", "description"}).
		AddRow("id", "name", "color", "logo", false, false, false, "description")

	mock.ExpectQuery(sqlStatement).
		WillReturnRows(mockRow)

	services, err := repo.FindAllServices()

	assert.NoError(test, err)
	assertService(test, services[0], "id", "name", "color", "logo", "description", false, false, false)

	err = mock.ExpectationsWereMet()
	if err != nil {
		test.Errorf("Expectation fail")
	}
}

func TestFindActionsServices(test *testing.T) {
	db, mock, repo := createMockDb(test)
	defer db.Close()

	sqlStatement := `SELECT \* FROM services WHERE hasactions = true`
	mockRow := sqlmock.NewRows([]string{"id", "name", "color", "logo", "hasactions", "hasreactions", "isauthneeded", "description"}).
		AddRow("id", "name", "color", "logo", false, false, false, "description")

	mock.ExpectQuery(sqlStatement).
		WillReturnRows(mockRow)

	services, err := repo.FindActionsServices()

	assert.NoError(test, err)
	assertService(test, services[0], "id", "name", "color", "logo", "description", false, false, false)

	err = mock.ExpectationsWereMet()
	if err != nil {
		test.Errorf("Expectation fail")
	}
}

func TestFindReactionsServices(test *testing.T) {
	db, mock, repo := createMockDb(test)
	defer db.Close()

	sqlStatement := `SELECT \* FROM services WHERE hasreactions = true`
	mockRow := sqlmock.NewRows([]string{"id", "name", "color", "logo", "hasactions", "hasreactions", "isauthneeded", "description"}).
		AddRow("id", "name", "color", "logo", false, false, false, "description")

	mock.ExpectQuery(sqlStatement).
		WillReturnRows(mockRow)

	services, err := repo.FindReactionsServices()

	assert.NoError(test, err)
	assertService(test, services[0], "id", "name", "color", "logo", "description", false, false, false)

	err = mock.ExpectationsWereMet()
	if err != nil {
		test.Errorf("Expectation fail")
	}
}
