package database

func (db *appdbimpl) GetIdUserByUsername(username string) (string, error) {
	var id string

	err := db.c.QueryRow("SELECT id FROM users WHERE name=?", username).Scan(&id)
	return id, err
}
