package userservice_repository

import (
	"database/sql"
	"fmt"

	"backend/src/entities"
)

type UserServiceRepository struct {
	db *sql.DB
}

func NewUserServiceRepository(db *sql.DB) *UserServiceRepository {
	return &UserServiceRepository{db: db}
}

func (self *UserServiceRepository) CreateUserService(userId, token, tokenRefresh, expiryDate, serviceId string) error {
	sqlStatement := `INSERT INTO userservices (userid, token, tokenrefresh, expiry, serviceid) VALUES ($1, $2, $3, $4, $5)`

	_, err := self.FindUserServiceByServiceIdandUserId(userId, serviceId)
	if err == nil {
		return fmt.Errorf("Service for user already exist")
	}

	_, err = self.db.Exec(sqlStatement, userId, token, tokenRefresh, expiryDate, serviceId)
	if err != nil {
		return err
	}
	return nil
}

func (self *UserServiceRepository) FindUserServiceByServiceIdandUserId(userId, serviceId string) (entities.UserService, error) {
	sqlStatement := `SELECT * FROM userservices WHERE userid = ($1) AND serviceid = ($2)`
	var userService entities.UserService

	row := self.db.QueryRow(sqlStatement, userId, serviceId)
	err := row.Scan(&userService.Id, &userService.UserId, &userService.AccessToken,
		&userService.RefreshToken, &userService.ExpiryDate, &userService.ServiceId)
	if err != nil {
		return userService, err
	}
	return userService, nil
}

func (self *UserServiceRepository) UpdateUserServiceByServiceIdAndUserId(userId, accessToken, refreshToken, expiryDate, serviceId string) error {
	sqlStatement := `UPDATE userservices SET token = ($1), tokenrefresh = ($2), expiry = ($3) WHERE userid = ($4) AND serviceid = ($5)`

	_, err := self.FindUserServiceByServiceIdandUserId(userId, serviceId)
	if err != nil {
		return fmt.Errorf("Service doesn't exist for user")
	}

	_, err = self.db.Exec(sqlStatement, accessToken, refreshToken, expiryDate, userId, serviceId)
	if err != nil {
		return err
	}
	return nil
}

func (self *UserServiceRepository) DeleteUserServiceByUserId(userId string) error {
	sqlStatement := `DELETE FROM userservices WHERE userid = ($1)`

	_, err := self.db.Exec(sqlStatement, userId)
	if err != nil {
		return err
	}
	return nil
}
