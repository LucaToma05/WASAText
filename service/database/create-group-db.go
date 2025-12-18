package database

import structure "github.com/LucaToma05/WASAText/service/structure"

func (db *appdbimpl) CreateGroup(group structure.Group, idGroup int) error {
	_, err := db.c.Exec("INSERT INTO groups (id, name, photo) VALUES (?,?,?);", idGroup, group.Name, group.Photo)
	return err
}
