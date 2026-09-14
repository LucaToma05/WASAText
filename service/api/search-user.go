package api

import (
	"encoding/json"
	"net/http"

	"github.com/LucaToma05/WASAText/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) searchUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	userId := ctx.UserId
	if userId == "" {
		http.Error(w, "Unauthorized User", http.StatusUnauthorized)
		return
	}

	// Get username from parameters
	username := ps.ByName("username")
	if username == "" {
		http.Error(w, "Parameters extraction error", http.StatusBadRequest)
		return
	}

	if len(username) < 3 || len(username) > 25 {
		http.Error(w, "Invalid username length", http.StatusBadRequest)
		return
	}

	// Get user id from username
	id, err := rt.db.GetIdUserByUsername(username)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if id == "" {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Get user info from user id
	userInfo, err := rt.db.GetMyInformation(id)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(userInfo); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
