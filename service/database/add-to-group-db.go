package database

import (
	"time"
)

func (db *appdbimpl) AddToGroup(idGroup int, idUser int) error {
	_, err := db.c.Exec("INSERT INTO userGroups (id_user, id_group, instant) VALUES (?,?,?);", idUser, idGroup, time.Now())
	return err
}
