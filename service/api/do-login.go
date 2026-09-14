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
	var req structure.LoginRequest

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
	match, _ := regexp.MatchString("^[A-Za-z0-9]+$", req.Name)
	if !match {
		http.Error(w, "Invalid username pattern", http.StatusBadRequest)
		return
	}

	// Set the flag to false
	userCreated := false

	// Search for the user in the database
	id, err := rt.db.GetIdUserByUsername(req.Name)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	} else if id == "" { // If the user does not exist
		userCreated = true
		// Generate ID
		id, err = generateId()
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Create the user
		var user structure.User
		user.Name = req.Name
		user.UserId = id
		err = rt.db.CreateUser(user)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}
	risposta := structure.LoginResponse{
		Token: id,
	}
	w.Header().Set("Content-Type", "application/json")
	if userCreated {
		w.WriteHeader(http.StatusCreated)
	} else {
		w.WriteHeader(http.StatusOK)
	}
	if err = json.NewEncoder(w).Encode(risposta); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// Function to generate a random id
func generateId() (string, error) {
	// Characters allowed in id
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	// Length of the id
	const length = 16

	// Allocate a byte slice of the specified length
	b := make([]byte, length)
	// Fill the byte slice with random bytes (from 0 to 255)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	// Map each byte to a character in the allowed characters string
	for i := 0; i < length; i++ {
		b[i] = letters[int(b[i])%len(letters)]
	}

	return string(b), nil
}
