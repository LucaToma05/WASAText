package database

import "database/sql"

func (db *appdbimpl) CheckCommentInMessage(commentId string, messageId string) (bool, error) {
	var check int
	err := db.c.QueryRow(`	SELECT 1
							FROM comments
							WHERE id = ? AND messageId = ?
							LIMIT 1;
							`, commentId, messageId).Scan(&check)

	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}
