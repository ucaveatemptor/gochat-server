package models

import "time"

type Message struct {
	ChatID    string    `bson:"chatId" json:"chatId"`
	From      string    `bson:"from" json:"from"`
	Content   string    `bson:"content" json:"content"`
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
}
