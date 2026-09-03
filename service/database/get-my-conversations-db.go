package database

import structure "github.com/LucaToma05/WASAText/service/structure"

func (db *appdbimpl) GetMyConversations(id string) ([]structure.ChatPreview, error) {
	rows, err := db.c.Query(`
        SELECT
			c.id,
			CASE 
				WHEN c.isGroup = 0 THEN u2.name
				ELSE c.name
			END AS name,

			CASE
				WHEN c.isGroup = 0 THEN u2.photo
				ELSE c.photo
			END AS photo,

			IFNULL(lm.text, '') AS lastMessage,
			IFNULL(strftime('%Y-%m-%d %H:%M:%S', lm.instant), '') AS lastMessageTime
		FROM chats c
		JOIN chat_members cm ON cm.chatId = c.id AND cm.userId = ?
		LEFT JOIN chat_members cm2 ON cm2.chatId = c.id AND cm2.userId != ?
		LEFT JOIN users u2 ON u2.userId = cm2.userId
		LEFT JOIN (
			SELECT m.chatId, m.text, m.instant
			FROM messages m
			JOIN (
				SELECT chatId, MAX(instant) AS maxInstant
				FROM messages
				GROUP BY chatId
			) last ON m.chatId = last.chatId AND m.instant = last.maxInstant
		) lm ON lm.chatId = c.id
		ORDER BY lastMessageTime DESC;



    `, id, id)
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
