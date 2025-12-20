package database

import structure "github.com/LucaToma05/WASAText/service/structure"

func (db *appdbimpl) SetGroupName(name structure.GroupNameRequest, userId int, groupId int) error {
	_, err := db.c.Exec(`	UPDATE groups 
							SET name = ? 
							WHERE id = ? 
							AND EXISTS (
								SELECT 1
								FROM userGrops
								WHERE id_group = ? AND id_user = ?)`, name.Name, groupId, groupId, userId)
	return err
}
