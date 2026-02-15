package db

import (
	"context"
	"fmt"
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

func (s *Storage) GetChatsByUserID(ctx context.Context, userId string) ([]models.ChatInfo, error) {
	filter := bson.M{"users": userId}
	cursor, err := s.Chats.Find(ctx, filter)
	if err != nil {
		log.Printf("db chat.go GetChatsByUserID %v", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var chats []models.ChatInfo

	for cursor.Next(ctx) {
		var chatDoc struct {
			ID    bson.ObjectID `bson:"_id"`
			Users [2]string     `bson:"users"`
		}
		if err := cursor.Decode(&chatDoc); err != nil {
			log.Printf("1 db chat.go GetChatsByUserID %v", err)
			continue
		}

		chatIdStr := chatDoc.ID.Hex()

		recipient, err := s.GetChatRecipient(ctx, chatIdStr, userId)
		if err != nil {
			log.Printf("2 db chat.go GetChatsByUserID %v", err)
			continue
		}

		lastMsg := s.GetLastMessageByChatID(ctx, chatIdStr)

		chats = append(chats, models.ChatInfo{
			ChatID:      chatIdStr,
			Recipient:   *recipient,
			LastMessage: *lastMsg,
		})
	}

	return chats, nil
}
func (s *Storage) GetChatRecipient(ctx context.Context, chatId string, myUserId string) (*models.ChatRecipient, error) {
	objID, err := bson.ObjectIDFromHex(chatId)
	if err != nil {
		return nil, fmt.Errorf("invalid chat id: %v", err)
	}

	var chat struct {
		Users []string `bson:"users"`
	}

	err = s.Chats.FindOne(ctx, bson.M{"_id": objID}).Decode(&chat)
	if err != nil {
		return nil, fmt.Errorf("chat not found: %v", err)
	}

	var recipientID string
	for _, uID := range chat.Users {
		if uID != myUserId {
			recipientID = uID
			break
		}
	}

	if recipientID == "" {
		return nil, fmt.Errorf("recipient not found in chat")
	}

	var user struct {
		Username string `bson:"username"`
	}

	objrID, err := bson.ObjectIDFromHex(recipientID)
	if err != nil {
		return nil, fmt.Errorf("invalid recipient ID format: %v", err)
	}
	err = s.Users.FindOne(ctx, bson.M{"_id": objrID}).Decode(&user)
	if err != nil {
		return nil, fmt.Errorf("user not found: %v", err)
	}

	return &models.ChatRecipient{
		ID:       recipientID,
		Username: user.Username,
	}, nil
}
func (s *Storage) GetLastMessageByChatID(ctx context.Context, chatId string) *models.ChatLastMessage {
	filter := bson.D{
		{Key: "chatId", Value: chatId},
	}
	opts := options.FindOne().SetSort(bson.D{{Key: "createdAt", Value: -1}})
	var msg models.Message
	err := s.Messages.FindOne(ctx, filter, opts).Decode(&msg)
	if err != nil {
		log.Printf("%v", err)
	}
	var lastmsg models.ChatLastMessage

	lastmsg.Content = msg.Content
	lastmsg.CreatedAt = msg.CreatedAt

	return &lastmsg

}
func (s *Storage) GetMessagesByChatID(ctx context.Context, chatId string) []*models.Message {
	filter := bson.D{
		{Key: "chatId", Value: chatId},
	}
	opts := options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}})
	cursor, err := s.Messages.Find(ctx, filter, opts)
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
