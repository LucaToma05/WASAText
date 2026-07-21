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

	var req structure.AddToGroupRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	//Check id length
	if len(req.UserId) < 3 || len(req.UserId) > 16 {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	//Check id pattern
	match, _ := regexp.MatchString("^.*?$", req.UserId)
	if !match {
		http.Error(w, "Invalid text pattern", http.StatusBadRequest)
		return
	}

	flag, err := rt.db.IsGroup(chatId)
	if !flag {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	err = rt.db.AddToGroup(req, chatId)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(req)
}
