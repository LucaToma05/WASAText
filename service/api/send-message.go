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

	var obj structure.MessageRequest
	err := json.NewDecoder(r.Body).Decode(&obj)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	//Check text length
	if len(obj.Text) < 1000 {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	//Check text pattern
	match, _ := regexp.MatchString("^.*?$", obj.Text)
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
		ID:       idMessage,
		Text:     obj.Text,
		State:    0,
		Instant:  time.Now(),
		ChatId:   chatId,
		IdSender: userId,
	}

	err = rt.db.SendMessage(newMessage)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(obj)
}
