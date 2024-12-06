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

func (self *TriggerRepository) CreateTrigger(name string, description string, serviceId string, nbParam int) error {
	sqlStatement := `INSERT INTO actions (name, description, serviceid, nbparam) VALUES ($1, $2, $3, $4);`

	_, err := self.FindTriggerByName(name)
	if err == nil {
		return fmt.Errorf("trigger already exist")
	}

	_, err = self.db.Exec(sqlStatement, name, description, serviceId, nbParam)
	if err != nil {
		return err
	}
	return nil
}

func (self *TriggerRepository) FindTriggerByName(name string) (entities.Trigger, error) {
	sqlStatement := `SELECT * FROM actions where name = ($1)`
	var retrievedTrigger entities.Trigger

	row := self.db.QueryRow(sqlStatement, name)
	err := row.Scan(&retrievedTrigger.Id, &retrievedTrigger.ServiceId, &retrievedTrigger.Name, &retrievedTrigger.Description, &retrievedTrigger.NbParam)
	if err != nil {
		return retrievedTrigger, err
	}
	return retrievedTrigger, nil
}

func (self *TriggerRepository) FindTriggerByIdService(idService string) ([]entities.Trigger, error) {
	sqlStatement := `SELECT * FROM actions where serviceid = ($1)`
	var triggers []entities.Trigger

	rows, errQuery := self.db.Query(sqlStatement, idService)
	if errQuery != nil {
		return nil, errQuery
	}

	defer rows.Close()

	for rows.Next() {
		var retrievedTrigger entities.Trigger
		err := rows.Scan(&retrievedTrigger.Id, &retrievedTrigger.ServiceId, &retrievedTrigger.Name, &retrievedTrigger.Description, &retrievedTrigger.NbParam)
		if err != nil {
			return nil, err
		}
		triggers = append(triggers, retrievedTrigger)
	}
	return triggers, nil
}
