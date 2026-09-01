package database

import structure "github.com/LucaToma05/WASAText/service/structure"

func (db *appdbimpl) GetMyConversations(id string) ([]structure.ChatPreview, error) {
	rows, err := db.c.Query(`
        SELECT
            c.id,
            c.name,
            c.photo,
            lm.text AS lastMessage,
            lm.instant AS lastMessageTime
        FROM chats c
        JOIN chat_members cm ON c.id = cm.chatId
        LEFT JOIN (
            SELECT chatId, text, instant
            FROM messages
            WHERE (chatId, instant) IN (
                SELECT chatId, MAX(instant)
                FROM messages
                GROUP BY chatId
            )
        ) lm ON lm.chatId = c.id
        WHERE cm.userId = ?
        ORDER BY lastMessageTime DESC;
    `, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chats []structure.ChatPreview

	for rows.Next() {
		var chat structure.ChatPreview
		if err = rows.Scan(
			&chat.ChatId,
			&chat.Name,
			&chat.Photo,
			&chat.LastMessage,
			&chat.LastMessageTime,
		); err != nil {
			return nil, err
		}
		chats = append(chats, chat)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return chats, nil
}
