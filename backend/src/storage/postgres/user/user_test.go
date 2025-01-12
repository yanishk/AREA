package user_repository

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"backend/src/entities"
)

func createMockDb(test *testing.T) (*sql.DB, sqlmock.Sqlmock, *UserRepository) {
	db, mock, err := sqlmock.New()
	if err != nil {
		test.Fatalf("Mock DB fail")
	}
	repo := NewUserRepository(db)
	return db, mock, repo
}

func assertUser(test *testing.T, user entities.User, email, connectionType, id, createdAt, timezone, password string) {
	assert.Equal(test, email, user.Email)
	assert.Equal(test, connectionType, user.ConnectionType)
	assert.Equal(test, id, user.Id)
	assert.Equal(test, createdAt, user.CreatedAt)
	assert.Equal(test, timezone, user.Timezone)
	assert.Equal(test, password, user.Password)
}

func TestCreateUser(test *testing.T) {
	db, mock, repo := createMockDb(test)
	defer db.Close()

	test.Run("Successful", func(test *testing.T) {
		sqlStatement := `INSERT INTO users \(email, password, connectiontype\) VALUES \(\$1, \$2, \$3\)`
		mock.ExpectExec(sqlStatement).
			WithArgs("email", "password", "connectiontype").
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.CreateUser("email", "password", "connectiontype")

		assert.NoError(test, err)

		err = mock.ExpectationsWereMet()
		if err != nil {
			test.Errorf("Expectation fail")
		}
	})

	test.Run("User already exist", func(test *testing.T) {
		findSqlStatement := `SELECT \* FROM users WHERE email = \(\$1\) AND connectiontype = \(\$2\)`
		mockRow := sqlmock.NewRows([]string{"email", "password", "id", "createdat", "timezone", "connectiontype"}).
			AddRow("email", "password", "id", "createdat", "timezone", "connectiontype")

		mock.ExpectQuery(findSqlStatement).
			WithArgs("email", "connectiontype").
			WillReturnRows(mockRow)

		sqlStatement := `DELETE FROM users WHERE email = \(\$1\)`
		mock.ExpectExec(sqlStatement).
			WithArgs("email").
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.CreateUser("email", "password", "connectiontype")

		err = mock.ExpectationsWereMet()
		if err == nil {
			test.Errorf("Expectation fail")
		}
	})
}

func TestFindUserByEmail(test *testing.T) {
	db, mock, repo := createMockDb(test)
	defer db.Close()

	sqlStatement := `SELECT \* FROM users WHERE email = \(\$1\) AND connectiontype = \(\$2\)`
	mockRow := sqlmock.NewRows([]string{"email", "password", "id", "createdat", "timezone", "connectiontype"}).
		AddRow("email", "password", "id", "createdat", "timezone", "connectiontype")

	mock.ExpectQuery(sqlStatement).
		WithArgs("email", "connectiontype").
		WillReturnRows(mockRow)

	user, err := repo.FindUserByEmail("email", "connectiontype")

	assert.NoError(test, err)
	assertUser(test, user, "email", "connectiontype", "id", "createdat", "timezone", "password")

	err = mock.ExpectationsWereMet()
	if err != nil {
		test.Errorf("Expectation fail")
	}
}

func TestFindUserById(test *testing.T) {
	db, mock, repo := createMockDb(test)
	defer db.Close()

	sqlStatement := `SELECT \* FROM users WHERE id = \(\$1\)`
	mockRow := sqlmock.NewRows([]string{"email", "password", "id", "createdat", "timezone", "connectiontype"}).
		AddRow("email", "password", "id", "createdat", "timezone", "connectiontype")

	mock.ExpectQuery(sqlStatement).
		WithArgs("id").
		WillReturnRows(mockRow)

	user, err := repo.FindUserById("id")

	assert.NoError(test, err)
	assertUser(test, user, "email", "connectiontype", "id", "createdat", "timezone", "password")

	err = mock.ExpectationsWereMet()
	if err != nil {
		test.Errorf("Expectation fail")
	}
}

func TestUpdateUser(test *testing.T) {
	db, mock, repo := createMockDb(test)
	defer db.Close()

	test.Run("Successful", func(test *testing.T) {
		findSqlStatement := `SELECT \* FROM users WHERE email = \(\$1\) AND connectiontype = \(\$2\)`
		mockRow := sqlmock.NewRows([]string{"email", "password", "id", "createdat", "timezone", "connectiontype"}).
			AddRow("email", "password", "id", "createdat", "timezone", "connectiontype")

		mock.ExpectQuery(findSqlStatement).
			WithArgs("email", "connectiontype").
			WillReturnRows(mockRow)

		sqlStatement := `UPDATE users SET password = \(\$1\) WHERE email = \(\$2\) AND connectiontype = \(\$3\)`
		mock.ExpectExec(sqlStatement).
			WithArgs("newpassword", "email", "connectiontype").
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.UpdateUser("email", "newpassword", "connectiontype")

		assert.NoError(test, err)

		err = mock.ExpectationsWereMet()
		if err != nil {
			test.Errorf("Expectation fail")
		}
	})

	test.Run("User doesn't exist", func(test *testing.T) {
		sqlStatement := `UPDATE users SET password = \(\$1\) WHERE email = \(\$2\)`
		mock.ExpectExec(sqlStatement).
			WithArgs("newpassword", "email").
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.UpdateUser("email", "newpassword", "connectiontype")

		err = mock.ExpectationsWereMet()
		if err == nil {
			test.Errorf("Expectation fail")
		}
	})
}

func TestDeleteUser(test *testing.T) {
	db, mock, repo := createMockDb(test)
	defer db.Close()

	test.Run("Successful", func(test *testing.T) {
		findSqlStatement := `SELECT \* FROM users WHERE email = \(\$1\) AND connectiontype = \(\$2\)`
		mockRow := sqlmock.NewRows([]string{"email", "password", "id", "createdat", "timezone", "connectiontype"}).
			AddRow("email", "password", "id", "createdat", "timezone", "connectiontype")

		mock.ExpectQuery(findSqlStatement).
			WithArgs("email", "connectiontype").
			WillReturnRows(mockRow)

		sqlStatement := `DELETE FROM users WHERE email = \(\$1\)`
		mock.ExpectExec(sqlStatement).
			WithArgs("email").
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.DeleteUser("email", "connectiontype")

		assert.NoError(test, err)

		err = mock.ExpectationsWereMet()
		if err != nil {
			test.Errorf("Expectation fail")
		}
	})

	test.Run("User doesn't exist", func(test *testing.T) {
		sqlStatement := `DELETE FROM users WHERE email = \(\$1\)`
		mock.ExpectExec(sqlStatement).
			WithArgs("email").
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.DeleteUser("email", "connectiontype")

		err = mock.ExpectationsWereMet()
		if err == nil {
			test.Errorf("Expectation fail")
		}
	})
}
