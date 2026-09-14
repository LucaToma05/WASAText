package database

import (
	structure "github.com/LucaToma05/WASAText/service/structure"
)

func (db *appdbimpl) AddToGroup(req structure.UserId, chatId string) error {
	_, err := db.c.Exec("INSERT INTO chat_members (chatId, userId) VALUES (?,?);", chatId, req.UserId)
	return err
}
