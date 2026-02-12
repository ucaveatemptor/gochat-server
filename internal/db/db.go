package db

import (
	"context"
	"fmt"
	"gochat-server/config"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Storage struct {
	Client   *mongo.Client
	Users    *mongo.Collection
	Messages *mongo.Collection
	Chats    *mongo.Collection
}

func NewStorage(ctx context.Context) (*Storage, error) {
	conf := config.LoadConfig()
	uri := conf.MongoConn
	dbName := conf.MongoDBName

	clientOptions := options.Client().
		ApplyURI(uri).
		SetMaxPoolSize(50).
		SetConnectTimeout(5 * time.Second)

	client, err := mongo.Connect(clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}

	// Connection check
	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("ping failed (check if docker is running): %w", err)
	}

	db := client.Database(dbName)

	return &Storage{
		Client:   client,
		Users:    db.Collection("users"),
		Messages: db.Collection("messages"),
		Chats:    db.Collection("chats"),
	}, nil
}
func (s *Storage) CreateIndexes(ctx context.Context) error {
	index := mongo.IndexModel{
		Keys: bson.D{
			{Key: "username", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	}

	_, err := s.Users.Indexes().CreateOne(ctx, index)
	return err
}
