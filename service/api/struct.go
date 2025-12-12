package api

type UsernameRequest struct {
	Username string
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
