package database

import structure "github.com/LucaToma05/WASAText/service/structure"

func (db *appdbimpl) SetGroupName(groupName structure.GroupNameRequest, chatId string) error {
	_, err := db.c.Exec(`	UPDATE chats 
							SET name = ? 
							WHERE id = ? 
							`, groupName.Name, chatId)
	return err
}
