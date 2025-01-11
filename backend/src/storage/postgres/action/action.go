package action_repository

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"backend/src/entities"
)

type ActionRepository struct {
	db *sql.DB
}

func NewActionRepository(db *sql.DB) *ActionRepository {
	return &ActionRepository{db: db}
}

func unmarshalParameters(parametersBytes []byte, action entities.Action) (entities.Action, error) {
	if len(parametersBytes) > 0 {
		err := json.Unmarshal(parametersBytes, &action.Parameters)
		if err != nil {
			return action, err
		}
	}
	return action, nil
}

func (self *ActionRepository) CreateAction(name, description, serviceId string, nbParam int) error {
	sqlStatement := `INSERT INTO actions (name, description, serviceid, nbparam) VALUES ($1, $2, $3, $4)`

	_, err := self.FindActionByName(name)
	if err == nil {
		return fmt.Errorf("Action already exist")
	}

	_, err = self.db.Exec(sqlStatement, name, description, serviceId, nbParam)
	if err != nil {
		return err
	}
	return nil
}

func (self *ActionRepository) FindActionById(id string) (entities.Action, error) {
	sqlStatement := `SELECT * FROM actions WHERE id = ($1)`
	var action entities.Action
	var parametersBytes []byte

	row := self.db.QueryRow(sqlStatement, id)
	err := row.Scan(&action.Id, &action.ServiceId, &action.Name, &action.Description, &action.NbParam, &parametersBytes)
	if err != nil {
		return action, err
	}

	action, err = unmarshalParameters(parametersBytes, action)
	if err != nil {
		return action, err
	}
	return action, nil
}

func (self *ActionRepository) FindActionByName(name string) (entities.Action, error) {
	sqlStatement := `SELECT * FROM actions WHERE name = ($1)`
	var action entities.Action
	var parametersBytes []byte

	row := self.db.QueryRow(sqlStatement, name)
	err := row.Scan(&action.Id, &action.ServiceId, &action.Name, &action.Description, &action.NbParam, &parametersBytes)
	if err != nil {
		return action, err
	}

	action, err = unmarshalParameters(parametersBytes, action)
	if err != nil {
		return action, err
	}
	return action, nil
}

func (self *ActionRepository) FindActionsByServiceId(serviceId string) ([]entities.Action, error) {
	sqlStatement := `SELECT * FROM actions WHERE serviceid = ($1)`
	var actions []entities.Action

	rows, errQuery := self.db.Query(sqlStatement, serviceId)
	if errQuery != nil {
		return nil, errQuery
	}
	defer rows.Close()

	for rows.Next() {
		var action entities.Action
		var parametersBytes []byte

		err := rows.Scan(&action.Id, &action.ServiceId, &action.Name, &action.Description, &action.NbParam, &parametersBytes)
		if err != nil {
			return nil, err
		}

		action, err = unmarshalParameters(parametersBytes, action)
		if err != nil {
			return actions, err
		}
		actions = append(actions, action)
	}
	return actions, nil
}

func (self *ActionRepository) FindActionByNameAndServiceId(name, serviceId string) (entities.Action, error) {
	sqlStatement := `SELECT * FROM actions WHERE name = ($1) AND serviceid = ($2)`
	var action entities.Action
	var parametersBytes []byte

	row := self.db.QueryRow(sqlStatement, name, serviceId)
	err := row.Scan(&action.Id, &action.ServiceId, &action.Name, &action.Description, &action.NbParam, &parametersBytes)
	if err != nil {
		return action, err
	}

	action, err = unmarshalParameters(parametersBytes, action)
	if err != nil {
		return action, err
	}
	return action, nil
}
