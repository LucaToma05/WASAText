/*
Package database is the middleware between the app database and the code. All data (de)serialization (save/load) from a
persistent database are handled here. Database specific logic should never escape this package.

To use this package you need to apply migrations to the database if needed/wanted, connect to it (using the database
data source name from config), and then initialize an instance of AppDatabase from the DB connection.

For example, this code adds a parameter in `webapi` executable for the database data source name (add it to the
main.WebAPIConfiguration structure):

	DB struct {
		Filename string `conf:""`
	}

This is an example on how to migrate the DB and connect to it:

	// Start Database
	logger.Println("initializing database support")
	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		logger.WithError(err).Error("error opening SQLite DB")
		return fmt.Errorf("opening SQLite: %w", err)
	}
	defer func() {
		logger.Debug("database stopping")
		_ = db.Close()
	}()

Then you can initialize the AppDatabase and pass it to the api package.
*/
package database

import (
	"database/sql"
	"errors"
	"fmt"

	//"github.com/LucaToma05/WASAText/service/api"
	structure "github.com/LucaToma05/WASAText/service/structure"
)

// AppDatabase is the high level interface for the DB
type AppDatabase interface {
	GetName() (string, error)
	SetName(name string) error

	GetIdUserByUsername(username string) (string, error)
	CreateUser(user structure.User) error

	CreateChat(infoChat structure.InfoChat, id string) error

	SetMyUserName(user structure.User) error
	SetMyPhoto(user structure.User) error
	GetMyConversations(id string) ([]structure.ChatPreview, error)
	GetConversation(chatId string) ([]structure.Message, error)
	GetMessage(messageId string) (string, error)
	getComment(messageId string) ([]structure.Comment, error)
	CheckConversation(chatId string, userId string) (bool, error)
	CommentMessage(comment structure.Comment) error
	UncommentMessage(messageId string, commentId string, userId string) error
	IsGroup(chatId string) (bool, error)

	SendMessage(message structure.Message) error
	DeleteMessage(chatId string, messageId string) error
	SetGroupName(name structure.GroupNameRequest, chatId string) error
	SetGroupPhoto(groupPhoto structure.GroupPhotoRequest, chatId string) error

	AddToGroup(req structure.AddToGroupRequest, chatId string) error
	LeaveGroup(chatId string, userId string) error

	Ping() error
}

type appdbimpl struct {
	c *sql.DB
}

// New returns a new instance of AppDatabase based on the SQLite connection `db`.
// `db` is required - an error will be returned if `db` is `nil`.
func New(db *sql.DB) (AppDatabase, error) {
	if db == nil {
		return nil, errors.New("database is required when building a AppDatabase")
	}

	//foreign key activation
	_, errPragma := db.Exec(`PRAGMA foreign_keys = ON`)
	if errPragma != nil {
		return nil, fmt.Errorf("foreign key activation error: %w", errPragma)
	}

	// Check if table exists. If not, the database is empty, and we need to create the structure

	//users table creation
	var users string
	errUsers := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='users';`).Scan(&users)
	if errors.Is(errUsers, sql.ErrNoRows) {
		sqlStmt := `CREATE TABLE users (
						userId TEXT PRIMARY KEY, 
						name TEXT NOT NULL, 
						photo TEXT);`
		_, errUsers = db.Exec(sqlStmt)
		if errUsers != nil {
			return nil, fmt.Errorf("error creating database structure: %w", errUsers)
		}
	}

	//chats table creation
	var chats string
	errChats := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='chats';`).Scan(&chats)
	if errors.Is(errChats, sql.ErrNoRows) {
		sqlStmt := `CREATE TABLE chats (
						id TEXT PRIMARY KEY,
						name TEXT NOT NULL,
						photo TEXT,
						isGroup BOOLEAN);`
		_, errChats = db.Exec(sqlStmt)
		if errChats != nil {
			return nil, fmt.Errorf("error creating database structure: %w", errChats)
		}
	}

	//chats table creation
	var chat_members string
	errConv := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='chat_members';`).Scan(&chat_members)
	if errors.Is(errConv, sql.ErrNoRows) {
		sqlStmt := `CREATE TABLE chat_members (
						chatId TEXT NOT NULL,
						userId TEXT NOT NULL,
						FOREIGN KEY (chatId) REFERENCES chats(id),
						FOREIGN KEY (userId) REFERENCES users(userId),
						PRIMARY KEY (chatId, userId)
					);`
		_, errConv = db.Exec(sqlStmt)
		if errConv != nil {
			return nil, fmt.Errorf("error creating database structure: %w", errConv)
		}
	}

	//messages table creation
	var messages string
	errMessages := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='messages';`).Scan(&messages)
	if errors.Is(errMessages, sql.ErrNoRows) {
		sqlStmt := `CREATE TABLE messages (
						id TEXT PRIMARY KEY, 
						text TEXT NOT NULL, 
						instant TIMESTAMP DEFAULT CURRENT_TIMESTAMP, 
						state INTEGER NOT NULL, 
						forwarded BOOLEAN NOT NULL DEFAULT 0,
						userId TEXT, 
						chatId TEXT,
						FOREIGN KEY (userId) REFERENCES users(userId),
						FOREIGN KEY (chatId) REFERENCES chats(id));`
		_, errMessages = db.Exec(sqlStmt)
		if errMessages != nil {
			return nil, fmt.Errorf("error creating database structure: %w", errMessages)
		}
	}

	//comments table creation
	var comments string
	errComments := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='comments';`).Scan(&comments)
	if errors.Is(errComments, sql.ErrNoRows) {
		sqlStmt := `CREATE TABLE comments (
						id TEXT PRIMARY KEY, 
						emoji TEXT NOT NULL, 
						userId TEXT, 
						messageId TEXT,
						FOREIGN KEY (userId) REFERENCES users(userId),
						FOREIGN KEY (messageId) REFERENCES messages(id));`
		_, errComments = db.Exec(sqlStmt)
		if errComments != nil {
			return nil, fmt.Errorf("error creating database structure: %w", errComments)
		}
	}

	return &appdbimpl{
		c: db,
	}, nil
}

func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}
