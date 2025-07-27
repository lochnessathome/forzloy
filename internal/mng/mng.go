package mng

import (
	"fmt"
	"os"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Database struct {
	*mongo.Database
}

func New() (*Database, error) {
	dbUrl := os.Getenv("MONGO_DATABASE_URL")
	if dbUrl == "" {
		return nil, fmt.Errorf("Empty MONGO_DATABASE_URL environment variable")
	}

	dbName := os.Getenv("MONGO_DATABASE_NAME")
	if dbName == "" {
		return nil, fmt.Errorf("Empty MONGO_DATABASE_NAME environment variable")
	}

	client, err := mongo.Connect(options.Client().ApplyURI(dbUrl))
	if err != nil {
		return nil, err
	}

	db := client.Database(dbName)

	return &Database{db}, nil
}
