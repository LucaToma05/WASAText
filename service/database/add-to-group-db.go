package database

import (
	structure "github.com/LucaToma05/WASAText/service/structure"
)

func (db *appdbimpl) AddToGroup(req structure.AddToGroupRequest, chatId string) error {
	_, err := db.c.Exec("INSERT INTO userGroups (chatId, userId) VALUES (?,?);", chatId, req.UserId)
	return err
}
