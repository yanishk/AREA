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

func (self *ReactionRepository) UnmarshalParameters(parametersBytes []byte, reaction entities.Reaction) (entities.Reaction, error) {
	if len(parametersBytes) > 0 {
		err := json.Unmarshal(parametersBytes, &reaction.Parameters)
		if err != nil {
			return reaction, err
		}
	}
	return reaction, nil
}

func (self *ReactionRepository) CreateReaction(name string, description string, serviceId string, nbParam int) error {
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
	var retrievedReaction entities.Reaction
	var parametersBytes []byte

	row := self.db.QueryRow(sqlStatement, id)
	err := row.Scan(&retrievedReaction.Id, &retrievedReaction.ServiceId, &retrievedReaction.Name, &retrievedReaction.Description, &retrievedReaction.NbParam, &parametersBytes)
	if err != nil {
		return retrievedReaction, err
	}

	retrievedReaction, err = self.UnmarshalParameters(parametersBytes, retrievedReaction)
	if err != nil {
		return retrievedReaction, err
	}
	return retrievedReaction, nil
}

func (self *ReactionRepository) FindReactionByName(name string) (entities.Reaction, error) {
	sqlStatement := `SELECT * FROM reactions WHERE name = ($1)`
	var retrievedReaction entities.Reaction
	var parametersBytes []byte

	row := self.db.QueryRow(sqlStatement, name)
	err := row.Scan(&retrievedReaction.Id, &retrievedReaction.ServiceId, &retrievedReaction.Name, &retrievedReaction.Description, &retrievedReaction.NbParam, &parametersBytes)
	if err != nil {
		return retrievedReaction, err
	}

	retrievedReaction, err = self.UnmarshalParameters(parametersBytes, retrievedReaction)
	if err != nil {
		return retrievedReaction, err
	}
	return retrievedReaction, nil
}

func (self *ReactionRepository) FindReactionByServiceId(serviceId string) ([]entities.Reaction, error) {
	sqlStatement := `SELECT * FROM reactions WHERE serviceid = ($1)`
	var reactions []entities.Reaction

	rows, errQuery := self.db.Query(sqlStatement, serviceId)
	if errQuery != nil {
		return nil, errQuery
	}
	defer rows.Close()

	for rows.Next() {
		var retrievedReaction entities.Reaction
		var parametersBytes []byte

		err := rows.Scan(&retrievedReaction.Id, &retrievedReaction.ServiceId, &retrievedReaction.Name, &retrievedReaction.Description, &retrievedReaction.NbParam, &parametersBytes)
		if err != nil {
			return nil, err
		}

		retrievedReaction, err = self.UnmarshalParameters(parametersBytes, retrievedReaction)
		if err != nil {
			return reactions, err
		}
		reactions = append(reactions, retrievedReaction)
	}
	return reactions, nil
}
