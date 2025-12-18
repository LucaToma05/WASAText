package database

func (db *appdbimpl) CreateConversation() (int, error) {
	aux, err := db.c.Exec("INSERT INTO conversations DEFAULT VALUES")
	if err != nil {
		return 0, err
	}
	conv_id, errore := aux.LastInsertId()
	if errore != nil {
		return 0, errore
	}
	return int(conv_id), nil
}
