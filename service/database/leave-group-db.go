package database

func (db *appdbimpl) LeaveGroup(chatId string, userId string) error {
	_, err := db.c.Exec(`	DELETE FROM chat_members
							WHERE chatId = ? AND userId = ? ; 
						`, chatId, userId)
	return err
}
