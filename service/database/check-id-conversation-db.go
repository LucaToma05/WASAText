package database

import (
	"database/sql"
	"errors"
)

// If the user is not found, return id ""
func (db *appdbimpl) CheckIdConversation(id string) (bool, error) {
	var flag bool
	err := db.c.QueryRow("SELECT 1 FROM chats WHERE id=?", id).Scan(&flag)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
