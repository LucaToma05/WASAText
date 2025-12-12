package database

import (
	"github.com/LucaToma05/WASAText/service/api"
)

func (db *appdbimpl) CreateUser(user api.User) error {
	_, err := db.c.Exec("INSERT INTO users (id,name,photo) VALUES (?,?,?)", user.ID, user.Name, user.Photo)
	return err
}
