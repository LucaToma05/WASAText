package api

import "time"

type LoginResponse struct {
	ID int
}

type UsernameRequest struct {
	Username string
}

type UsernameQuery struct {
	ID       int
	Username string
}

type PhotoRequest struct {
	Photo string
}

type PhotoQuery struct {
	ID    int
	Photo string
}

type User struct {
	ID    int
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
	ID       int
	Messages []Message
}

type Comment struct {
	ID       int
	IdSender string
	Emoji    string
}

type Message struct {
	ID       int
	Text     string
	Instant  time.Time
	Comments []Comment
	State    int
	ChatId   int
	IdSender int
}

type Group struct {
	Name  string
	Photo string
}

type GroupResponse struct {
	ID    int
	Name  string
	Photo string
}
