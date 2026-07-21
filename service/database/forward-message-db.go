package database

func (db *appdbimpl) ForwardMessage(chatId string, messageId string) error {
	_, err := db.c.Exec(`	DELETE FROM messages
							WHERE id = ? AND id_user = ? AND id_conversation = ?; 
						`)
	return err
}
