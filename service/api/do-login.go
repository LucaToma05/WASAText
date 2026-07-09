package api

import (
	"encoding/json"
	"net/http"
	"regexp"

	"github.com/golang-jwt/jwt/v5"
	"github.com/julienschmidt/httprouter"

	"github.com/LucaToma05/WASAText/service/api/reqcontext"
	structure "github.com/LucaToma05/WASAText/service/structure"
)

var secretKey = []byte("a5jvs496djd15yj")

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
	} else if id == 0 {
		// Create the user
		err := rt.db.CreateUser(req)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}

	id, err = rt.db.GetIdUserByUsername(req.Name)
	if err != nil {
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
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Authorization", tokenString)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(risposta)
}
