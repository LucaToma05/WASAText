package database

import structure "github.com/LucaToma05/WASAText/service/structure"

func (db *appdbimpl) SendMessage(message structure.Message) error {
	_, err := db.c.Exec(`INSERT INTO messages(id,text,instant,state,userId,chatId, forwarded)
						VALUES (?,?,?,?,?,?,?);`, message.ID, message.Text, message.Instant, message.State, message.IdSender, message.ChatId, message.Forwarded)
	return err
}
