package database

import (
	"database/sql"
	"errors"
)

// If the user is not found, return id 0
func (db *appdbimpl) GetIdUserByUsername(username string) (int, error) {
	var id int

	err := db.c.QueryRow("SELECT id FROM users WHERE name=?", username).Scan(&id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}
		return 0, err
	}
	return id, nil
}
