package database

import structure "github.com/LucaToma05/WASAText/service/structure"

func (db *appdbimpl) GetMyConversations(id string) ([]structure.ChatPreview, error) {
	rows, err := db.c.Query(`	SELECT
									c.chatId,
									c.name,
									c.photo,
									m.text AS last_message,
									m.instant AS last_instant
								FROM chats c
								JOIN chat_members cm ON c.chatId = cm.chatId
								LEFT JOIN messages m ON m.chatId = c.chatId AND m.instant = (
									SELECT MAX(instant)
									FROM messages
									WHERE chatId = c.chatId
								)
								WHERE cm.userId = ?
								GROUP BY c.chatId
								ORDER BY last_instant DESC;
								`, id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var chats []structure.ChatPreview

	// Scan the query result
	for rows.Next() {
		var chat structure.ChatPreview
		if err = rows.Scan(&chat.ChatId, &chat.Name, &chat.Photo, &chat.LastMessage, &chat.LastMessageInstant); err != nil {
			return nil, err
		}
		chats = append(chats, chat)
	}

	// rows.Next() non restituisce direttamente un errore, ma lo registra internamente
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return chats, nil
}
