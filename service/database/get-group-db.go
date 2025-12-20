package database

import structure "github.com/LucaToma05/WASAText/service/structure"

func (db *appdbimpl) GetGroup(groupId int) structure.Group {
	rows, err := db.c.Query(`	SELECT name, photo
								FROM groups
								WHERE id = ?;`, groupId)
	if err != nil {
		return nil
	}

	var group structure.Group
}
