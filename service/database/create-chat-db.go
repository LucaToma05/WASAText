package database

//import "github.com/LucaToma05/WASAText/service/api"
import structure "github.com/LucaToma05/WASAText/service/structure"

func (db *appdbimpl) CreateChat(infoChat structure.InfoChat, id string, participants []structure.Username) error {
	_, err := db.c.Exec("INSERT INTO chats (id, name, isGroup) VALUES (?, ?, ?)", id, infoChat.Name, infoChat.IsGroup)
	if err != nil {
		return err
	}

	for _, p := range participants {
		_, err := db.c.Exec(`
            INSERT INTO chat_members (chatId, userId)
            VALUES (?, ?)
        `, id, p.Username)
		if err != nil {
			return err
		}
	}
	return nil
}
