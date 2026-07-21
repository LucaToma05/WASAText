package database

import structure "github.com/LucaToma05/WASAText/service/structure"

func (db *appdbimpl) SetGroupPhoto(groupPhoto structure.GroupPhotoRequest, chatId string) error {
	_, err := db.c.Exec(`	UPDATE chats 
							SET photo = ? 
							WHERE id = ? 
							`, groupPhoto.Url, chatId)
	return err
}
