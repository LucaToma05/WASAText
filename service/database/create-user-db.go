package database

//import "github.com/LucaToma05/WASAText/service/api"
import structure "github.com/LucaToma05/WASAText/service/structure"

func (db *appdbimpl) CreateUser(user structure.User) error {
	_, err := db.c.Exec("INSERT INTO users (userId, name, photo) VALUES (?, ?, ?)", user.UserId, user.Name, user.Photo)
	return err
}
