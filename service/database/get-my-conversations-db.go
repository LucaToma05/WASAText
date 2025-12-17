package database

import structure "github.com/LucaToma05/WASAText/service/structure"

func (db *appdbimpl) GetMyConversations(id int) []structure.ChatQueryReturn {
	rows, err := db.c.Query(`
								/* Nome chat */
								CASE
									WHEN g.id IS NOT NULL THEN g.name
									ELSE u.name
								END AS chat_name,

								/* Foto chat */
								CASE
									WHEN g.id IS NOT NULL THEN g.photo
									ELSE u.photo
								END AS chat_photo,

								m.text AS last_message,
								m.instant AS last_message_time

							FROM conversations c

							/* Chat singole */
							LEFT JOIN singleChats sc ON sc.id = c.id
							LEFT JOIN users u ON (
								(u.id = sc.id_user_1 AND sc.id_user_2 = ?)
							OR (u.id = sc.id_user_2 AND sc.id_user_1 = ?)
							)

							/* Gruppi */
							LEFT JOIN groups g ON g.id = c.id
							LEFT JOIN userGroups ug ON ug.id_group = g.id AND ug.id_user = ?

							/* Ultimo messaggio */
							LEFT JOIN messages m ON m.id = (
								SELECT m2.id
								FROM messages m2
								WHERE m2.id_conversation = c.id
								ORDER BY m2.instant DESC
								LIMIT 1
							)

							/* Chat dell'utente */
							WHERE
								sc.id_user_1 = ?
							OR sc.id_user_2 = ?
							OR ug.id_user = ?

							ORDER BY m.instant DESC;`, id, id, id, id, id, id)
	if err != nil {
		return nil
	}

	var chats []structure.ChatQueryReturn

	//Scan the query result
	for rows.Next() {
		var chat structure.ChatQueryReturn
		if err = rows.Scan(&chat.Name, &chat.Photo, &chat.LastMessage, &chat.LastMessageInstant); err != nil {
			return nil
		}
		chats = append(chats, chat)
	}

	//rows.Next() non restituisce direttamente un errore, ma lo registra internamente
	if err = rows.Err(); err != nil {
		return nil
	}

	return chats
}
