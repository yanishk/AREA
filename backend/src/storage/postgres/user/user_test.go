package user_repository

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCreateUserService(test *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		test.Fatalf("Mock DB fail")
	}
	defer db.Close()

	repo := NewUserRepository(db)

	sqlStatement := `INSERT INTO users \(email, password, connectiontype\) VALUES \(\$1, \$2, \$3\)`
	mock.ExpectExec(sqlStatement).
		WithArgs("email", "password", "connectiontype").
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.CreateUser("email", "password", "connectiontype")

	assert.NoError(test, err)

	err = mock.ExpectationsWereMet()
	if err != nil {
		test.Errorf("Expectation fail")
	}
}

func TestFindUserByEmail(test *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		test.Fatalf("Mock DB fail")
	}
	defer db.Close()

	repo := NewUserRepository(db)

	sqlStatement := `SELECT \* FROM users WHERE email = \(\$1\) AND connectiontype = \(\$2\)`
	mockRow := sqlmock.NewRows([]string{"email", "password", "id", "createdat", "timezone", "connectiontype"}).
		AddRow("email", "password", "id", "createdat", "timezone", "connectiontype")

	mock.ExpectQuery(sqlStatement).
		WithArgs("email", "connectiontype").
		WillReturnRows(mockRow)

	user, err := repo.FindUserByEmail("email", "connectiontype")

	assert.NoError(test, err)
	assert.Equal(test, "email", user.Email)
	assert.Equal(test, "connectiontype", user.ConnectionType)
	assert.Equal(test, "id", user.Id)
	assert.Equal(test, "createdat", user.CreatedAt)
	assert.Equal(test, "timezone", user.Timezone)
	assert.Equal(test, "connectiontype", user.ConnectionType)

	err = mock.ExpectationsWereMet()
	if err != nil {
		test.Errorf("Expectation fail")
	}
}

func TestFindUserById(test *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		test.Fatalf("Mock DB fail")
	}
	defer db.Close()

	repo := NewUserRepository(db)

	sqlStatement := `SELECT \* FROM users WHERE id = \(\$1\)`
	mockRow := sqlmock.NewRows([]string{"email", "password", "id", "createdat", "timezone", "connectiontype"}).
		AddRow("email", "password", "id", "createdat", "timezone", "connectiontype")

	mock.ExpectQuery(sqlStatement).
		WithArgs("id").
		WillReturnRows(mockRow)

	user, err := repo.FindUserById("id")

	assert.NoError(test, err)
	assert.Equal(test, "email", user.Email)
	assert.Equal(test, "connectiontype", user.ConnectionType)
	assert.Equal(test, "id", user.Id)
	assert.Equal(test, "createdat", user.CreatedAt)
	assert.Equal(test, "timezone", user.Timezone)

	err = mock.ExpectationsWereMet()
	if err != nil {
		test.Errorf("Expectation fail")
	}
}

func TestUpdateUser(test *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		test.Fatalf("Mock DB fail")
	}
	defer db.Close()

	repo := NewUserRepository(db)

	test.Run("Successful", func(test *testing.T) {
		findSqlStatement := `SELECT \* FROM users WHERE email = \(\$1\) AND connectiontype = \(\$2\)`
		mockRow := sqlmock.NewRows([]string{"email", "password", "id", "createdat", "timezone", "connectiontype"}).
			AddRow("email", "password", "id", "createdat", "timezone", "connectiontype")

		mock.ExpectQuery(findSqlStatement).
			WithArgs("email", "connectiontype").
			WillReturnRows(mockRow)

		sqlStatement := `UPDATE users SET password = \(\$1\) WHERE email = \(\$2\)`
		mock.ExpectExec(sqlStatement).
			WithArgs("newpassword", "email").
			WillReturnResult(sqlmock.NewResult(1, 1))

		err = repo.UpdateUser("email", "newpassword", "connectiontype")

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

		err = repo.UpdateUser("email", "newpassword", "connectiontype")

		assert.NoError(test, err)

		err = mock.ExpectationsWereMet()
		if err != nil {
			test.Errorf("Expectation fail")
		}
	})
}
