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

func (self *UserRepository) CreateUser(email, password, connectionType string) error {
	sqlStatement := `INSERT INTO users (email, password, connectiontype) VALUES ($1, $2, $3)`

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

func (self *UserRepository) FindUserByEmail(email, connectionType string) (entities.User, error) {
	sqlStatement := `SELECT * FROM users WHERE email = ($1) AND connectiontype = ($2)`
	var user entities.User

	row := self.db.QueryRow(sqlStatement, email, connectionType)
	err := row.Scan(&user.Email, &user.Password, &user.Id, &user.CreatedAt, &user.Timezone, &user.ConnectionType)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (self *UserRepository) FindUserById(userId string) (entities.User, error) {
	sqlStatement := `SELECT * FROM users WHERE id = ($1)`
	var user entities.User

	row := self.db.QueryRow(sqlStatement, userId)
	err := row.Scan(&user.Email, &user.Password, &user.Id, &user.CreatedAt, &user.Timezone, &user.ConnectionType)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (self *UserRepository) UpdateUser(email, password, connectionType string) error {
	sqlStatement := `UPDATE users SET password = ($1) WHERE email = ($2) AND connectiontype = ($3)`

	_, err := self.FindUserByEmail(email, connectionType)
	if err != nil {
		return fmt.Errorf("User doesn't exist")
	}

	_, err = self.db.Exec(sqlStatement, password, email, connectionType)
	if err != nil {
		return err
	}
	return nil
}

func (self *UserRepository) DeleteUser(email, connectionType string) error {
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
