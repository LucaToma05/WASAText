package api

import "time"

type UsernameRequest struct {
	Username string
}

type UsernameQuery struct {
	ID       string
	Username string
}

type PhotoRequest struct {
	Photo string
}

type PhotoQuery struct {
	ID    string
	Photo string
}

type User struct {
	ID    string
	Name  string
	Photo string
}

type ChatQueryReturn struct {
	Name               string
	Photo              string
	LastMessageInstant time.Time
	LastMessage        string
}

type Conversation struct {
	ID       string
	Messages []Message
}

type Comment struct {
	ID       string
	IdSender string
	Emoji    string
}

type Message struct {
	ID       string
	Text     string
	Instant  time.Time
	Comments []Comment
	State    int
}
