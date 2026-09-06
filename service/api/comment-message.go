package api

import (
	"encoding/json"
	"net/http"

	"github.com/LucaToma05/WASAText/service/api/reqcontext"
	structure "github.com/LucaToma05/WASAText/service/structure"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) commentMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
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

	var req structure.CommentRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	commentId, err := generateId()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	senderInfo, errInfo := rt.db.GetMyInformation(userId)
	senderName := userId
	if errInfo == nil && senderInfo.Name != "" {
		senderName = senderInfo.Name
	}

	newComment := structure.Comment{
		ID:             commentId,
		IdSender:       userId,
		SenderUsername: senderName,
		SenderName:     senderName,
		Emoji:          req.Emoji,
		MessageId:      messageId,
	}

	err = rt.db.CommentMessage(newComment)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err = json.NewEncoder(w).Encode(newComment); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
