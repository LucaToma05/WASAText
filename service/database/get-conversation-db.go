package database

import structure "github.com/LucaToma05/WASAText/service/structure"
import "time"

func (db *appdbimpl) GetConversation(chatId string, currentUserId string) ([]structure.Message, error) {
	// Assicura che la tabella message_reads esista
	_, _ = db.c.Exec(`
		CREATE TABLE IF NOT EXISTS message_reads (
			messageId TEXT NOT NULL,
			userId TEXT NOT NULL,
			readInstant TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (messageId, userId)
		);
	`)

	// Assicura che la colonna deleted esista (soft-delete)
	_, _ = db.c.Exec(`ALTER TABLE messages ADD COLUMN deleted INTEGER NOT NULL DEFAULT 0;`)

	if currentUserId != "" {
		// 1. Registra nel database che currentUserId ha visualizzato tutti i messaggi della chat inviati da altri
		_, _ = db.c.Exec(`
			INSERT OR IGNORE INTO message_reads (messageId, userId)
			SELECT id, ?
			FROM messages
			WHERE chatId = ? AND userId != ?;
		`, currentUserId, chatId, currentUserId)

		// 2. Aggiorna lo stato dei messaggi a 2 (Letto) se e solo se TUTTI i destinatari della chat
		// (tutti i membri di chat_members tranne il mittente del messaggio) hanno visualizzato il messaggio.
		// Funziona perfettamente sia per le chat 1-a-1 (1 destinatario) sia per i gruppi (N destinatari).
		_, _ = db.c.Exec(`
			UPDATE messages
			SET state = 2
			WHERE chatId = ?
			  AND state != 2
			  AND (
				  SELECT COUNT(*)
				  FROM message_reads mr
				  WHERE mr.messageId = messages.id
				    AND mr.userId != messages.userId
			  ) >= (
				  SELECT COUNT(*)
				  FROM chat_members cm
				  WHERE cm.chatId = messages.chatId
				    AND cm.userId != messages.userId
			  )
			  AND (
				  SELECT COUNT(*)
				  FROM chat_members cm2
				  WHERE cm2.chatId = messages.chatId
				    AND cm2.userId != messages.userId
			  ) > 0;
		`, chatId)

		// 3. Normalizza eventuali messaggi con state 0 a state 1 (Ricevuto/Consegnato)
		_, _ = db.c.Exec(`UPDATE messages SET state = 1 WHERE chatId = ? AND state = 0;`, chatId)
	}

    rows, err := db.c.Query(`
        SELECT messages.id,
               messages.text,
               messages.instant,
               messages.state,
               messages.userId,
               messages.forwarded,
               IFNULL(users.name, messages.userId) AS senderName,
               COALESCE(messages.deleted, 0) AS deleted
        FROM messages
        LEFT JOIN users ON users.userId = messages.userId
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
		var senderName string
		var deleted int

		if err = rows.Scan(
			&message.ID,
			&message.Text,
			&instantStr,
			&message.State,
			&message.IdSender,
			&message.Forwarded,
			&senderName,
			&deleted,
		); err != nil {
			return nil, err
		}
		// Se il messaggio è stato eliminato (soft-delete), sovrascriviamo il testo
		// con il tombstone. In questo modo tutti gli utenti vedono la scritta.
		if deleted == 1 {
			message.Text = "Questo messaggio è stato eliminato"
		}

		message.SenderUsername = senderName
		message.SenderName = senderName

		// Se state è ancora 0, normalizza a 1 (Ricevuto/Consegnato)
		if message.State < 1 {
			message.State = 1
		}

		message.Instant, err = time.Parse(time.RFC3339, instantStr)
		if err != nil {
			return nil, err
		}

		if message.Comments, err = db.getComment(message.ID); err != nil {
			return nil, err
		}

		message.ReadBy = db.getMessageReaders(message.ID)

		messages = append(messages, message)
	}

    if err = rows.Err(); err != nil {
        return nil, err
    }

    return messages, nil
}

func (db *appdbimpl) getMessageReaders(messageId string) []string {
	rows, err := db.c.Query(`
		SELECT COALESCE(users.name, message_reads.userId)
		FROM message_reads
		LEFT JOIN users ON users.userId = message_reads.userId
		WHERE message_reads.messageId = ?;
	`, messageId)
	if err != nil {
		return []string{}
	}
	defer rows.Close()

	var readers []string
	for rows.Next() {
		var reader string
		if err := rows.Scan(&reader); err == nil && reader != "" {
			readers = append(readers, reader)
		}
	}
	return readers
}


func (db *appdbimpl) getComment(messageId string) ([]structure.Comment, error) {
	rows, err := db.c.Query(`	SELECT 
									COALESCE(comments.id, ''),
									COALESCE(comments.emoji, ''),
									COALESCE(comments.userId, ''),
									COALESCE(users.name, comments.userId)
								FROM comments
								LEFT JOIN users ON users.userId = comments.userId
								WHERE comments.messageId = ?
								;`, messageId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []structure.Comment

	for rows.Next() {
		var comment structure.Comment
		var senderName string
		if err = rows.Scan(&comment.ID, &comment.Emoji, &comment.IdSender, &senderName); err != nil {
			return nil, err
		}
		comment.SenderUsername = senderName
		comment.SenderName = senderName
		comments = append(comments, comment)
	}
	return comments, nil
}
