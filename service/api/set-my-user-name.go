package api

import (
	"encoding/json"
	"net/http"
	"regexp"

	"github.com/julienschmidt/httprouter"

	structure "github.com/LucaToma05/WASAText/service/structure"
)

func (rt *_router) setMyUserName(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	var req structure.User

	// Read body content
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	// Check username length
	if len(req.Name) < 3 || len(req.Name) > 16 {
		http.Error(w, "Invalid username length", http.StatusBadRequest)
		return
	}

	// Check text pattern
	match, _ := regexp.MatchString("^.*?$", req.Name)
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
		return
	}

	// Add user id
	req.ID = userId

	// Check that the username doesn't exist
	idCheck, err := rt.db.GetIdUserByUsername(req.Name)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	} else if idCheck != 0 {
		http.Error(w, "Username already exists", http.StatusBadRequest)
		return
	}

	err = rt.db.SetMyUserName(req)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Create the response
	response := structure.UsernameResponse{
		Name: req.Name,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
