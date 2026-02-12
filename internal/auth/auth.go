package auth

import (
	"encoding/json"
	"gochat-server/internal/db"
	"gochat-server/internal/models"
	"log"
	"net/http"

	"github.com/tailscale/golang-x-crypto/bcrypt"
)

func Register(s *db.Storage, w http.ResponseWriter, r *http.Request) {
	var nu models.ReqUser
	if err := json.NewDecoder(r.Body).Decode(&nu); err != nil {
		http.Error(w, "Wrong format JSON (Register)", http.StatusBadRequest)
		return
	}
	nu.Password, _ = HashPassword(nu.Password)
	if err := s.AddUser(r.Context(), nu.Username, nu.Password); err != nil {
		http.Error(w, "Error", http.StatusNotFound)
		log.Print("Add user error")
	}
}
func Login(s *db.Storage, w http.ResponseWriter, r *http.Request) {
	var ru models.ReqUser
	if err := json.NewDecoder(r.Body).Decode(&ru); err != nil {
		http.Error(w, "Wrong format JSON (Login)", http.StatusBadRequest)
		return
	}
	u, err := s.GetUserByName(r.Context(), ru.Username)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		log.Print("Get user error")
		return
	}

	// check if hash & password match
	if err = bcrypt.CompareHashAndPassword(
		[]byte(u.PasswordHash),
		[]byte(ru.Password),
	); err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	resp := u.ToJSON()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	return string(hash), err
}
func CheckPassword(password string, hash string) error {
	return bcrypt.CompareHashAndPassword(
		[]byte(hash),
		[]byte(password),
	)
}
