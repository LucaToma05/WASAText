package api

import (
	"encoding/json"
	"net/http"
	"regexp"

	"github.com/julienschmidt/httprouter"

	structure "github.com/LucaToma05/WASAText/service/structure"
)

func (rt *_router) setMyUserName(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

	match, _ := regexp.MatchString("^.*?$", req.Username)
	if !match {
		http.Error(w, "Invalid username pattern", http.StatusBadRequest)
		return
	}

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Unauthorized User", http.StatusUnauthorized)
		return
	}
	userId, err := extractUserId(authHeader)
	if err != nil {
		http.Error(w, "Token extraction error", http.StatusInternalServerError)
	}

	newUser := structure.UsernameQuery{
		ID:       userId,
		Username: req.Username,
	}

	err = rt.db.SetMyUserName(newUser)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

	//req è la risposta
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(req)
}
