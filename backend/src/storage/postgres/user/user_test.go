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
