package database

//import "github.com/LucaToma05/WASAText/service/api"
import structure "github.com/LucaToma05/WASAText/service/structure"

func (db *appdbimpl) CreateUser(user structure.User) error {
	_, err := db.c.Exec("INSERT INTO users (id,name,photo) VALUES (?,?,?)", user.ID, user.Name, user.Photo)
	return err
}
