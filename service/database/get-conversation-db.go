package database

import (
	structure "github.com/LucaToma05/WASAText/service/structure"
)

func (db *appdbimpl) GetConversation(chatId string, currentUserId string) ([]structure.Message, error) {
	if currentUserId != "" {
		// Registra i messaggi letti in message_reads
		_, _ = db.c.Exec(`
			INSERT OR IGNORE INTO message_reads (messageId, userId)
			SELECT id, ?
			FROM messages
			WHERE chatId = ? AND userId != ?;
		`, currentUserId, chatId, currentUserId)

		// Aggiorna lo stato dei messaggi a 2 (Letto) se tutti gli altri partecipanti hanno letto (numero lettori = numero membri gruppo)
		_, _ = db.c.Exec(`
			UPDATE messages
			SET state = 2
			WHERE chatId = ?
			  AND state != 2
			  AND (
			  -- Conta quante persone hanno letto il messaggio
				  SELECT COUNT(*)
				  FROM message_reads mr
				  WHERE mr.messageId = messages.id AND mr.userId != messages.userId
			  ) >= (
			  -- Conta il numero di partecipanti al gruppo
				  SELECT COUNT(*)
				  FROM chat_members cm
				  WHERE cm.chatId = messages.chatId AND cm.userId != messages.userId
			  );
		`, chatId)
	}

	// 1. Recupera tutti i messaggi della conversazione
	rows, err := db.c.Query(`
		SELECT messages.id,
		       messages.text,
		       messages.instant,
		       messages.state,
		       messages.userId,
		       messages.forwarded,
		       messages.replyTo,
		       users.name AS senderName,
		       messages.deleted
		FROM messages
		JOIN users ON users.userId = messages.userId
		WHERE messages.chatId = ?
		ORDER BY messages.instant ASC;
	`, chatId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Inizializza slice vuota
	messages := make([]structure.Message, 0)

	for rows.Next() {
		var msg structure.Message
		var deleted int

		// Mappa ogni colonna della query nella struttura Message
		if err = rows.Scan(
			&msg.ID,
			&msg.Text,
			&msg.Instant,
			&msg.State,
			&msg.IdSender,
			&msg.Forwarded,
			&msg.ReplyTo,
			&msg.SenderUsername,
			&deleted,
		); err != nil {
			return nil, err
		}

		msg.ChatId = chatId
		msg.Reply = (msg.ReplyTo != "")
		msg.Deleted = (deleted == 1)
		if msg.Deleted {
			msg.Text = "Questo messaggio è stato eliminato"
		}
		msg.Comments = make([]structure.Comment, 0)

		messages = append(messages, msg)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	rows.Close()

	// Se non ci sono messaggi, esce subito
	if len(messages) == 0 {
		return messages, nil
	}

	// 2. Recupera tutti i commenti della chat in una mappa map[messageId][]Comment
	commentsMap, err := db.getChatComments(chatId)
	if err == nil {
		// Inserisce i commenti nei rispettivi messaggi
		for i := range messages {
			if comms, ok := commentsMap[messages[i].ID]; ok {
				messages[i].Comments = comms
			}
		}
	}

	return messages, nil
}

// getChatComments recupera tutti i commenti per i messaggi della chat
func (db *appdbimpl) getChatComments(chatId string) (map[string][]structure.Comment, error) {
	rows, err := db.c.Query(`
		SELECT comments.id,
		       comments.emoji,
		       comments.userId,
		       users.name,
		       comments.messageId
		FROM comments
		JOIN users ON users.userId = comments.userId
		JOIN messages ON messages.id = comments.messageId
		WHERE messages.chatId = ?;
	`, chatId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := make(map[string][]structure.Comment)
	// Itera sulle righe del database
	for rows.Next() {
		var c structure.Comment

		if err := rows.Scan(&c.ID, &c.Emoji, &c.IdSender, &c.SenderUsername, &c.MessageId); err == nil {
			// Aggiunge il commento alla mappa
			res[c.MessageId] = append(res[c.MessageId], c)
		}
	}
	return res, rows.Err()
}
