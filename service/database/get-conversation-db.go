package database

import structure "github.com/LucaToma05/WASAText/service/structure"

func (db *appdbimpl) GetConversation(chatId string) ([]structure.Message, error) {
	rows, err := db.c.Query(`	SELECT messages.id, 
									messages.text,
									messages.instant,
									messages.state,
									messages.userId,
									messages.forwarded,
									users.name,
									users.photo
								FROM messages
								JOIN users ON messages.userId = users.userId
								WHERE messages.chatId = ?
								ORDER BY messages.instant DESC;
								`, chatId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var messages []structure.Message

	for rows.Next() {
		var message structure.Message
		if err = rows.Scan(&message.ID, &message.Text, &message.Instant, &message.State, &message.IdSender, &message.Forwarded); err != nil {
			return nil, err
		}
		if message.Comments, err = db.getComment(message.ID); err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}

	return messages, nil
}

func (db *appdbimpl) getComment(messageId string) ([]structure.Comment, error) {
	rows, err := db.c.Query(`	SELECT id, emoji, userId
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
