package database

import structure "github.com/LucaToma05/WASAText/service/structure"

func (db *appdbimpl) SetMyUserName(User structure.User) error {
	_, err := db.c.Exec("UPDATE users SET name = ? WHERE id=?", User.Name, User.ID)
	return err
}
