package database

import "database/sql"

func (db *appdbimpl) CheckConversation(chatId int) (bool, error) {
	x := 0
	err := db.c.QueryRow(`	SELECT 1
							FROM conversations
							LEFT JOIN singleChats  ON singleChats.id = conversations.id
							LEFT JOIN groups ON groups.id = conversations.id
							LEFT JOIN userGroups ON userGroups.id_group = groups.id
							WHERE conversations.id = ?
							AND (
								singleChats.id_user_1 = ?
								OR singleChats.id_user_2 = ?
								OR userGroups.id_user = ?
							)
							LIMIT 1;
							`, chatId, chatId, chatId, chatId).Scan(&x)

	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil

}
