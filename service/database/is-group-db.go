package database

// Verifica se una chat è un gruppo
func (db *appdbimpl) IsGroup(chatId string) (bool, error) {
	var isGroup bool
	err := db.c.QueryRow(`	SELECT isGroup
							FROM chats
							WHERE id = ?`, chatId).Scan(&isGroup)
	if err != nil {
		return false, err // errore
	}
	return isGroup, nil // restituisce risultato
}
