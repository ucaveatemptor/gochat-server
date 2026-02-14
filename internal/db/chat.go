package db

import (
	"context"
	"gochat-server/internal/models"

	"log"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func (s *Storage) GetOrCreateChatAndReturnChatID(ctx context.Context, users [2]string) (string, error) {
	// users = users IDs

	// sort
	if users[0] > users[1] {
		users[0], users[1] = users[1], users[0]
	}
	for _, id := range users {
		_, err := s.GetUserByID(ctx, id)
		if err != nil {
			log.Print("GetOrCreateChat - user doesnt exist")
			return "", err
		}
	}

	filter := bson.M{
		"users": users,
	}

	update := bson.M{
		"$setOnInsert": bson.M{
			"users": users,
		},
	}

	opts := options.FindOneAndUpdate().
		SetUpsert(true).
		SetReturnDocument(options.After)

	var chat struct {
		ID bson.ObjectID `bson:"_id"`
	}

	err := s.Chats.FindOneAndUpdate(
		ctx,
		filter,
		update,
		opts,
	).Decode(&chat)

	if err != nil {
		return "", err
	}

	return chat.ID.Hex(), nil
}

func (s *Storage) GetMessagesByChatID(ctx context.Context, chatId string) []*models.Message {
	filter := bson.D{
		{Key: "chatId", Value: chatId},
	}
	cursor, err := s.Messages.Find(ctx, filter)
	if err != nil {
		log.Print("GetMessagesByChatID err message find")
	}
	defer cursor.Close(ctx)
	messages := make([]*models.Message, 0)
	for cursor.Next(ctx) {
		var m models.Message
		if err := cursor.Decode(&m); err != nil {
			return nil
		}
		messages = append(messages, &m)
	}
	return messages
}
func (s *Storage) SaveMessage(ctx context.Context, msg models.Message) {
	doc := bson.M{
		"chatId":    msg.ChatID,
		"from":      msg.From,
		"content":   msg.Content,
		"createdAt": time.Now(),
	}

	_, err := s.Messages.InsertOne(ctx, doc)
	if err != nil {
		log.Printf("SaveMessage error %v", err)
		return
	}
}
