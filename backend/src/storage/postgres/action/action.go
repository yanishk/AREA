package action_repository

import (
	"database/sql"
	"fmt"

	"backend/src/entities"
)

type ReactionRepository struct {
	db *sql.DB
}

func NewReactionRepository(db *sql.DB) *ReactionRepository {
	return &ReactionRepository{db: db}
}

func (self *ReactionRepository) CreateReaction(name string, description string, serviceId string, nbParam int) error {
	sqlStatement := `INSERT INTO reactions (name, description, serviceid, nbparam) VALUES ($1, $2, $3, $4);`

	_, err := self.FindReactionByName(name)
	if err == nil {
		return fmt.Errorf("Action already exist")
	}

	_, err = self.db.Exec(sqlStatement, name, description, serviceId, nbParam)
	if err != nil {
		return err
	}
	return nil
}

func (self *ReactionRepository) FindReactionByName(name string) (entities.Action, error) {
	sqlStatement := `SELECT * FROM reactions where name = ($1)`
	var retrievedAction entities.Action

	row := self.db.QueryRow(sqlStatement, name)
	err := row.Scan(&retrievedAction.Id, &retrievedAction.ServiceId, &retrievedAction.Name, &retrievedAction.Description, &retrievedAction.NbParam)
	if err != nil {
		return retrievedAction, err
	}
	return retrievedAction, nil
}

func (self *ReactionRepository) FindReactionByID(id string) (entities.Action, error) {
	sqlStatement := `SELECT * FROM reactions where id = ($1)`
	var retrievedAction entities.Action

	row := self.db.QueryRow(sqlStatement, id)
	err := row.Scan(&retrievedAction.Id, &retrievedAction.ServiceId, &retrievedAction.Name, &retrievedAction.Description, &retrievedAction.NbParam)
	if err != nil {
		return retrievedAction, err
	}
	return retrievedAction, nil
}
