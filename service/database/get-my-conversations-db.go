package database

import structure "github.com/LucaToma05/WASAText/service/structure"

func (db *appdbimpl) GetMyConversations(id string) []structure.ChatQueryReturn {
	rows, err := db.c.Query(`	SELECT users.name, users.photo, messages.text, messages.instant
								FROM singleChats
						
								JOIN users
									ON users.id = CASE
										WHEN singleChats.id_user_1 = ? THEN singleChats.id_user_2
										ELSE singleChats.id_user_1
									END

								LEFT JOIN messages
									ON messages.id_single_chat = singleChats.id
								AND messages.instant = (
										SELECT MAX(instant)
										FROM messages
										WHERE id_single_chat = singleChats.id
								)

								WHERE singleChats.id_user_1 = ?
								OR singleChats.id_user_2 = ?

								UNION ALL

								SELECT groups.name, groups.photo, messages.text, messages.instant
								FROM groups

								JOIN userGroups
								ON userGroups.id.group = groups.id

								LEFT JOIN messages
								ON messages.id_group = groups.id
								AND messages.instant = (
									SELECT MAX(instant)
									FROM messages
									WHERE id_group = groups.id
								)

								WHERE userGroups.id_user = ?

								ORDER BY messages.instant DESC
								;`, id, id, id, id)
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
