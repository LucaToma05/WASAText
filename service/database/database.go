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
)

// AppDatabase is the high level interface for the DB
type AppDatabase interface {
	GetName() (string, error)
	SetName(name string) error

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
						id INTEGER NOT NULL PRIMARY KEY, 
						name TEXT, 
						photo TEXT);`
		_, errUsers = db.Exec(sqlStmt)
		if errUsers != nil {
			return nil, fmt.Errorf("error creating database structure: %w", errUsers)
		}
	}

	//single chats table creation
	var singleChats string
	errSingleChats := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='singleChats';`).Scan(&singleChats)
	if errors.Is(errSingleChats, sql.ErrNoRows) {
		sqlStmt := `CREATE TABLE singleChats (
						id INTEGER NOT NULL PRIMARY KEY, 
						id_user_1 INTEGER,
						id_user_2 INTEGER, 
						FOREIGN KEY (id_user_1) REFERENCES users(id),
						FOREIGN KEY (id_user_2) REFERENCES users(id),
						UNIQUE(id_user_1, id_user_2));`
		_, errSingleChats = db.Exec(sqlStmt)
		if errSingleChats != nil {
			return nil, fmt.Errorf("error creating database structure: %w", errSingleChats)
		}
	}

	//groups table creation
	var groups string
	errGroups := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='groups';`).Scan(&groups)
	if errors.Is(errGroups, sql.ErrNoRows) {
		sqlStmt := `CREATE TABLE groups (
						id INTEGER NOT NULL PRIMARY KEY, 
						name TEXT, 
						photo TEXT);`
		_, errGroups = db.Exec(sqlStmt)
		if errGroups != nil {
			return nil, fmt.Errorf("error creating database structure: %w", errGroups)
		}
	}

	//user_groups table creation (association table)
	var userGroups string
	errUserGroups := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='userGroups';`).Scan(&userGroups)
	if errors.Is(errUserGroups, sql.ErrNoRows) {
		sqlStmt := `CREATE TABLE userGroups (
						id_user INTEGER NOT NULL, 
						id_group INTEGER NOT NULL,
						instant TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
						FOREIGN KEY (id_user) REFERENCES users(id),
						FOREIGN KEY (id_group) REFERENCES groups(id),
						PRIMARY KEY (id_user, id_group));`
		_, errUserGroups = db.Exec(sqlStmt)
		if errUserGroups != nil {
			return nil, fmt.Errorf("error creating database structure: %w", errUserGroups)
		}
	}

	//messages table creation
	var messages string
	errMessages := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='messages';`).Scan(&messages)
	if errors.Is(errMessages, sql.ErrNoRows) {
		sqlStmt := `CREATE TABLE messages (
						id INTEGER NOT NULL PRIMARY KEY, 
						text TEXT, 
						instant TIMESTAMP DEFAULT CURRENT_TIMESTAMP, 
						state INTEGER, 
						id_user INTEGER, 
						id_single_chat INTEGER,
						id_group INTEGER,
						FOREIGN KEY (id_user) REFERENCES users(id),
						FOREIGN KEY (id_single_chat) REFERENCES singleChats(id),
						FOREIGN KEY (id_group) REFERENCES groups(id));`
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
						id INTEGER NOT NULL PRIMARY KEY, 
						text TEXT, 
						id_user INTEGER, 
						id_message INTEGER,
						FOREIGN KEY (id_user) REFERENCES users(id),
						FOREIGN KEY (id_message) REFERENCES messages(id));`
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
