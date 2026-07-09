package api

import (
	"net/http"
	"strconv"

	"github.com/LucaToma05/WASAText/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) deleteMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
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

	messageId, flag := strconv.Atoi(ps.ByName("message_id"))
	if flag != nil {
		http.Error(w, "Conversion Error", http.StatusInternalServerError)
		return
	}
	if messageId == 0 {
		http.Error(w, "Parameters extraction error", http.StatusInternalServerError)
		return
	}

	err = rt.db.DeleteMessage(userId, chatId, messageId)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
