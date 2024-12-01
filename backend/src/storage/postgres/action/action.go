package action_repository

import (
	"database/sql"
	"fmt"

	"backend/src/entities"
)

type ActionRepository struct {
	db *sql.DB
}

func NewActionRepository(db *sql.DB) *ActionRepository {
	return &ActionRepository{db: db}
}

func (self *ActionRepository) CreateAction(name string, description string, serviceId string) error {
	sqlStatement := `INSERT INTO actions (name, description, serviceid) VALUES ($1, $2, $3);`

	_, err := self.FindActionByName(name)
	if err == nil {
		return fmt.Errorf("Action already exist")
	}

	_, err = self.db.Exec(sqlStatement, name, description, serviceId)
	if err != nil {
		return err
	}
	return nil
}

func (self *ActionRepository) FindActionByName(name string) (entities.Action, error) {
	sqlStatement := `SELECT * FROM actions where name = ($1)`
	var retrievedAction entities.Action

	row := self.db.QueryRow(sqlStatement, name)
	err := row.Scan(&retrievedAction.Id, &retrievedAction.ServiceId, &retrievedAction.Name, &retrievedAction.Description)
	if err != nil {
		return retrievedAction, err
	}
	return retrievedAction, nil
}
