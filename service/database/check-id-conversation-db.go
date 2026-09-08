package database

import (
	"database/sql"
	"errors"
)

// Verifica se l'id della chat esiste
func (db *appdbimpl) CheckIdConversation(id string) (bool, error) {
	var flag bool
	err := db.c.QueryRow("SELECT 1 FROM chats WHERE id=?", id).Scan(&flag)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil // id inesistente
		}
		return false, err // errore
	}
	return true, nil // id esistente
}
