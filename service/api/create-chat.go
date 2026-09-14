package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/LucaToma05/WASAText/service/api/reqcontext"
	structure "github.com/LucaToma05/WASAText/service/structure"
)

func (rt *_router) createChat(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	userId := ctx.UserId
	if userId == "" {
		http.Error(w, "Unauthorized User", http.StatusUnauthorized)
		return
	}

	var req structure.InfoChat
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	// Distinguish between single chat and group chat
	isGroup := req.IsGroup || len(req.Participants) > 1

	if !isGroup {
		// --- LOGICA CHAT SINGOLA 1-A-1 ---
		targetUser := req.Name

		if len(targetUser) < 3 || len(targetUser) > 25 {
			http.Error(w, "Invalid name length", http.StatusBadRequest)
			return
		}

		id, err := rt.db.GetIdUserByUsername(targetUser)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		if id == "" {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		exists, err := rt.db.CheckChatExistsBetween(userId, id)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		if exists {
			http.Error(w, "Chat already exists", http.StatusConflict)
			return
		}

		chatId, err := generateId()
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Users participants of the single chat
		singleChatParticipants := []structure.UserId{
			{UserId: userId},
			{UserId: id},
		}

		// New chat
		infoChat := structure.InfoChat{
			IsGroup: false,
			Name:    targetUser,
			Photo:   "",
		}

		chatErr := rt.db.CreateChat(infoChat, chatId, singleChatParticipants)
		if chatErr != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Response
		response := structure.Chat{
			ChatId:  chatId,
			Name:    targetUser,
			Photo:   "",
			IsGroup: false,
			Participants: []structure.Username{
				{Username: targetUser},
			},
			Messages: []structure.Message{},
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err = json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

	} else {
		// --- LOGICA CHAT DI GRUPPO ---
		groupName := req.Name

		if len(groupName) < 3 || len(groupName) > 25 {
			http.Error(w, "Invalid group name length (must be 3-25 characters)", http.StatusBadRequest)
			return
		}

		chatId, err := generateId()
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Initialize the members including the group creator
		memberMap := map[string]bool{userId: true}
		chatMembers := []structure.UserId{
			{UserId: userId},
		}

		for _, p := range req.Participants {
			if p.Username == "" {
				continue
			}

			// Get user id by username
			pId, err := rt.db.GetIdUserByUsername(p.Username)
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			// If the user doesn't exist
			if pId == "" {
				http.Error(w, fmt.Sprintf("User '%s' not found", p.Username), http.StatusNotFound)
				return
			}
			// If the user is not already in the group, add it
			if !memberMap[pId] {
				memberMap[pId] = true
				chatMembers = append(chatMembers, structure.UserId{UserId: pId})
			}
		}

		// Group info
		infoChat := structure.InfoChat{
			IsGroup: true,
			Name:    groupName,
			Photo:   req.Photo,
		}

		chatErr := rt.db.CreateChat(infoChat, chatId, chatMembers)
		if chatErr != nil {
			http.Error(w, "Internal Server Error creating group", http.StatusInternalServerError)
			return
		}

		// Response
		response := structure.Chat{
			ChatId:       chatId,
			Name:         groupName,
			Photo:        req.Photo,
			IsGroup:      true,
			Participants: req.Participants,
			Messages:     []structure.Message{},
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err = json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "Internal Server Error encoding response", http.StatusInternalServerError)
			return
		}
	}
}
