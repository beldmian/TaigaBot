package db

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

// DB provides type for datebase
type DB struct {
	DbURL string
}

// Connect is function for db connection
func (db DB) Connect() (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(db.DbURL))
	if err != nil {
		return nil, err
	}
	return client, nil
}