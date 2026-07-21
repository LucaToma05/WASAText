package database

func (db *appdbimpl) DeleteMessage(chatId string, messageId string) error {
	_, err := db.c.Exec(`	DELETE FROM messages
							WHERE id = ? AND id_user = ? AND id_conversation = ?; 
						`)
	return err
}
