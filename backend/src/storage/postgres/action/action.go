package action_repository

import (
	"database/sql"
	"fmt"

	"backend/src/entities"
)

type actionRepository struct {
	db *sql.DB
}

func NewActionRepository(db *sql.DB) *actionRepository {
	return &actionRepository{db: db}
}

func (self *actionRepository) CreateAction(name string, description string, serviceId string, nbParam int) error {
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

func (self *actionRepository) FindActionByName(name string) (entities.Action, error) {
	sqlStatement := `SELECT * FROM actions where name = ($1)`
	var retrievedaction entities.Action

	row := self.db.QueryRow(sqlStatement, name)
	err := row.Scan(&retrievedaction.Id, &retrievedaction.ServiceId, &retrievedaction.Name, &retrievedaction.Description, &retrievedaction.NbParam)
	if err != nil {
		return retrievedaction, err
	}
	return retrievedaction, nil
}

func (self *actionRepository) FindActionByServiceId(serviceId string) ([]entities.Action, error) {
	sqlStatement := `SELECT * FROM actions where serviceid = ($1)`
	var actions []entities.Action

	rows, errQuery := self.db.Query(sqlStatement, serviceId)
	if errQuery != nil {
		return nil, errQuery
	}

	defer rows.Close()

	for rows.Next() {
		var retrievedaction entities.Action
		err := rows.Scan(&retrievedaction.Id, &retrievedaction.ServiceId, &retrievedaction.Name, &retrievedaction.Description, &retrievedaction.NbParam)
		if err != nil {
			return nil, err
		}
		actions = append(actions, retrievedaction)
	}
	return actions, nil
}
