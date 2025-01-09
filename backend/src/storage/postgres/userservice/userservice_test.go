package userservice_repository

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

	repo := NewUserServiceRepository(db)

	sqlStatement := `INSERT INTO userservices \(userid, token, tokenrefresh, expiry, serviceid\) VALUES \(\$1, \$2, \$3, \$4, \$5\)`
	mock.ExpectExec(sqlStatement).
		WithArgs("userid", "token", "tokenrefresh", "expiry", "serviceid").
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.CreateUserService("userid", "token", "tokenrefresh", "expiry", "serviceid")

	assert.NoError(test, err)

	err = mock.ExpectationsWereMet()
	if err != nil {
		test.Errorf("Expectation fail")
	}
}
