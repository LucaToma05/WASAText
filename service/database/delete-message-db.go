package database

func (db *appdbimpl) DeleteMessage(chatId string, messageId string, userId string) error {
	_, err := db.c.Exec(`	DELETE FROM messages
							WHERE id = ? AND userId = ? AND chatId = ?; 
						`)
	return err
}
