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

func TestFindUserServiceByServiceIdandUserId(test *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		test.Fatalf("Mock DB fail")
	}
	defer db.Close()

	repo := NewUserServiceRepository(db)

	sqlStatement := `SELECT \* FROM userservices WHERE userid = \(\$1\) AND serviceid = \(\$2\)`
	mockRow := sqlmock.NewRows([]string{"id", "userid", "accesstoken", "refreshtoken", "expirydate", "serviceid"}).
		AddRow("id", "userid", "accesstoken", "refreshtoken", "expirydate", "serviceid")

	mock.ExpectQuery(sqlStatement).
		WithArgs("userid", "serviceid").
		WillReturnRows(mockRow)

	userService, err := repo.FindUserServiceByServiceIdandUserId("userid", "serviceid")

	assert.NoError(test, err)
	assert.Equal(test, "id", userService.Id)
	assert.Equal(test, "userid", userService.UserId)
	assert.Equal(test, "accesstoken", userService.AccessToken)
	assert.Equal(test, "refreshtoken", userService.RefreshToken)
	assert.Equal(test, "expirydate", userService.ExpiryDate)
	assert.Equal(test, "serviceid", userService.ServiceId)

	err = mock.ExpectationsWereMet()
	if err != nil {
		test.Errorf("Expectation fail")
	}
}
