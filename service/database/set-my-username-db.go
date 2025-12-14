package database

import structure "github.com/LucaToma05/WASAText/service/structure"

func (db *appdbimpl) SetMyUserName(User structure.UsernameQuery) error {
	_, err := db.c.Exec("UPDATE users SET name = ? WHERE id=?", User.Username, User.ID)
	return err
}
