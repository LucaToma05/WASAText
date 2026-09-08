package database

import structure "github.com/LucaToma05/WASAText/service/structure"

func (db *appdbimpl) GetMyConversations(id string) ([]structure.ChatPreview, error) {

	rows, err := db.c.Query(`
        SELECT
			c.id,
			-- Se non è un gruppo, prende il nome dell'altro utente
			CASE 
				WHEN c.isGroup = 0 THEN MAX(u2.name)
				ELSE c.name
			END AS name,

			-- Se non è un gruppo, prende la foto dell'altro utente
			CASE
				WHEN c.isGroup = 0 THEN MAX(u2.photo)
				ELSE c.photo
			END AS photo,

			-- Se l'ultimo messaggio è stato eliminato, mostra 'Questo messaggio è stato eliminato'
			CASE
				WHEN COALESCE(lm.deleted, 0) = 1 THEN 'Questo messaggio è stato eliminato'
				ELSE IFNULL(lm.text, '')
			END AS lastMessage,
			-- Data e ora dell'ultimo messaggio
			IFNULL(strftime('%Y-%m-%dT%H:%M:%SZ', lm.instant), '') AS lastMessageTime,
			-- Indica se la conversazione è un gruppo
			c.isGroup,
			-- Indica se l'utente ha messaggi non letti
			CASE 
				WHEN lm.messageId IS NOT NULL 
				     AND lm.userId != ? 
				     AND NOT EXISTS (
				         SELECT 1 FROM message_reads mr 
				         WHERE mr.messageId = lm.messageId AND mr.userId = ?
				     )
				THEN 1
				ELSE 0
			END AS hasUnread,
			-- Numero di messaggi non letti
			IFNULL((
				SELECT COUNT(*)
				FROM messages m2
				WHERE m2.chatId = c.id
				  AND m2.userId != ?
				  AND NOT EXISTS (
				      SELECT 1 FROM message_reads mr2
				      WHERE mr2.messageId = m2.id AND mr2.userId = ?
				  )
			), 0) AS unreadCount,
			-- ID dell'autore dell'ultimo messaggio
			IFNULL(lm.userId, '') AS lastMessageUserId
		FROM chats c
		JOIN chat_members cm ON cm.chatId = c.id AND cm.userId = ?
		LEFT JOIN chat_members cm2 ON cm2.chatId = c.id AND cm2.userId != ? AND c.isGroup = 0
		LEFT JOIN users u2 ON u2.userId = cm2.userId
		LEFT JOIN (
			SELECT m.chatId, m.text, m.instant, m.userId, m.id AS messageId, COALESCE(m.deleted, 0) AS deleted
			FROM messages m
			JOIN (
				SELECT chatId, MAX(instant) AS maxInstant
				FROM messages
				GROUP BY chatId
			) last ON m.chatId = last.chatId AND m.instant = last.maxInstant
		) lm ON lm.chatId = c.id
		GROUP BY c.id
		ORDER BY lastMessageTime DESC;
    `, id, id, id, id, id, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chats []structure.ChatPreview

	// scan dei campi
	for rows.Next() {
		var chat structure.ChatPreview
		var hasUnreadInt int
		if err = rows.Scan(
			&chat.ChatId,
			&chat.Name,
			&chat.Photo,
			&chat.LastMessage,
			&chat.LastMessageTime,
			&chat.IsGroup,
			&hasUnreadInt,
			&chat.UnreadCount,
			&chat.LastMessageUserId,
		); err != nil {
			return nil, err
		}
		chat.HasUnread = hasUnreadInt > 0 || chat.UnreadCount > 0
		chats = append(chats, chat)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return chats, nil
}
