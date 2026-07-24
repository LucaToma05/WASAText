package database

import structure "github.com/LucaToma05/WASAText/service/structure"

func (db *appdbimpl) SetMyUserName(user structure.User) error {
	_, err := db.c.Exec("UPDATE users SET name = ? WHERE userId=?", user.Name, user.UserId)
	return err
}
