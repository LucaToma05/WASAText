package api

import (
	"encoding/json"
	"net/http"
	"net/url"
	"regexp"

	"github.com/LucaToma05/WASAText/service/api/reqcontext"
	structure "github.com/LucaToma05/WASAText/service/structure"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) createChat(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	var req structure.Group

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

	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if len(req.Name) < 3 || len(req.Name) > 50 {
		http.Error(w, "Invalid group name length", http.StatusBadRequest)
		return
	}

	match, _ := regexp.MatchString("^.*?$", req.Name)
	if !match {
		http.Error(w, "Invalid username pattern", http.StatusBadRequest)
		return
	}

	if _, flag := url.ParseRequestURI(req.Photo); flag != nil {
		http.Error(w, "Invalid URI", http.StatusBadRequest)
		return
	}

	idGroup, err := rt.db.CreateConversation()
	if idGroup == 0 {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = rt.db.CreateGroup(req, idGroup)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusBadRequest)
		return
	}

	err = rt.db.AddToGroup(idGroup, userId)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusBadRequest)
		return
	}

	resp := structure.GroupResponse{
		ID:    idGroup,
		Name:  req.Name,
		Photo: req.Photo,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}
