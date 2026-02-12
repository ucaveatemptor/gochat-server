package chat

import (
	"encoding/json"
	"gochat-server/internal/db"
	"log"
	"net/http"
)

func GetDMChatID(s *db.Storage, w http.ResponseWriter, r *http.Request) { // POST
	// for request {
	//	users: [string ID,string ID]
	//}
	ctx := r.Context()

	var ChatRequest struct {
		Users [2]string `json:"users"` // users = userIDs
	}

	if err := json.NewDecoder(r.Body).Decode(&ChatRequest); err != nil {
		http.Error(w, "Wrong format JSON", http.StatusBadRequest)
		log.Print("GetDMChatID - failed to decode")
		return
	}
	for _, u := range ChatRequest.Users {
		if u == "" {
			http.Error(w, "2 users required", http.StatusBadRequest)
			log.Print("GetDMChatID - wrong request")
			return
		}
	}
	id, err := s.GetOrCreateChatAndReturnChatID(ctx, ChatRequest.Users)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		log.Print("GetDMChatID - getchat failed")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(id)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		log.Print("GetDMChatID - failed to encode response")
		return
	}

}
func FindUsers(s *db.Storage, w http.ResponseWriter, r *http.Request) { // GET
	ctx := r.Context()

	prefix := r.URL.Query().Get("username")
	if prefix == "" {
		http.Error(w, "username query param is required", http.StatusBadRequest)
		return
	}

	users, err := s.FindUsersByPrefix(ctx, prefix)
	if err != nil {
		http.Error(w, "failed to find users", http.StatusInternalServerError)
		log.Print("Failed to find users")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(users)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		log.Print("FindUsers - failed to encode response")
		return
	}
}
