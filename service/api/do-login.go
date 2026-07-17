package api

import (
	"crypto/rand"
	"encoding/json"
	"net/http"
	"regexp"

	"github.com/julienschmidt/httprouter"

	"github.com/LucaToma05/WASAText/service/api/reqcontext"
	structure "github.com/LucaToma05/WASAText/service/structure"
)

// If the user does not exist, it will be created, and an identifier is returned.
// If the user exists, the user identifier is returned
func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// defines the User struct
	var req structure.User

	// Read body content
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	// Check username lenght
	if len(req.Name) < 3 || len(req.Name) > 16 {
		http.Error(w, "Invalid name length", http.StatusBadRequest)
		return
	}

	// Check username pattern
	match, _ := regexp.MatchString("^.*?$", req.Name)
	if !match {
		http.Error(w, "Invalid username pattern", http.StatusBadRequest)
		return
	}

	// Search for the user in the database
	id, err := rt.db.GetIdUserByUsername(req.Name)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	} else if id == "" {
		// Generate ID
		id, err = generateId()
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Create the user
		req.UserId = id
		err := rt.db.CreateUser(req)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}

	risposta := structure.LoginResponse{
		UserId: id,
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Authorization", "Bearer "+id)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(risposta)
}

func generateId() (string, error) {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const length = 16

	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	for i := 0; i < length; i++ {
		b[i] = letters[int(b[i])%len(letters)]
	}

	return string(b), nil
}
