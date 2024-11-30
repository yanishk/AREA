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

func (self *ServiceRepository) CreateService(name string) error {
	sqlStatement := `INSERT INTO services (name) VALUES ($1);`

	_, err := self.FindServiceByName(name)
	if err == nil {
		return fmt.Errorf("service already exist")
	}

	_, err = self.db.Exec(sqlStatement, name)
	if err != nil {
		return err
	}
	return nil
}

func (self *ServiceRepository) FindServiceByName(name string) (entities.Service, error) {
	sqlStatement := `SELECT * FROM services where name = ($1)`
	var retrievedService entities.Service

	row := self.db.QueryRow(sqlStatement, name)
	err := row.Scan(&retrievedService.Id, &retrievedService.Name)
	if err != nil {
		return retrievedService, err
	}
	return retrievedService, nil
}

func (self *ServiceRepository) FindUserServiceByName(name string, userId string) (entities.Service, error) {
	sqlStatement := `SELECT * FROM services where name = ($1) AND iduser = ($2)`
	var retrievedService entities.Service

	row := self.db.QueryRow(sqlStatement, name, userId)
	err := row.Scan(&retrievedService.Id, &retrievedService.Name, &retrievedService.AccessToken, &retrievedService.RefreshToken, &retrievedService.IdUser, &retrievedService.ExpiryToken)
	if err != nil {
		return retrievedService, err
	}
	return retrievedService, nil
}

func (self *ServiceRepository) UpdateServiceByUserAndName(name string, accessToken string, refreshToken string, uuidUser string, expiryToken string) error {
	sqlStatement := `UPDATE services SET accessToken = ($1), refreshToken = ($2), expirytoken = ($3) WHERE name = ($4);`

	_, err := self.FindUserServiceByName(name, uuidUser)
	if err != nil {
		return fmt.Errorf("service doesn't exist")
	}

	_, err = self.db.Exec(sqlStatement, accessToken, refreshToken, expiryToken, name)
	if err != nil {
		return err
	}
	return nil
}
