package api

import "time"

type LoginRequest struct { // UTILIZZATA IN do-login.go
	Name string `json:"name"`
}

type User struct { // UTILIZZATA IN do-login.go
	UserId string `json:"userId"`
	Name   string `json:"name"`
	Photo  string `json:"photo"`
}

type LoginResponse struct { //UTILIZZATA IN do-login.go
	UserId string `json:"userId"`
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

type ForwardMessageRequest struct {
	UserID string
}

type CommentRequest struct {
	SenderId string
	Emoji    string
}

type Comment struct {
	ID        string
	IdSender  string
	Emoji     string
	MessageId string
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

type Message struct {
	ID          string
	Text        string
	Instant     time.Time
	Comments    []Comment
	State       int
	Forwarded   bool
	IdSender    string
	NameSender  string
	PhotoSender string
	ChatId      string
}

type MessageRequest struct {
	Text string
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

type GroupPhotoRequest struct {
	Url string
}

type AddToGroupRequest struct {
	UserId string
}
