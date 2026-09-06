package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/LucaToma05/WASAText/service/api/reqcontext"
	structure "github.com/LucaToma05/WASAText/service/structure"
)

type CreateChatPayload struct {
	Type           *bool           `json:"type"`
	IsGroup        *bool           `json:"isGroup"`
	Name           string          `json:"name"`
	GroupName      string          `json:"groupName"`
	Photo          string          `json:"photo"`
	Participants   json.RawMessage `json:"participants"`
	Users          []string        `json:"users"`
	TargetUsername string          `json:"targetUsername"`
}

type CreateChatResponse struct {
	ID           string               `json:"id"`
	ChatId       string               `json:"chatId"`
	Name         string               `json:"name"`
	Photo        string               `json:"photo"`
	IsGroup      bool                 `json:"isGroup"`
	Type         bool                 `json:"type"`
	Participants []structure.Username `json:"participants"`
}

func (rt *_router) createChat(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	userId := ctx.UserId
	if userId == "" {
		http.Error(w, "Unauthorized User", http.StatusUnauthorized)
		return
	}

	var req CreateChatPayload
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	// Extract participants list if provided
	var rawParticipants []string
	if len(req.Participants) > 0 {
		var strList []string
		if err := json.Unmarshal(req.Participants, &strList); err == nil {
			rawParticipants = append(rawParticipants, strList...)
		} else {
			var userList []structure.Username
			if err := json.Unmarshal(req.Participants, &userList); err == nil {
				for _, u := range userList {
					if u.Username != "" {
						rawParticipants = append(rawParticipants, u.Username)
					}
				}
			}
		}
	}
	if len(req.Users) > 0 {
		rawParticipants = append(rawParticipants, req.Users...)
	}

	// Distinguish between single chat and group chat
	isGroup := false
	if req.Type != nil && *req.Type {
		isGroup = true
	} else if req.IsGroup != nil && *req.IsGroup {
		isGroup = true
	} else if req.GroupName != "" {
		isGroup = true
	} else if len(rawParticipants) > 1 {
		isGroup = true
	}

	if !isGroup {
		// --- LOGICA CHAT SINGOLA 1-A-1 ---
		targetUser := req.Name
		if targetUser == "" && req.TargetUsername != "" {
			targetUser = req.TargetUsername
		}
		if targetUser == "" && len(rawParticipants) == 1 {
			targetUser = rawParticipants[0]
		}

		if len(targetUser) < 3 || len(targetUser) > 25 {
			http.Error(w, "Invalid name length", http.StatusBadRequest)
			return
		}

		id, err := rt.db.GetIdUserByUsername(targetUser)
		if err != nil {
			http.Error(w, "Internal 1 Server Error", http.StatusInternalServerError)
			return
		}
		if id == "" {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		exists, err := rt.db.CheckChatExistsBetween(userId, id)
		if err != nil {
			http.Error(w, "Internal 2 Server Error", http.StatusInternalServerError)
			return
		}
		if exists {
			http.Error(w, "Chat already exists", http.StatusConflict)
			return
		}

		chatId, err := generateId()
		if err != nil {
			http.Error(w, "Internal 3 Server Error", http.StatusInternalServerError)
			return
		}

		participants := []structure.UserId{
			{UserId: userId},
			{UserId: id},
		}

		infoChat := structure.InfoChat{
			IsGroup: false,
			Name:    targetUser,
			Photo:   "",
		}

		chatErr := rt.db.CreateChat(infoChat, chatId, participants)
		if chatErr != nil {
			http.Error(w, "Internal 4 Server Error", http.StatusInternalServerError)
			return
		}

		response := CreateChatResponse{
			ID:      chatId,
			ChatId:  chatId,
			Name:    targetUser,
			Photo:   "",
			IsGroup: false,
			Type:    false,
			Participants: []structure.Username{
				{Username: targetUser},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err = json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "Internal 5 Server Error", http.StatusInternalServerError)
			return
		}

	} else {
		// --- LOGICA CHAT DI GRUPPO ---
		groupName := req.Name
		if groupName == "" && req.GroupName != "" {
			groupName = req.GroupName
		}

		if len(groupName) < 3 || len(groupName) > 25 {
			http.Error(w, "Invalid group name length (must be 3-25 characters)", http.StatusBadRequest)
			return
		}

		chatId, err := generateId()
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Build members list, always including the creator
		memberMap := make(map[string]bool)
		memberMap[userId] = true

		for _, p := range rawParticipants {
			if p == "" {
				continue
			}
			pId, err := rt.db.GetIdUserByUsername(p)
			if err == nil && pId != "" {
				memberMap[pId] = true
			} else {
				// Check if 'p' is already a valid userId
				uInfo, err2 := rt.db.GetMyInformation(p)
				if err2 == nil && uInfo.UserId != "" {
					memberMap[uInfo.UserId] = true
				} else {
					http.Error(w, fmt.Sprintf("User '%s' not found", p), http.StatusNotFound)
					return
				}
			}
		}

		var participants []structure.UserId
		for mId := range memberMap {
			participants = append(participants, structure.UserId{UserId: mId})
		}

		infoChat := structure.InfoChat{
			IsGroup: true,
			Name:    groupName,
			Photo:   req.Photo,
		}

		chatErr := rt.db.CreateChat(infoChat, chatId, participants)
		if chatErr != nil {
			http.Error(w, "Internal Server Error creating group", http.StatusInternalServerError)
			return
		}

		var respParticipants []structure.Username
		for _, p := range rawParticipants {
			respParticipants = append(respParticipants, structure.Username{Username: p})
		}

		response := CreateChatResponse{
			ID:           chatId,
			ChatId:       chatId,
			Name:         groupName,
			Photo:        req.Photo,
			IsGroup:      true,
			Type:         true,
			Participants: respParticipants,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err = json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "Internal Server Error encoding response", http.StatusInternalServerError)
			return
		}
	}
}
