package api

import (
	"net/http"

	"github.com/LucaToma05/WASAText/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) deleteMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
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

	messageId := ps.ByName("messageId")
	if messageId == "" {
		http.Error(w, "Parameters extraction error", http.StatusBadRequest)
		return
	}

	err := rt.db.DeleteMessage(chatId, messageId)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
