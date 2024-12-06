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

func (self *ServiceRepository) FindServiceByName(name string) (entities.Service, error) {
	sqlStatement := `SELECT * FROM services where name = ($1)`
	var retrievedService entities.Service

	row := self.db.QueryRow(sqlStatement, name)
	err := row.Scan(&retrievedService.Id, &retrievedService.Name, &retrievedService.Color, &retrievedService.Logo, &retrievedService.HasActions, &retrievedService.HasReactions)
	if err != nil {
		return retrievedService, err
	}
	return retrievedService, nil
}

func (self *ServiceRepository) FindActionsServices() ([]entities.Service, error) {
	sqlStatement := `SELECT name FROM services where hasactions = true`
	var retrievedServices []entities.Service

	rows, err := self.db.Query(sqlStatement)
	if (err != nil) {
		return retrievedServices, err
	}
	defer rows.Close()

	for rows.Next() {
		var retrievedService entities.Service
		err := rows.Scan(&retrievedService.Name)
		if err != nil {
		    return retrievedServices, err
		}
		retrievedServices = append(retrievedServices, retrievedService)
	}
	return retrievedServices, nil
}

func (self *ServiceRepository) FindReactionsServices() ([]entities.Service, error) {
	sqlStatement := `SELECT name FROM services where hasreactions = true`
	var retrievedServices []entities.Service

	rows, err := self.db.Query(sqlStatement)
	if (err != nil) {
		return retrievedServices, err
	}
	defer rows.Close()

	for rows.Next() {
		var retrievedService entities.Service
		err := rows.Scan(&retrievedService.Name)
		if err != nil {
		    return retrievedServices, err
		}
		retrievedServices = append(retrievedServices, retrievedService)
	}
	return retrievedServices, nil
}

