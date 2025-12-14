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
