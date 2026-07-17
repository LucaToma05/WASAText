package database

import structure "github.com/LucaToma05/WASAText/service/structure"

func (db *appdbimpl) SetMyPhoto(user structure.User) error {
	_, err := db.c.Exec("UPDATE users SET photo = ? WHERE userId=?", user.Photo, user.UserId)
	return err
}
