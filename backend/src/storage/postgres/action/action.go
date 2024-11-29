package action_repository

import (
	"fmt"
	"database/sql"

	"backend/src/entities"
)

type ActionRepository struct {
	db *sql.DB
}

func NewActionRepository(db *sql.DB) *ActionRepository {
	return &ActionRepository {db: db}
}
