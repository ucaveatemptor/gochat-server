package models

type Message struct {
	ChatID  string `json:"chatId"`
	From    string `json:"from"`
	Content string `json:"content"`
}
