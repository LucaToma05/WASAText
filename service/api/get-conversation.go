package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) getConversation(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

	chatId, err := strconv.Atoi(ps.ByName("chat_id"))
	if err != nil {
		http.Error(w, "Conversion Error", http.StatusInternalServerError)
		return
	}
	if chatId == 0 {
		http.Error(w, "Parameters extraction error", http.StatusInternalServerError)
		return
	}

	ok, err := rt.db.CheckConversation(chatId)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if !ok {
		http.Error(w, "Resource not found", http.StatusNotFound)
		return
	}

	conversation, err := rt.db.GetConversation(userId, chatId)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(conversation)
}
