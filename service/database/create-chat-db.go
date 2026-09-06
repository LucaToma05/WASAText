package database

//import "github.com/LucaToma05/WASAText/service/api"
import structure "github.com/LucaToma05/WASAText/service/structure"

func (db *appdbimpl) CreateChat(infoChat structure.InfoChat, id string, participants []structure.UserId) error {
	_, err := db.c.Exec("INSERT INTO chats (id, name, photo, isGroup) VALUES (?, ?, ?, ?)", id, infoChat.Name, infoChat.Photo, infoChat.IsGroup)
	if err != nil {
		return err
	}

	for _, p := range participants {
		_, err := db.c.Exec(`
            INSERT INTO chat_members (chatId, userId)
            VALUES (?, ?)
        `, id, p.UserId)
		if err != nil {
			return err
		}
	}
	return nil
}
