package service_repository

import (
	"database/sql"
	"fmt"

	"backend/src/entities"
)

type ServiceRepository struct {
	db *sql.DB
}

func NewServiceRepository(db *sql.DB) *ServiceRepository {
	return &ServiceRepository{db: db}
}

func (self *ServiceRepository) CreateService(name string, color string, logo string) error {
	sqlStatement := `INSERT INTO services (name, color, logo) VALUES ($1, $2, $3);`

	_, err := self.FindServiceByName(name)
	if err == nil {
		return fmt.Errorf("service already exist")
	}

	_, err = self.db.Exec(sqlStatement, name, color, logo)
	if err != nil {
		return err
	}
	return nil
}

func (self *ServiceRepository) CreateUserService(idUser string, token string, refreshToken string, expiryDate string, idService string) error {
	sqlStatement := `INSERT INTO userservices (iduser, token, tokenrefresh, expiry, idservice) VALUES ($1, $2, $3, $4, $5);`

	_, err := self.FindUserServiceByIdServiceandIdUser(idService, idUser)
	if err == nil {
		return fmt.Errorf("service for user already exist")
	}

	_, err = self.db.Exec(sqlStatement, idUser, token, refreshToken, expiryDate, idService)
	if err != nil {
		return err
	}
	return nil
}

func (self *ServiceRepository) FindServiceByName(name string) (entities.Service, error) {
	sqlStatement := `SELECT * FROM services where name = ($1)`
	var retrievedService entities.Service

	row := self.db.QueryRow(sqlStatement, name)
	err := row.Scan(&retrievedService.Id, &retrievedService.Name, &retrievedService.Color, &retrievedService.Logo)
	if err != nil {
		return retrievedService, err
	}
	return retrievedService, nil
}

func (self *ServiceRepository) FindUserServiceByIdServiceandIdUser(idservice string, userId string) (entities.UserService, error) {
	sqlStatement := `SELECT * FROM userservices where iduser = ($1) AND idservice = ($2)`
	var retrievedUserService entities.UserService

	row := self.db.QueryRow(sqlStatement, userId, idservice)
	err := row.Scan(&retrievedUserService.Id, &retrievedUserService.IdUser, &retrievedUserService.AccessToken,
		&retrievedUserService.RefreshToken, &retrievedUserService.ExpiryToken, &retrievedUserService.IdService)
	if err != nil {
		return retrievedUserService, err
	}
	return retrievedUserService, nil
}

func (self *ServiceRepository) UpdateUserServiceByServiceIdAndUserId(uuidService string, uuidUser string, accessToken string, refreshToken string, expiryToken string) error {
	sqlStatement := `UPDATE userservices SET accessToken = ($1), refreshToken = ($2), expirytoken = ($3) WHERE iduser = ($4) AND idservice = ($5);`

	_, err := self.FindUserServiceByIdServiceandIdUser(uuidService, uuidUser)
	if err != nil {
		return fmt.Errorf("service doesn't exist for user")
	}

	_, err = self.db.Exec(sqlStatement, accessToken, refreshToken, expiryToken, uuidUser, uuidService)
	if err != nil {
		return err
	}
	return nil
}
