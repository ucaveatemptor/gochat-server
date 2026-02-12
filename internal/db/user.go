package db

import (
	"context"
	"gochat-server/internal/models"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func (s *Storage) AddUser(ctx context.Context, username string, password string) error {
	user := bson.D{
		{Key: "username", Value: username},
		{Key: "password", Value: password},
	}
	_, err := s.Users.InsertOne(ctx, user)
	return err
}
func (s *Storage) GetUserByName(ctx context.Context, username string) (*models.UserResponse, error) {
	filter := bson.D{
		{Key: "username", Value: username},
	}
	var user models.UserResponse
	err := s.Users.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
func (s *Storage) GetUserByID(ctx context.Context, id string) (*models.UserResponse, error) {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	filter := bson.D{
		{Key: "_id", Value: objID},
	}
	var user models.UserResponse
	err = s.Users.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *Storage) FindUsersByPrefix(ctx context.Context, prefix string) ([]*models.UserResponseJSON, error) {
	filter := bson.D{
		{Key: "username", Value: bson.D{
			{Key: "$regex", Value: "^" + prefix}, // starts with
			{Key: "$options", Value: "i"},        // case-insensitive
		}},
	}

	cursor, err := s.Users.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []*models.UserResponseJSON
	for cursor.Next(ctx) {
		var u models.UserResponse
		if err := cursor.Decode(&u); err != nil {
			return nil, err
		}
		users = append(users, u.ToJSON())
	}

	return users, nil
}
