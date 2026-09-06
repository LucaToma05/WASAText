package database

import "database/sql"

func (db *appdbimpl) CheckMessageInChat(chatId string, messageId string) (bool, error) {
	var check int
	err := db.c.QueryRow(`	SELECT 1
							FROM messages
							WHERE chatId = ? AND id = ?
							LIMIT 1;
							`, chatId, messageId).Scan(&check)

	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}
