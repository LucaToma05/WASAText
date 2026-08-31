package database

func (db *appdbimpl) CheckChatExistsBetween(userA string, userB string) (bool, error) {
	var count int

	err := db.c.QueryRow(`
        SELECT COUNT(*)
        FROM chats c
        JOIN chat_members cm1 ON c.id = cm1.chatId
        JOIN chat_members cm2 ON c.id = cm2.chatId
        WHERE cm1.userId = ?
          AND cm2.userId = ?
          AND c.isGroup = 0
          AND (
                SELECT COUNT(*) 
                FROM chat_members 
                WHERE chatId = c.id
              ) = 2;
    `, userA, userB).Scan(&count)

	if err != nil {
		return false, err
	}

	return count > 0, nil
}
