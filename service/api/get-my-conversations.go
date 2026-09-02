package api

import (
	"encoding/json"
	"net/http"
	"log"

	"github.com/LucaToma05/WASAText/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) getMyConversations(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	/*
	userId := ctx.UserId
	if userId == "" {
		http.Error(w, "Unauthorized User", http.StatusUnauthorized)
		return
	}
	*/
	ctx.Logger.Info("DEBUG: ctx.UserId = " + ctx.UserId)

	chats, err := rt.db.GetMyConversations(ctx.UserId)
	if err != nil {
		log.Println("ERRORE GET:", err)
		http.Error(w, "Internal Server Error 1", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(chats); err != nil {
		http.Error(w, "Internal Server Error 2", http.StatusInternalServerError)
		return
	}
}
