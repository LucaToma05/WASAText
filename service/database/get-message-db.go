package database

func (db *appdbimpl) GetMessage(messageId string) (string, error) {
	var message string
	err := db.c.QueryRow("SELECT text FROM messages WHERE id = ?", messageId).Scan(&message)
	return message, err
}
