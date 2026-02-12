package models

import "go.mongodb.org/mongo-driver/v2/bson"

// response from mongodb for handlers
type UserResponse struct {
	Id           bson.ObjectID `bson:"_id"`
	Username     string        `bson:"username"`
	PasswordHash string        `bson:"password"`
}

// response for client
type UserResponseJSON struct {
	Id       string `json:"id"`
	Username string `json:"username"`
}

// request user (register + login)
type ReqUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (u *UserResponse) ToJSON() *UserResponseJSON {
	return &UserResponseJSON{
		Id:       u.Id.Hex(),
		Username: u.Username,
	}
}
