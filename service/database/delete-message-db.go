package database

func (db *appdbimpl) DeleteMessage(chatId string, messageId string, userId string) error {

	// Soft-delete: marca il messaggio come eliminato invece di cancellarlo fisicamente.
	// In questo modo il tombstone rimane visibile a tutti gli utenti della chat
	// e non si creano problemi di foreign key con message_reads.
	_, err := db.c.Exec(`
		UPDATE messages
		SET deleted = 1
		WHERE id = ? AND userId = ? AND chatId = ?;
	`, messageId, userId, chatId)
	return err
}
