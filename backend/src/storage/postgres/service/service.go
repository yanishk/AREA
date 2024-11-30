package service_repository

import (
	"fmt"
	"database/sql"

	"backend/src/entities"
)

type ServiceRepository struct {
	db *sql.DB
}

func NewServiceRepository(db *sql.DB) *ServiceRepository {
	return &ServiceRepository {db: db}
}

func (self *ServiceRepository) CreateService(name string) error {
	sqlStatement := `INSERT INTO services (name) VALUES ($1);`

	_, err := self.FindServiceByName(name)
	if err == nil {
		return fmt.Errorf("Service already exist")
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
