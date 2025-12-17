package database

import structure "github.com/LucaToma05/WASAText/service/structure"

func (db *appdbimpl) SendMessage(message structure.Message) error {
	_, err := db.c.Exec(`INSERT INTO messages(text,instant,state,id_user,id_conversation)
						VALUES (?,?,?,?,?);`, message.Text, message.Instant, message.State, message.IdSender, message.ChatId)
	return err
}
