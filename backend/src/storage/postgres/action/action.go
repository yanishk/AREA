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

func (self *ActionRepository) CreateAction(name string, description string, serviceId string, nbParam int) error {
	sqlStatement := `INSERT INTO actions (name, description, serviceid, nbparam) VALUES ($1, $2, $3, $4);`

	_, err := self.FindActionByName(name)
	if err == nil {
		return fmt.Errorf("action already exist")
	}

	_, err = self.db.Exec(sqlStatement, name, description, serviceId, nbParam)
	if err != nil {
		return err
	}
	return nil
}

func (self *ActionRepository) FindActionByName(name string) (entities.Action, error) {
	sqlStatement := `SELECT * FROM actions where name = ($1)`
	var retrievedAction entities.Action

	row := self.db.QueryRow(sqlStatement, name)
	err := row.Scan(&retrievedAction.Id, &retrievedAction.ServiceId, &retrievedAction.Name, &retrievedAction.Description, &retrievedAction.NbParam)
	if err != nil {
		return retrievedAction, err
	}
	return retrievedAction, nil
}

func (self *ActionRepository) FindActionByServiceId(serviceId string) ([]entities.Action, error) {
	sqlStatement := `SELECT * FROM actions where serviceid = ($1)`
	var actions []entities.Action

	rows, errQuery := self.db.Query(sqlStatement, serviceId)
	if errQuery != nil {
		return nil, errQuery
	}

	defer rows.Close()

	for rows.Next() {
		var retrievedAction entities.Action
		err := rows.Scan(&retrievedAction.Id, &retrievedAction.ServiceId, &retrievedAction.Name, &retrievedAction.Description, &retrievedAction.NbParam)
		if err != nil {
			return nil, err
		}
		actions = append(actions, retrievedAction)
	}
	return actions, nil
}
