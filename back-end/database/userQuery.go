package database

import (
	"database/sql"
	"fmt"
)

type User struct {
	Id        string
	Email     string
	Password  string
	CreatedAt string
}

func CreateUser(db *sql.DB, email string, password string) error {
	sqlStatement := `INSERT INTO users (email, password) VALUES ($1, $2);`

	_, err := FindUserByEmail(db, email)
	if err == nil {
		return fmt.Errorf("User already exist")
	}

	_, err = db.Exec(sqlStatement, email, password)
	if err != nil {
		return err
	}
	return nil
}

func FindUserByEmail(db *sql.DB, email string) (User, error) {
	sqlStatement := `SELECT * FROM users where email = ($1)`
	var retrievedUser User

	row := db.QueryRow(sqlStatement, email)
	err := row.Scan(&retrievedUser.Email, &retrievedUser.Password, &retrievedUser.Id, &retrievedUser.CreatedAt)
	if err != nil {
		return retrievedUser, err
	}
	return retrievedUser, nil
}

func UpdateUser(db *sql.DB, email string, password string) error {
	sqlStatement := `UPDATE users SET password = ($1) WHERE email = ($2);`

	_, err := FindUserByEmail(db, email)
	if err != nil {
		return fmt.Errorf("User doesn't exist")
	}

	_, err = db.Exec(sqlStatement, password, email)
	if err != nil {
		return err
	}
	return nil
}

func DeleteUser(db *sql.DB, email string) error {
	sqlStatement := `DELETE FROM users WHERE email = ($1)`

	_, err := FindUserByEmail(db, email)
	if err != nil {
		return fmt.Errorf("User doesn't exist")
	}

	_, err = db.Exec(sqlStatement, email)
	if err != nil {
		return err
	}
	return nil
}
