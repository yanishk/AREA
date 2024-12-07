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

func (self *UserServiceRepository) CreateUserService(idUser string, token string, tokenRefresh string, expiryDate string, idService string) error {
	sqlStatement := `INSERT INTO userservices (iduser, token, tokenrefresh, expiry, idservice) VALUES ($1, $2, $3, $4, $5);`

	_, err := self.FindUserServiceByIdServiceandIdUser(idUser, idService)
	if err == nil {
		return fmt.Errorf("service for user already exist")
	}

	_, err = self.db.Exec(sqlStatement, idUser, token, tokenRefresh, expiryDate, idService)
	if err != nil {
		return err
	}
	return nil
}

func (self *UserServiceRepository) FindUserServiceByIdServiceandIdUser(idUser, idService string) (entities.UserService, error) {
	sqlStatement := `SELECT * FROM userservices where iduser = ($1) AND idservice = ($2)`
	var retrievedUserService entities.UserService

	row := self.db.QueryRow(sqlStatement, idUser, idService)
	err := row.Scan(&retrievedUserService.Id, &retrievedUserService.IdUser, &retrievedUserService.AccessToken,
		&retrievedUserService.RefreshToken, &retrievedUserService.ExpiryDate, &retrievedUserService.IdService)
	if err != nil {
		return retrievedUserService, err
	}
	return retrievedUserService, nil
}

func (self *UserServiceRepository) UpdateUserServiceByServiceIdAndUserId(uuidService string, uuidUser string, accessToken string, refreshToken string, expiryToken string) error {
	sqlStatement := `UPDATE userservices SET token = ($1), tokenrefresh = ($2), expiry = ($3) WHERE iduser = ($4) AND idservice = ($5);`

	_, err := self.FindUserServiceByIdServiceandIdUser(idUser, idService)
	if err != nil {
		return fmt.Errorf("service doesn't exist for user")
	}

	_, err = self.db.Exec(sqlStatement, token, tokenRefresh, expiryDate, idUser, idService)
	if err != nil {
		return err
	}
	return nil
}
