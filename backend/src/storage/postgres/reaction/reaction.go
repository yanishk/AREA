package reaction_repository

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

func (self *ReactionRepository) FindReactionByName(name string) (entities.Reaction, error) {
	sqlStatement := `SELECT * FROM reactions where name = ($1)`
	var retrievedReaction entities.Reaction

	row := self.db.QueryRow(sqlStatement, name)
	err := row.Scan(&retrievedReaction.Id, &retrievedReaction.ServiceId, &retrievedReaction.Name, &retrievedReaction.Description, &retrievedReaction.NbParam)
	if err != nil {
		return retrievedReaction, err
	}
	return retrievedReaction, nil
}

func (self *ReactionRepository) FindReactionByServiceId(serviceId string) ([]entities.Reaction, error) {
	sqlStatement := `SELECT * FROM reactions where serviceid = ($1)`
	var reactions []entities.Reaction

	rows, errQuery := self.db.Query(sqlStatement, serviceId)
	if errQuery != nil {
		return nil, errQuery
	}

	defer rows.Close()

	for rows.Next() {
		var retrievedReaction entities.Reaction
		err := rows.Scan(&retrievedReaction.Id, &retrievedReaction.ServiceId, &retrievedReaction.Name, &retrievedReaction.Description, &retrievedReaction.NbParam)
		if err != nil {
			return nil, err
		}
		reactions = append(reactions, retrievedReaction)
	}
	return reactions, nil
}

func (self *ReactionRepository) FindReactionByID(id string) (entities.Reaction, error) {
	sqlStatement := `SELECT * FROM reactions where id = ($1)`
	var retrievedReaction entities.Reaction

	row := self.db.QueryRow(sqlStatement, id)
	err := row.Scan(&retrievedReaction.Id, &retrievedReaction.ServiceId, &retrievedReaction.Name, &retrievedReaction.Description, &retrievedReaction.NbParam)
	if err != nil {
		return retrievedReaction, err
	}
	return retrievedReaction, nil
}
