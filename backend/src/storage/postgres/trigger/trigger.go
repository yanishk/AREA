package trigger_repository

import (
	"database/sql"
	"fmt"

	"backend/src/entities"
)

type TriggerRepository struct {
	db *sql.DB
}

func NewTriggerRepository(db *sql.DB) *TriggerRepository {
	return &TriggerRepository{db: db}
}

func (self *TriggerRepository) CreateTrigger(name string, description string, serviceId string) error {
	sqlStatement := `INSERT INTO triggers (name, description, serviceid) VALUES ($1, $2, $3);`

	_, err := self.FindTriggerByName(name)
	if err == nil {
		return fmt.Errorf("Trigger already exist")
	}

	_, err = self.db.Exec(sqlStatement, name, description, serviceId)
	if err != nil {
		return err
	}
	return nil
}

func (self *TriggerRepository) FindTriggerByName(name string) (entities.Trigger, error) {
	sqlStatement := `SELECT * FROM triggers where name = ($1)`
	var retrievedTrigger entities.Trigger

	row := self.db.QueryRow(sqlStatement, name)
	err := row.Scan(&retrievedTrigger.Id, &retrievedTrigger.ServiceId, &retrievedTrigger.Name, &retrievedTrigger.Description)
	if err != nil {
		return retrievedTrigger, err
	}
	return retrievedTrigger, nil
}
