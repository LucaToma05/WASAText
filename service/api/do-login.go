package api

import (
	"database/sql"
	"encoding/json"
	"math/rand"
	"net/http"
	"regexp"
	"time"

	"github.com/julienschmidt/httprouter"
)

// getHelloWorld is an example of HTTP endpoint that returns "Hello world!" as a plain text
func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var req UsernameRequest

	//Read body content
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	//Check username lenght
	if len(req.Username) < 3 || len(req.Username) > 25 {
		http.Error(w, "Invalid username length", http.StatusBadRequest)
		return
	}

	//Check username pattern
	match, _ := regexp.MatchString("^.*?$", req.Username)
	if !match {
		http.Error(w, "Invalid username pattern", http.StatusBadRequest)
		return
	}

	//User search in database
	id, err := rt.db.GetIdUserByUsername(req.Username)
	if err != nil && err != sql.ErrNoRows {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if id == "" {
		newId := RandString()
		id = newId
		newUser := User{newId, req.Username, ""}
		err := rt.db.CreateUser(newUser)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}

	risposta := UsernameRequest{id}
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(risposta)
}

func RandString() string {
	rand.Seed(time.Now().UnixNano())
	stringa := make([]byte, 10)
	for i := range stringa {
		if rand.Intn(2) == 0 {
			stringa[i] = byte(rand.Intn(26) + 65)
		} else {
			stringa[i] = byte(rand.Intn(26) + 97)
		}
	}
	return string(stringa)
}
