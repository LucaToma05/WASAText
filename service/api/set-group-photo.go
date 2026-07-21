package api

import (
	"encoding/json"
	"net/http"
	"regexp"

	"github.com/LucaToma05/WASAText/service/api/reqcontext"
	structure "github.com/LucaToma05/WASAText/service/structure"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) setGroupPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	var req structure.GroupPhotoRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if len(req.Url) < 3 || len(req.Url) > 1000 {
		http.Error(w, "Invalid name group length", http.StatusBadRequest)
		return
	}

	match, _ := regexp.MatchString("^.*?$", req.Url)
	if !match {
		http.Error(w, "Invalid name pattern", http.StatusBadRequest)
		return
	}

	userId := ctx.UserId
	if userId == "" {
		http.Error(w, "Unauthorized User", http.StatusUnauthorized)
		return
	}

	chatId := ps.ByName("chatId")
	if chatId == "" {
		http.Error(w, "Parameters extraction error", http.StatusInternalServerError)
		return
	}

	flag, err := rt.db.IsGroup(chatId)
	if !flag {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	err = rt.db.SetGroupPhoto(req, chatId)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(req)
}
