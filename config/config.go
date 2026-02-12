package config

import (
	"os"

	_ "github.com/joho/godotenv/autoload"
)

type MongoDBConfig struct {
	MongoConn   string
	MongoDBName string
}

func LoadConfig() *MongoDBConfig {
	return &MongoDBConfig{
		MongoConn:   os.Getenv("MONGOCONN"),
		MongoDBName: os.Getenv("MONGODB"),
	}
}
