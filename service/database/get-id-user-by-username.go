package database

import (
	"database/sql"
	"errors"
)

// If the user is not found, return id 0
func (db *appdbimpl) GetIdUserByUsername(username string) (string, error) {
	var id string

	err := db.c.QueryRow("SELECT id FROM users WHERE name=?", username).Scan(&id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil
		}
		return "", err
	}
	return id, nil
}
