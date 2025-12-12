package api

import (
	"encoding/json"
	"net/http"
	"regexp"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) setMyUserName(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var req UsernameRequest

	//Read body content
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	//Check username lenght
	if len(req.Username) < 3 || len(req.Username) > 25 {
		http.Error(w, "Invalid username length", http.StatusBadRequest)
		return
	}

	match, _ := regexp.MatchString("^.*?$", req.Username)
	if !match {
		http.Error(w, "Invalid username pattern", http.StatusBadRequest)
		return
	}

	userId := extractBearer(r.Header.Get("Authorization"))

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(body)

}
