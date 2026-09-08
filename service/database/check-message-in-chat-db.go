package database

import "database/sql"

// Verifica se un messaggio appartiene ad una chat
func (db *appdbimpl) CheckMessageInChat(chatId string, messageId string) (bool, error) {
	var check int
	err := db.c.QueryRow(`	SELECT 1
							FROM messages
							WHERE chatId = ? AND id = ?
							LIMIT 1;
							`, chatId, messageId).Scan(&check)

	if err == sql.ErrNoRows {
		return false, nil // messaggio non appartenente alla chat
	} else if err != nil {
		return false, err // errore
	}
	return true, nil // messaggio appartenente alla chat
}
