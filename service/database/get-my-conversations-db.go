package database

import structure "github.com/LucaToma05/WASAText/service/structure"

func (db *appdbimpl) GetMyConversations(id string) ([]structure.ChatPreview, error) {

	rows, err := db.c.Query(`
		SELECT
			c.id,
			-- Se non è un gruppo, prende nome e foto dell'altro utente
			CASE WHEN c.isGroup = 0 THEN MAX(u2.name) ELSE c.name END AS name,
			CASE WHEN c.isGroup = 0 THEN MAX(u2.photo) ELSE c.photo END AS photo,

			-- Se l'ultimo messaggio è stato eliminato, mostra il testo informativo
			CASE 
				WHEN lm.deleted = 1 THEN 'Questo messaggio è stato eliminato' 
				ELSE IFNULL(lm.text, '') 
			END AS lastMessage,

			-- Data dell'ultimo messaggio
			IFNULL(strftime('%Y-%m-%dT%H:%M:%SZ', lm.instant), '') AS lastMessageTime,
			c.isGroup,

			-- Indica se ci sono messaggi non letti per l'utente
			EXISTS (
				SELECT 1
				FROM messages m2
				LEFT JOIN message_reads mr ON mr.messageId = m2.id AND mr.userId = ?
				WHERE m2.chatId = c.id
				  AND m2.userId != ?
				  AND mr.messageId IS NULL
			) AS hasUnread,

			-- ID autore dell'ultimo messaggio
			IFNULL(lm.userId, '') AS lastMessageUserId
		FROM chats c
		JOIN chat_members cm ON cm.chatId = c.id AND cm.userId = ?
		LEFT JOIN chat_members cm2 ON cm2.chatId = c.id AND cm2.userId != ? AND c.isGroup = 0
		LEFT JOIN users u2 ON u2.userId = cm2.userId
		LEFT JOIN (
			SELECT chatId, text, instant, userId, id AS messageId, deleted, MAX(instant)
			FROM messages
			GROUP BY chatId
		) lm ON lm.chatId = c.id
		GROUP BY c.id
		ORDER BY lastMessageTime DESC;
	`, id, id, id, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Crea lo slice di chat previews
	chats := make([]structure.ChatPreview, 0)

	// Inserisce i dati dentro lo slice
	for rows.Next() {
		var chat structure.ChatPreview
		if err = rows.Scan(
			&chat.ChatId,
			&chat.Name,
			&chat.Photo,
			&chat.LastMessage,
			&chat.LastMessageTime,
			&chat.IsGroup,
			&chat.HasUnread,
			&chat.LastMessageUserId,
		); err != nil {
			return nil, err
		}
		chats = append(chats, chat)
	}

	return chats, rows.Err()
}
