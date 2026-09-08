package api

import (
	"encoding/json"
	"net/http"
	"regexp"
	"time"

	"github.com/LucaToma05/WASAText/service/api/reqcontext"
	structure "github.com/LucaToma05/WASAText/service/structure"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) sendMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	userId := ctx.UserId
	if userId == "" {
		http.Error(w, "Unauthorized User", http.StatusUnauthorized)
		return
	}

	chatId := ps.ByName("chatId")
	if chatId == "" {
		http.Error(w, "Parameters extraction error", http.StatusBadRequest)
		return
	}

	ok, err := rt.db.CheckIdConversation(chatId)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if !ok {
		http.Error(w, "Chat not found", http.StatusNotFound)
		return
	}

	ok, err = rt.db.CheckConversation(chatId, userId)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if !ok {
		http.Error(w, "User not authorized for chat", http.StatusForbidden)
		return
	}

	var obj structure.MessageRequest
	err = json.NewDecoder(r.Body).Decode(&obj)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Check text length
	if len(obj.Text) < 1 || len(obj.Text) > 1000 {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Check text pattern
	match, _ := regexp.MatchString("\\S", obj.Text)
	if !match {
		http.Error(w, "Invalid text pattern", http.StatusBadRequest)
		return
	}

	idMessage, err := generateId()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	newMessage := structure.Message{
		ID:        idMessage,
		Text:      obj.Text,
		State:     1,
		Instant:   time.Now(),
		ChatId:    chatId,
		IdSender:  userId,
		Forwarded: false,
	}

	err = rt.db.SendMessage(newMessage)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err = json.NewEncoder(w).Encode(newMessage); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
