package api

import (
	"encoding/json"
	"net/http"

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
	ctx.Logger.Info("DEBUG: ctx.UserId = " + ctx.UserId)

	userId := ctx.UserId
	if userId == "" {
		http.Error(w, "Unauthorized User", http.StatusUnauthorized)
		return
	}
	/*
		// Check text pattern
		match, _ := regexp.MatchString("^(?=.*\\S).+$", req.Name)
		if !match {
			http.Error(w, "Invalid username pattern", http.StatusBadRequest)
			return
		}
	*/

	id, err := rt.db.GetIdUserByUsername(req.Name)
	if err != nil {
		http.Error(w, "Internal 1 Server Error", http.StatusInternalServerError)
		return
	}
	if id == "" {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	exists, err := rt.db.CheckChatExistsBetween(userId, id)
	if err != nil {
		http.Error(w, "Internal 2 Server Error", http.StatusInternalServerError)
		return
	}
	if exists {
		http.Error(w, "Chat already exists", http.StatusConflict)
		return
	}

	chatId, err := generateId()
	if err != nil {
		http.Error(w, "Internal 3 Server Error", http.StatusInternalServerError)
		return
	}

	participants := []structure.UserId{
		{UserId: userId},
		{UserId: id},
	}

	chatErr := rt.db.CreateChat(req, chatId, participants)
	if chatErr != nil {
		http.Error(w, "Internal 4 Server Error", http.StatusInternalServerError)
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
		http.Error(w, "Internal 5 Server Error", http.StatusInternalServerError)
		return
	}
}
