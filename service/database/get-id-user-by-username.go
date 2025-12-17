package database

func (db *appdbimpl) GetIdUserByUsername(username string) (int, error) {
	var id int

	err := db.c.QueryRow("SELECT id FROM users WHERE name=?", username).Scan(&id)
	return id, err
}
