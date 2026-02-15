package models

import "time"

type ChatInfo struct {
	ChatID      string          `json:"chatId"`
	Recipient   ChatRecipient   `json:"recipient"`
	LastMessage ChatLastMessage `json:"lastMessage"`
}

type ChatRecipient struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

type ChatLastMessage struct {
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"createdAt"`
}
