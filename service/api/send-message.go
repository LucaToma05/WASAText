package api

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strconv"
	"time"

	structure "github.com/LucaToma05/WASAText/service/structure"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) sendMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

	chatId, flag := strconv.Atoi(ps.ByName("chat_id"))
	if flag != nil {
		http.Error(w, "Conversion Error", http.StatusInternalServerError)
		return
	}
	if chatId == 0 {
		http.Error(w, "Parameters extraction error", http.StatusInternalServerError)
		return
	}

	var obj structure.Message
	err = json.NewDecoder(r.Body).Decode(&obj)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	obj.State = 0
	obj.Instant = time.Now()
	obj.ChatId = chatId
	obj.IdSender = userId

	//Check text length
	if len(obj.Text) < 1 {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	//Check text pattern
	match, _ := regexp.MatchString("^.*?$", obj.Text)
	if !match {
		http.Error(w, "Invalid text pattern", http.StatusBadRequest)
		return
	}

	err = rt.db.SendMessage(obj)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(obj)
}
