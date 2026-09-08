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

type Username struct {
	Username string `json:"username"`
}

type UserId struct {
	UserId string `json:"userId"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type NewUsername struct {
	Name string `json:"name"`
}

type Photo struct {
	Url string `json:"url"`
}

type ChatPreview struct {
	ChatId            string `json:"chatId"`
	Name              string `json:"name"`
	Photo             string `json:"photo"`
	LastMessage       string `json:"lastMessage"`
	LastMessageTime   string `json:"lastMessageTime"`
	IsGroup           bool   `json:"isGroup"`
	HasUnread         bool   `json:"hasUnread"`
	UnreadCount       int    `json:"unreadCount"`
	LastMessageUserId string `json:"lastMessageUserId"`
}

type InfoChat struct {
	IsGroup      bool       `json:"type"`
	Name         string     `json:"name"`
	Photo        string     `json:"photo"`
	Participants []Username `json:"participants"`
}

type Chat struct {
	ChatId       string     `json:"id"`
	IsGroup      bool       `json:"isGroup"`
	Name         string     `json:"name"`
	Photo        string     `json:"photo"`
	Participants []Username `json:"participants"`
	Messages     []Message  `json:"messages"`
}

type MessageRequest struct {
	Text string `json:"text"`
}

type Message struct {
	ID             string    `json:"id"`
	Text           string    `json:"text"`
	Instant        time.Time `json:"instant"`
	Comments       []Comment `json:"comments"`
	State          int       `json:"state"`
	Forwarded      bool      `json:"forwarded"`
	IdSender       string    `json:"sender"`
	SenderUsername string    `json:"senderUsername"`
	SenderName     string    `json:"senderName"`
	ChatId         string    `json:"chatId"`
	ReadBy         []string  `json:"readBy,omitempty"`
	Deleted        bool      `json:"deleted"`
}

type ForwardMessageRequest struct {
	UserID string `json:"userId"`
}

type CommentRequest struct {
	SenderId string `json:"sender"`
	Emoji    string `json:"emoji"`
}

type Comment struct {
	ID             string `json:"id"`
	IdSender       string `json:"sender"`
	SenderUsername string `json:"senderUsername,omitempty"`
	SenderName     string `json:"senderName,omitempty"`
	Emoji          string `json:"emoji"`
	MessageId      string `json:"messageId"`
}

type GroupNameRequest struct {
	Name string `json:"name"`
}

type GroupPhotoRequest struct {
	Url string `json:"url"`
}

type AddToGroupRequest struct {
	UserId string `json:"userId"`
}
