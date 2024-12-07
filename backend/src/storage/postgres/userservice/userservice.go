package userservice_repository

import (
	"backend/src/entities"
	"database/sql"
	"fmt"
)

type UserServiceRepository struct {
	db *sql.DB
}

func NewUserServiceRepository(db *sql.DB) *UserServiceRepository {
	return &UserServiceRepository{db: db}
}

func (self *UserServiceRepository) CreateUserService(userId string, token string, tokenRefresh string, expiryDate string, serviceId string) error {
	sqlStatement := `INSERT INTO userservices (iduser, token, tokenrefresh, expiry, idservice) VALUES ($1, $2, $3, $4, $5);`

	_, err := self.FindUserServiceByServiceIdandUserId(userId, serviceId)
	if err == nil {
		return fmt.Errorf("service for user already exist")
	}

	_, err = self.db.Exec(sqlStatement, userId, token, tokenRefresh, expiryDate, serviceId)
	if err != nil {
		return err
	}
	return nil
}

func (self *UserServiceRepository) FindUserServiceByServiceIdandUserId(userId, serviceId string) (entities.UserService, error) {
	sqlStatement := `SELECT * FROM userservices where iduser = ($1) AND idservice = ($2)`
	var retrievedUserService entities.UserService

	row := self.db.QueryRow(sqlStatement, userId, serviceId)
	err := row.Scan(&retrievedUserService.Id, &retrievedUserService.UserId, &retrievedUserService.AccessToken,
		&retrievedUserService.RefreshToken, &retrievedUserService.ExpiryDate, &retrievedUserService.ServiceId)
	if err != nil {
		return retrievedUserService, err
	}
	return retrievedUserService, nil
}

func (self *UserServiceRepository) UpdateUserServiceByServiceIdAndUserId(uuserviceId string, uuuserId string, accessToken string, refreshToken string, expiryToken string) error {
	sqlStatement := `UPDATE userservices SET token = ($1), tokenrefresh = ($2), expiry = ($3) WHERE iduser = ($4) AND idservice = ($5);`

	_, err := self.FindUserServiceByServiceIdandUserId(uuuserId, uuserviceId)
	if err != nil {
		return fmt.Errorf("service doesn't exist for user")
	}

	_, err = self.db.Exec(sqlStatement, accessToken, refreshToken, expiryToken, uuuserId, uuserviceId)
	if err != nil {
		return err
	}
	return nil
}
