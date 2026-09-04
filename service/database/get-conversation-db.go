package database

import structure "github.com/LucaToma05/WASAText/service/structure"
import "time"

func (db *appdbimpl) GetConversation(chatId string) ([]structure.Message, error) {
    rows, err := db.c.Query(`
        SELECT messages.id,
               messages.text,
               messages.instant,
               messages.state,
               messages.userId,
               messages.forwarded
        FROM messages
        WHERE messages.chatId = ?
        ORDER BY messages.instant ASC;
    `, chatId)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var messages []structure.Message

    for rows.Next() {
		var message structure.Message
		var instantStr string

		if err = rows.Scan(
			&message.ID,
			&message.Text,
			&instantStr,
			&message.State,
			&message.IdSender,
			&message.Forwarded,
		); err != nil {
			return nil, err
		}

		message.Instant, err = time.Parse(time.RFC3339, instantStr)
		if err != nil {
			return nil, err
		}

		if message.Comments, err = db.getComment(message.ID); err != nil {
			return nil, err
		}

		messages = append(messages, message)
	}


    if err = rows.Err(); err != nil {
        return nil, err
    }

    return messages, nil
}


func (db *appdbimpl) getComment(messageId string) ([]structure.Comment, error) {
	rows, err := db.c.Query(`	SELECT 
									COALESCE(id, ''),
									COALESCE(emoji, ''),
									COALESCE(userId, '')
								FROM comments
								WHERE messageId = ?
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
