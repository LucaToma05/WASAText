package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"regexp"

	"github.com/golang-jwt/jwt/v5"
	"github.com/julienschmidt/httprouter"

	structure "github.com/LucaToma05/WASAText/service/structure"
)

var secretKey = []byte("a5jvs496djd15yj")

// getHelloWorld is an example of HTTP endpoint that returns "Hello world!" as a plain text
func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var req structure.UsernameRequest

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

	//Creation new user
	if id == 0 {
		newUser := structure.User{
			Name:  req.Username,
			Photo: "",
		}
		err := rt.db.CreateUser(newUser)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}

	id, err = rt.db.GetIdUserByUsername(req.Username)
	if err != nil && err != sql.ErrNoRows {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": id,
	})

	tokenString, err := claims.SignedString(secretKey)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	risposta := structure.LoginResponse{
		ID: id,
	}
	w.Header().Set("content-type", "application/json")
	w.Header().Set("Authorization", tokenString)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(risposta)
}
