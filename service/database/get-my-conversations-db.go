package database

import structure "github.com/LucaToma05/WASAText/service/structure"

func (db *appdbimpl) GetMyConversations(id int) ([]structure.ChatPreview, error) {
	rows, err := db.c.Query(`	-- Single chats
								SELECT
									-- sc.id AS conversation_id,
									u.name AS chat_name,
									u.photo AS chat_photo,
									m.text AS last_text,
									m.instant AS last_instant
								FROM singleChats sc
								JOIN users u ON u.id = CASE 
														WHEN sc.id_user_1 = ? THEN sc.id_user_2
														ELSE sc.id_user_1
														END
								LEFT JOIN (
									SELECT id_conversation, text, instant
									FROM messages m1
									WHERE instant = (
										SELECT MAX(instant)
										FROM messages m2
										WHERE m2.id_conversation = m1.id_conversation
									)
								) m ON m.id_conversation = sc.id
								WHERE sc.id_user_1 = ? OR sc.id_user_2 = ?

								UNION ALL

								-- Groups
								SELECT
									-- g.id AS conversation_id,
									g.name AS chat_name,
									g.photo AS chat_photo,
									m.text AS last_text,
									m.instant AS last_instant
								FROM groups g
								JOIN userGroups ug ON ug.id_group = g.id AND ug.id_user = ?
								LEFT JOIN (
									SELECT id_conversation, text, instant
									FROM messages m1
									WHERE instant = (
										SELECT MAX(instant)
										FROM messages m2
										WHERE m2.id_conversation = m1.id_conversation
									)
								) m ON m.id_conversation = g.id

								ORDER BY last_instant DESC;
								`, id, id, id, id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var chats []structure.ChatPreview

	// Scan the query result
	for rows.Next() {
		var chat structure.ChatPreview
		if err = rows.Scan(&chat.Name, &chat.Photo, &chat.LastMessage, &chat.LastMessageInstant); err != nil {
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
