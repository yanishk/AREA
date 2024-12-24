package user_repository

import (
	"database/sql"
	"fmt"

	"backend/src/entities"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (self *UserRepository) CreateUser(email string, password string, connectionType string) error {
	sqlStatement := `INSERT INTO users (email, password, connectiontype) VALUES ($1, $2, $3);`

	_, err := self.FindUserByEmail(email, connectionType)
	if err == nil {
		return fmt.Errorf("User already exist")
	}

	_, err = self.db.Exec(sqlStatement, email, password, connectionType)
	if err != nil {
		return err
	}
	return nil
}

func (self *UserRepository) FindUserByEmail(email string, connectionType string) (entities.User, error) {
	sqlStatement := `SELECT * FROM users WHERE email = ($1) AND connectiontype = ($2)`
	var retrievedUser entities.User

	row := self.db.QueryRow(sqlStatement, email, connectionType)
	err := row.Scan(&retrievedUser.Email, &retrievedUser.Password, &retrievedUser.Id, &retrievedUser.CreatedAt, &retrievedUser.Timezone, &retrievedUser.ConnectionType)
	if err != nil {
		return retrievedUser, err
	}
	return retrievedUser, nil
}

func (self *UserRepository) FindUserById(userId string) (entities.User, error) {
	sqlStatement := `SELECT * FROM users WHERE id = ($1)`
	var retrievedUser entities.User

	row := self.db.QueryRow(sqlStatement, userId)
	err := row.Scan(&retrievedUser.Email, &retrievedUser.Password, &retrievedUser.Id, &retrievedUser.CreatedAt, &retrievedUser.Timezone, &retrievedUser.ConnectionType)
	if err != nil {
		return retrievedUser, err
	}
	return retrievedUser, nil
}

func (self *UserRepository) UpdateUser(email string, password string, connectionType string) error {
	sqlStatement := `UPDATE users SET password = ($1) WHERE email = ($2);`

	_, err := self.FindUserByEmail(email, connectionType)
	if err != nil {
		return fmt.Errorf("User doesn't exist")
	}

	_, err = self.db.Exec(sqlStatement, password, email)
	if err != nil {
		return err
	}
	return nil
}

func (self *UserRepository) DeleteUser(email string, connectionType string) error {
	sqlStatement := `DELETE FROM users WHERE email = ($1)`

	_, err := self.FindUserByEmail(email, connectionType)
	if err != nil {
		return fmt.Errorf("User doesn't exist")
	}

	_, err = self.db.Exec(sqlStatement, email)
	if err != nil {
		return err
	}
	return nil
}
