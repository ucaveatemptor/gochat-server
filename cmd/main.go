package main

import (
	"context"
	"gochat-server/internal/auth"
	"gochat-server/internal/chat"
	"gochat-server/internal/db"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
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
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/ws", func(w http.ResponseWriter, r *http.Request) {
		chat.ServeWs(hub, w, r)
	})
	r.Post("/auth/register", func(w http.ResponseWriter, r *http.Request) {
		auth.Register(s, w, r)
	})
	r.Post("/auth/login", func(w http.ResponseWriter, r *http.Request) {
		auth.Login(s, w, r)
	})
	r.Get("/users/search", func(w http.ResponseWriter, r *http.Request) {
		chat.FindUsers(s, w, r)
	})
	r.Post("/chat/get", func(w http.ResponseWriter, r *http.Request) {
		chat.GetDMChatID(s, w, r)
	})
	r.Get("/chat/messages/get", func(w http.ResponseWriter, r *http.Request) {
		chat.GetMessages(s, w, r)
	})
	r.Get("/chats/get", func(w http.ResponseWriter, r *http.Request) {
		chat.GetChats(s, w, r)
	})
	log.Print("Server is running")
	http.ListenAndServe(":8080", r)
}
