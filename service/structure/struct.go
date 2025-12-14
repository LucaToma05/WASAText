package api

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

type SingleChat struct {
	ID    int
	User1 int
	User2 int
}

type Group struct {
	ID    int
	Name  string
	Photo string
}
