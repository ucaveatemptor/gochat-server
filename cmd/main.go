package main

import (
	"context"
	"gochat-server/internal/auth"
	"gochat-server/internal/chat"
	"gochat-server/internal/db"
	"log"
	"net/http"
)

func main() {
	ctx := context.TODO()
	s, err := db.NewStorage(ctx)
	if err != nil {
		log.Print("DB Connection err")
		log.Fatal()
	}
	s.CreateIndexes(ctx)
	hub := chat.NewHub(s)

	go hub.Run()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		chat.ServeWs(hub, w, r)
	})
	http.HandleFunc("POST /auth/register", func(w http.ResponseWriter, r *http.Request) {
		auth.Register(s, w, r)
	})
	http.HandleFunc("POST /auth/login", func(w http.ResponseWriter, r *http.Request) {
		auth.Login(s, w, r)
	})
	http.HandleFunc("GET /users/search", func(w http.ResponseWriter, r *http.Request) {
		chat.FindUsers(s, w, r)
	})
	http.HandleFunc("POST /chat/get", func(w http.ResponseWriter, r *http.Request) {
		chat.GetDMChatID(s, w, r)
	})
	log.Print("Server is running")
	http.ListenAndServe(":8080", nil)
}
