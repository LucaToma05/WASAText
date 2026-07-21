package database

import structure "github.com/LucaToma05/WASAText/service/structure"

func (db *appdbimpl) CommentMessage(comment structure.Comment) error {
	_, err := db.c.Exec(`INSERT INTO comments(id,emoji,userId,messageId)
						VALUES (?,?,?,?);`, comment.ID, comment.Emoji, comment.IdSender, comment.MessageId)
	return err
}
