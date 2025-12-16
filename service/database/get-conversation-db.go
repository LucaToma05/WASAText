package database

import structure "github.com/LucaToma05/WASAText/service/structure"

func (db *appdbimpl) GetConversation(userId string, chatId string) []structure.Message {
	rows, err := db.c.Query(`	SELECT messages.id, messages.text, messages.instant, messages.state
								FROM singleChats
								JOIN messages ON messages.id_single_chat = ? 
								WHERE singleChats.id = ? AND (singleChats.id_user_1 = ? OR singleChats.id_user_2 = ?)
								;`, chatId, chatId, userId, userId)
	if err != nil {
		return nil
	}

	var messages []structure.Message

	for rows.Next() {
		var message structure.Message
		if err = rows.Scan(&message.ID, &message.Text, &message.Instant, &message.State); err != nil {
			return nil
		}
		if message.Comments = db.getComment(message.ID); message.Comments == nil {
			return nil
		}
		messages = append(messages, message)
	}

	return messages
}

func (db *appdbimpl) getComment(messageId string) []structure.Comment {
	rows, err := db.c.Query(`	SELECT id, text, id_user
								FROM comments
								WHERE id_message = ?
								;`, messageId)
	if err != nil {
		return nil
	}

	var comments []structure.Comment

	for rows.Next() {
		var comment structure.Comment
		if err = rows.Scan(&comment.ID, &comment.Emoji, &comment.IdSender); err != nil {
			return nil
		}
		comments = append(comments, comment)
	}
	return comments
}
