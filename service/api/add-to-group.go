package api

import (
	"encoding/json"
	"net/http"
	"regexp"

	"github.com/LucaToma05/WASAText/service/api/reqcontext"
	structure "github.com/LucaToma05/WASAText/service/structure"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) addToGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
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

	var req structure.UserId
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Check id length (allow 3 to 25 to support both username and userId)
	if len(req.UserId) < 3 || len(req.UserId) > 25 {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Check id pattern
	match, _ := regexp.MatchString("^.*?$", req.UserId)
	if !match {
		http.Error(w, "Invalid text pattern", http.StatusBadRequest)
		return
	}

	// Check if the chat is a group chat
	flag, err := rt.db.IsGroup(chatId)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if !flag {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Check if the user is in the group
	ok, err := rt.db.CheckConversation(chatId, userId)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if !ok {
		http.Error(w, "User not in the group", http.StatusForbidden)
		return
	}

	// Resolve target user: req.UserId can be either a username (e.g. "bob") or a userId UUID
	var targetUserId string
	uId, err := rt.db.GetIdUserByUsername(req.UserId)
	if err == nil && uId != "" {
		targetUserId = uId
	} else {
		uInfo, err2 := rt.db.GetMyInformation(req.UserId)
		if err2 != nil || uInfo.UserId == "" {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		targetUserId = uInfo.UserId
	}

	// Check if target user is already in the group
	alreadyMember, err := rt.db.CheckConversation(chatId, targetUserId)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if alreadyMember {
		http.Error(w, "User already in group", http.StatusBadRequest)
		return
	}

	req.UserId = targetUserId
	err = rt.db.AddToGroup(req, chatId)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err = json.NewEncoder(w).Encode(req); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
