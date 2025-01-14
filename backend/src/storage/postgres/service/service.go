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

func (self *ServiceRepository) CreateService(name, color, logo string) error {
	sqlStatement := `INSERT INTO services (name, color, logo) VALUES ($1, $2, $3)`

	_, err := self.FindServiceByName(name)
	if err == nil {
		return fmt.Errorf("Service already exist")
	}

	_, err = self.db.Exec(sqlStatement, name, color, logo)
	if err != nil {
		return err
	}
	return nil
}

func (self *ServiceRepository) findServiceWithParam(sqlStatement, param string) (entities.Service, error) {
	var service entities.Service

	row := self.db.QueryRow(sqlStatement, param)
	err := row.Scan(&service.Id, &service.Name, &service.Color, &service.Logo, &service.HasActions, &service.HasReactions, &service.IsAuthNeeded, &service.Description)
	if err != nil {
		return service, err
	}
	return service, nil
}

func (self *ServiceRepository) findServices(sqlStatement string) ([]entities.Service, error) {
	var services []entities.Service

	rows, err := self.db.Query(sqlStatement)
	if err != nil {
		return services, err
	}
	defer rows.Close()

	for rows.Next() {
		var service entities.Service
		err := rows.Scan(&service.Id, &service.Name, &service.Color, &service.Logo, &service.HasActions, &service.HasReactions, &service.IsAuthNeeded, &service.Description)
		if err != nil {
			return services, err
		}
		services = append(services, service)
	}
	return services, nil
}

func (self *ServiceRepository) FindServiceById(id string) (entities.Service, error) {
	sqlStatement := `SELECT * FROM services WHERE id = ($1)`
	return self.findServiceWithParam(sqlStatement, id)
}

func (self *ServiceRepository) FindServiceByName(name string) (entities.Service, error) {
	sqlStatement := `SELECT * FROM services WHERE name = ($1)`
	return self.findServiceWithParam(sqlStatement, name)
}

func (self *ServiceRepository) FindAllServices() ([]entities.Service, error) {
	sqlStatement := `SELECT * FROM services`
	return self.findServices(sqlStatement)
}

func (self *ServiceRepository) FindActionsServices() ([]entities.Service, error) {
	sqlStatement := `SELECT * FROM services WHERE hasactions = true`
	return self.findServices(sqlStatement)
}

func (self *ServiceRepository) FindReactionsServices() ([]entities.Service, error) {
	sqlStatement := `SELECT * FROM services WHERE hasreactions = true`
	return self.findServices(sqlStatement)
}
