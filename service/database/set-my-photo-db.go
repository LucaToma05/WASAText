package database

import structure "github.com/LucaToma05/WASAText/service/structure"

func (db *appdbimpl) SetMyPhoto(Photo structure.PhotoQuery) error {
	_, err := db.c.Exec("UPDATE users SET photo = ? WHERE id=?", Photo.Photo, Photo.ID)
	return err
}
