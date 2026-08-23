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
	UserId string `json:"userId"`
}

type NewUsername struct {
	Name string `json:"name"`
}

type Photo struct {
	Url string `json:"url"`
}

type ChatPreview struct {
	ChatId             string    `json:"id"`
	Name               string    `json:"name"`
	Photo              string    `json:"photo"`
	LastMessageInstant time.Time `json:"lastMessageInstant"`
	LastMessage        string    `json:"lastMessage"`
}

type InfoChat struct {
	IsGroup      bool   `json:"type"`
	Name         string `json:"name"`
	Participants []User `json:"participants"`
}

type Chat struct {
	ChatId       string    `json:"id"`
	IsGroup      bool      `json:"isGroup"`
	Name         string    `json:"name"`
	Photo        string    `json:"photo"`
	Participants []User    `json:"participants"`
	Messages     []Message `json:"messages"`
}

type MessageRequest struct {
	Text string `json:"text"`
}

type Message struct {
	ID        string    `json:"id"`
	Text      string    `json:"text"`
	Instant   time.Time `json:"instant"`
	Comments  []Comment `json:"comments"`
	State     int       `json:"state"`
	Forwarded bool      `json:"forwarded"`
	IdSender  string    `json:"sender"`
	ChatId    string    `json:"chatId"`
}

type ForwardMessageRequest struct {
	UserID string `json:"userId"`
}

type CommentRequest struct {
	SenderId string `json:"sender"`
	Emoji    string `json:"emoji"`
}

type Comment struct {
	ID        string `json:"id"`
	IdSender  string `json:"sender"`
	Emoji     string `json:"emoji"`
	MessageId string `json:"messageId"`
}

type GroupNameRequest struct {
	Name string `json:"name"`
}

type GroupPhotoRequest struct {
	Url string `json:"url"`
}

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

type Group struct {
	Name  string
	Photo string
}

type GroupResponse struct {
	ID    int
	Name  string
	Photo string
}

type AddToGroupRequest struct {
	UserId string
}
