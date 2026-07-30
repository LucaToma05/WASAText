package api

import (
	"encoding/json"
	"net/http"
	"regexp"

	"github.com/julienschmidt/httprouter"

	"github.com/LucaToma05/WASAText/service/api/reqcontext"
	structure "github.com/LucaToma05/WASAText/service/structure"
)

func (rt *_router) createChat(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	var req structure.InfoChat

	// Read body content
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	// Check name length
	if len(req.Name) < 3 || len(req.Name) > 25 {
		http.Error(w, "Invalid name length", http.StatusBadRequest)
		return
	}

	// Check text pattern
	match, _ := regexp.MatchString("^(?=.*\\S).+$", req.Name)
	if !match {
		http.Error(w, "Invalid username pattern", http.StatusBadRequest)
		return
	}

	userId := ctx.UserId
	if userId == "" {
		http.Error(w, "Unauthorized User", http.StatusUnauthorized)
		return
	}

	chatId, err := generateId()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	chatErr := rt.db.CreateChat(req, chatId)
	if chatErr != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Create the response
	response := structure.Chat{
		ChatId:       chatId,
		Name:         req.Name,
		IsGroup:      req.IsGroup,
		Participants: req.Participants,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err = json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
