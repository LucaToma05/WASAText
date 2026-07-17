package database

//import "github.com/LucaToma05/WASAText/service/api"
import structure "github.com/LucaToma05/WASAText/service/structure"

func (db *appdbimpl) CreateChat(infoChat structure.InfoChat, id string) error {
	_, err := db.c.Exec("INSERT INTO chats (id, name, isGroup) VALUES (?, ?, ?)", id, infoChat.Name, infoChat.IsGroup)
	return err
}
