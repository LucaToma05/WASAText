package database

import "database/sql"

// Verifica se un commento appartiene a un messaggio
func (db *appdbimpl) CheckCommentInMessage(commentId string, messageId string) (bool, error) {
	var check int
	err := db.c.QueryRow(`	SELECT 1
							FROM comments
							WHERE id = ? AND messageId = ?
							LIMIT 1;
							`, commentId, messageId).Scan(&check)

	if err == sql.ErrNoRows {
		return false, nil // commento non appartenente al messaggio
	} else if err != nil {
		return false, err // errore
	}
	return true, nil // commento appartenente al messaggio
}
