package database

func (db *appdbimpl) UncommentMessage(messageId string, commentId string, userId string) error {
	_, err := db.c.Exec(`	DELETE FROM comments
							WHERE id = ? AND userId = ? AND messageId = ?; 
						`, commentId, userId, messageId)
	return err
}
