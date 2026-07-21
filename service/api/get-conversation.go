package api

import (
	"encoding/json"
	"net/http"

	"github.com/LucaToma05/WASAText/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) getConversation(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
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

	ok, err := rt.db.CheckConversation(chatId, userId)
	if ok == false && err == nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	conversation, err := rt.db.GetConversation(chatId)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(conversation)
}
