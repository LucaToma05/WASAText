package api

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strconv"

	"github.com/LucaToma05/WASAText/service/api/reqcontext"
	structure "github.com/LucaToma05/WASAText/service/structure"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) setGroupName(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	var req structure.GroupNameRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if len(req.Name) < 3 || len(req.Name) > 50 {
		http.Error(w, "Invalid name group length", http.StatusBadRequest)
		return
	}

	match, _ := regexp.MatchString("^.*?$", req.Name)
	if !match {
		http.Error(w, "Invalid name pattern", http.StatusBadRequest)
		return
	}

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Unauthorized User", http.StatusUnauthorized)
		return
	}
	userId, err := extractUserId(authHeader)
	if err != nil {
		http.Error(w, "Token extraction error", http.StatusInternalServerError)
	}

	groupId, flag := strconv.Atoi(ps.ByName("group_id"))
	if flag != nil {
		http.Error(w, "Conversion Error", http.StatusInternalServerError)
		return
	}
	if groupId == 0 {
		http.Error(w, "Parameters extraction error", http.StatusInternalServerError)
		return
	}

	err = rt.db.SetGroupName(req, userId, groupId)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	//Ritorna il gruppo
	//!!!
}
