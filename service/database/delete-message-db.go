package database

func (db *appdbimpl) DeleteMessage(chatId string, messageId string, userId string) error {
	// Aggiunge la colonna deleted se non esiste ancora (idempotente)
	_, _ = db.c.Exec(`ALTER TABLE messages ADD COLUMN deleted INTEGER NOT NULL DEFAULT 0;`)

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
