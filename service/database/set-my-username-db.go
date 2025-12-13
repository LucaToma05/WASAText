package database

func (db *appdbimpl) setMyUserName(Username string) (string, error) {
	var name string
	err := db.c.QueryRow("SELECT name FROM example_table WHERE id=1").Scan(&name)
	return name, err
}
