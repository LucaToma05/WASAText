package api

import (
	"net/http"

	"github.com/LucaToma05/WASAText/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) uncommentMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
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

	commentId := ps.ByName("commentId")
	if messageId == "" {
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

	ok, err = rt.db.CheckMessageInChat(chatId, messageId)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if !ok {
		http.Error(w, "Message not found in chat", http.StatusNotFound)
		return
	}

	ok, err = rt.db.CheckCommentInMessage(commentId, messageId)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if !ok {
		http.Error(w, "Commment not found in message", http.StatusNotFound)
		return
	}

	err = rt.db.UncommentMessage(messageId, commentId, userId)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
