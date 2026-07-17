package api

import (
	"encoding/json"
	"net/http"
	"regexp"

	"github.com/julienschmidt/httprouter"

	"github.com/LucaToma05/WASAText/service/api/reqcontext"
	structure "github.com/LucaToma05/WASAText/service/structure"
)

func (rt *_router) setMyUserName(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	var req structure.NewUsername

	// Read body content
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	// Check username length
	if len(req.Name) < 3 || len(req.Name) > 25 {
		http.Error(w, "Invalid username length", http.StatusBadRequest)
		return
	}

	// Check text pattern
	match, _ := regexp.MatchString("^.*?$", req.Name)
	if !match {
		http.Error(w, "Invalid username pattern", http.StatusBadRequest)
		return
	}

	userId := ctx.UserId
	if userId == "" {
		http.Error(w, "Unauthorized User", http.StatusUnauthorized)
		return
	}

	// Check that the username doesn't exist
	idCheck, err := rt.db.GetIdUserByUsername(req.Name)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	} else if idCheck != "" {
		http.Error(w, "Username already exists", http.StatusBadRequest)
		return
	}

	newUser := structure.User{
		UserId: userId,
		Name:   req.Name,
	}

	err = rt.db.SetMyUserName(newUser)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Create the response
	response := structure.NewUsername{
		Name: req.Name,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
