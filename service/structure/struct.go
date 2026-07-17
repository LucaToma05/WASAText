package api

import "time"

type LoginRequest struct {
	Name string `json:"name"`
}

type User struct {
	UserId string `json:"userId"`
	Name   string `json:"name"`
	Photo  string `json:"photo"`
}

type LoginResponse struct {
	UserId string
}

type NewUsername struct {
	Name string `json:"name"`
}

type Photo struct {
	Url string
}

type ChatPreview struct {
	ChatId             string
	Name               string
	Photo              string
	LastMessageInstant time.Time
	LastMessage        string
}

type InfoChat struct {
	IsGroup      bool
	Name         string
	Participants []User
}

type Chat struct {
	ChatId       string
	IsGroup      bool
	Name         string
	Photo        string
	Participants []User
	Messages     []Message
}

// ##################################

type UsernameRequest struct {
	Username string
}

type UsernameQuery struct {
	ID       int
	Username string
}

type PhotoQuery struct {
	ID    int
	Photo string
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

type GroupNameRequest struct {
	Name string
}
