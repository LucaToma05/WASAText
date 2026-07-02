package database

import structure "github.com/LucaToma05/WASAText/service/structure"

func (db *appdbimpl) GetConversation(userId int, chatId int) ([]structure.Message, error) {
	rows, err := db.c.Query(`	SELECT messages.id, messages.text, messages.instant, messages.state
								FROM messages 
								JOIN conversations  ON messages.id_conversation = conversations.id
								LEFT JOIN singleChats  ON singleChats.id = conversations.id
								LEFT JOIN groups ON groups.id = conversations.id
								LEFT JOIN userGroups  ON userGroups.id_group = groups.id
								WHERE messages.id_conversation = ?
								AND (
									singleChats.id_user_1 = ?
									OR singleChats.id_user_2 = ?
									OR userGroups.id_user = ?
								)
								ORDER BY messages.instant;
								`, chatId, userId, userId, userId)
	if err != nil {
		return nil, err
	}

	var messages []structure.Message

	for rows.Next() {
		var message structure.Message
		if err = rows.Scan(&message.ID, &message.Text, &message.Instant, &message.State); err != nil {
			return nil, err
		}
		if message.Comments, err = db.getComment(message.ID); err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}

	return messages, nil
}

func (db *appdbimpl) getComment(messageId int) ([]structure.Comment, error) {
	rows, err := db.c.Query(`	SELECT id, text, id_user
								FROM comments
								WHERE id_message = ?
								;`, messageId)
	if err != nil {
		return nil, err
	}

	var comments []structure.Comment

	for rows.Next() {
		var comment structure.Comment
		if err = rows.Scan(&comment.ID, &comment.Emoji, &comment.IdSender); err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	return comments, nil
}
