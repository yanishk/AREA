package reaction_repository

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"backend/src/entities"
)

type ReactionRepository struct {
	db *sql.DB
}

func NewReactionRepository(db *sql.DB) *ReactionRepository {
	return &ReactionRepository{db: db}
}

func unmarshalParameters(parametersBytes []byte, reaction entities.Reaction) (entities.Reaction, error) {
	if len(parametersBytes) > 0 {
		err := json.Unmarshal(parametersBytes, &reaction.Parameters)
		if err != nil {
			return reaction, err
		}
	}
	return reaction, nil
}

func (self *ReactionRepository) CreateReaction(name, description, serviceId string, nbParam int) error {
	sqlStatement := `INSERT INTO reactions (name, description, serviceid, nbparam) VALUES ($1, $2, $3, $4)`

	_, err := self.FindReactionByName(name)
	if err == nil {
		return fmt.Errorf("Reaction already exist")
	}

	_, err = self.db.Exec(sqlStatement, name, description, serviceId, nbParam)
	if err != nil {
		return err
	}
	return nil
}

func (self *ReactionRepository) FindReactionById(id string) (entities.Reaction, error) {
	sqlStatement := `SELECT * FROM reactions WHERE id = ($1)`
	var reaction entities.Reaction
	var parametersBytes []byte

	row := self.db.QueryRow(sqlStatement, id)
	err := row.Scan(&reaction.Id, &reaction.ServiceId, &reaction.Name, &reaction.Description, &reaction.NbParam, &parametersBytes)
	if err != nil {
		return reaction, err
	}

	reaction, err = unmarshalParameters(parametersBytes, reaction)
	if err != nil {
		return reaction, err
	}
	return reaction, nil
}

func (self *ReactionRepository) FindReactionByName(name string) (entities.Reaction, error) {
	sqlStatement := `SELECT * FROM reactions WHERE name = ($1)`
	var reaction entities.Reaction
	var parametersBytes []byte

	row := self.db.QueryRow(sqlStatement, name)
	err := row.Scan(&reaction.Id, &reaction.ServiceId, &reaction.Name, &reaction.Description, &reaction.NbParam, &parametersBytes)
	if err != nil {
		return reaction, err
	}

	reaction, err = unmarshalParameters(parametersBytes, reaction)
	if err != nil {
		return reaction, err
	}
	return reaction, nil
}

func (self *ReactionRepository) FindReactionsByServiceId(serviceId string) ([]entities.Reaction, error) {
	sqlStatement := `SELECT * FROM reactions WHERE serviceid = ($1)`
	var reactions []entities.Reaction

	rows, errQuery := self.db.Query(sqlStatement, serviceId)
	if errQuery != nil {
		return nil, errQuery
	}
	defer rows.Close()

	for rows.Next() {
		var reaction entities.Reaction
		var parametersBytes []byte

		err := rows.Scan(&reaction.Id, &reaction.ServiceId, &reaction.Name, &reaction.Description, &reaction.NbParam, &parametersBytes)
		if err != nil {
			return nil, err
		}

		reaction, err = unmarshalParameters(parametersBytes, reaction)
		if err != nil {
			return reactions, err
		}
		reactions = append(reactions, reaction)
	}
	return reactions, nil
}
