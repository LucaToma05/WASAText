package database

import structure "github.com/LucaToma05/WASAText/service/structure"

func (db *appdbimpl) SetMyPhoto(photo structure.Photo, userId int) error {
	_, err := db.c.Exec("UPDATE users SET photo = ? WHERE id=?", photo.Url, userId)
	return err
}
