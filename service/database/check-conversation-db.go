package database

import "database/sql"

// Verifica se un utente fa parte di una chat
func (db *appdbimpl) CheckConversation(chatId string, userId string) (bool, error) {
	var check int
	err := db.c.QueryRow(`	SELECT 1
							FROM chat_members
							WHERE chatId = ? AND userId = ?
							LIMIT 1;
							`, chatId, userId).Scan(&check)

	if err == sql.ErrNoRows {
		return false, nil // utente non partecipe
	} else if err != nil {
		return false, err // errore
	}
	return true, nil // utente partecipe
}
