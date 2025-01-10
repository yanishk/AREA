package userservice_repository

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func createMockDb(test *testing.T) (*sql.DB, sqlmock.Sqlmock, *UserServiceRepository) {
	db, mock, err := sqlmock.New()
	if err != nil {
		test.Fatalf("Mock DB fail")
	}
	repo := NewUserServiceRepository(db)
	return db, mock, repo
}

func TestCreateUserService(test *testing.T) {
	db, mock, repo := createMockDb(test)
	defer db.Close()

	sqlStatement := `INSERT INTO userservices \(userid, token, tokenrefresh, expiry, serviceid\) VALUES \(\$1, \$2, \$3, \$4, \$5\)`
	mock.ExpectExec(sqlStatement).
		WithArgs("userid", "token", "tokenrefresh", "expiry", "serviceid").
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.CreateUserService("userid", "token", "tokenrefresh", "expiry", "serviceid")

	assert.NoError(test, err)

	err = mock.ExpectationsWereMet()
	if err != nil {
		test.Errorf("Expectation fail")
	}
}

func TestFindUserServiceByServiceIdandUserId(test *testing.T) {
	db, mock, repo := createMockDb(test)
	defer db.Close()

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

func TestUpdateUserServiceByServiceIdAndUserId(test *testing.T) {
	db, mock, repo := createMockDb(test)
	defer db.Close()

	test.Run("Successful", func(test *testing.T) {
		findSqlStatement := `SELECT \* FROM userservices WHERE userid = \(\$1\) AND serviceid = \(\$2\)`
		mockRow := sqlmock.NewRows([]string{"id", "userid", "accesstoken", "refreshtoken", "expirydate", "serviceid"}).
			AddRow("id", "userid", "oldtoken", "oldtokenrefresh", "expirydate", "serviceid")

		mock.ExpectQuery(findSqlStatement).
			WithArgs("userid", "serviceid").
			WillReturnRows(mockRow)

		updateSqlStatement := `UPDATE userservices SET token = \(\$1\), tokenrefresh = \(\$2\), expiry = \(\$3\) WHERE userid = \(\$4\) AND serviceid = \(\$5\)`
		mock.ExpectExec(updateSqlStatement).
			WithArgs("token", "tokenrefresh", "expiry", "userid", "serviceid").
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.UpdateUserServiceByServiceIdAndUserId("userid", "token", "tokenrefresh", "expiry", "serviceid")

		assert.NoError(test, err)

		err = mock.ExpectationsWereMet()
		if err != nil {
			test.Errorf("Expectation fail")
		}
	})

	test.Run("Service doesn't exist for user", func(test *testing.T) {
		updateSqlStatement := `UPDATE userservices SET token = \(\$1\), tokenrefresh = \(\$2\), expiry = \(\$3\) WHERE userid = \(\$4\) AND serviceid = \(\$5\)`
		mock.ExpectExec(updateSqlStatement).
			WithArgs("token", "tokenrefresh", "expiry", "userid", "serviceid").
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.UpdateUserServiceByServiceIdAndUserId("userid", "token", "tokenrefresh", "expiry", "serviceid")

		err = mock.ExpectationsWereMet()
		if err == nil {
			test.Errorf("Expectation fail")
		}
	})
}

func TestDeleteUserServiceByUserId(test *testing.T) {
	db, mock, repo := createMockDb(test)
	defer db.Close()

	sqlStatement := `DELETE FROM userservices WHERE userid = \(\$1\)`
	mock.ExpectExec(sqlStatement).
		WithArgs("userid").
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.DeleteUserServiceByUserId("userid")

	assert.NoError(test, err)

	err = mock.ExpectationsWereMet()
	if err != nil {
		test.Errorf("Expectation fail")
	}
}
