package database

import (
	"database/sql"
	"errors"

	structure "github.com/LucaToma05/WASAText/service/structure"
)

// If the user is not found, return id ""
func (db *appdbimpl) GetMyInformation(userId string) (structure.User, error) {
	var info structure.User

	err := db.c.QueryRow("SELECT * FROM users WHERE userId=?", userId).Scan(&info.UserId, &info.Name, &info.Photo)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return structure.User{}, nil // utente non trovato
		}
		return structure.User{}, err // errore DB
	}
	return info, nil // restituisce informazioni
}
