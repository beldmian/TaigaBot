package db

import (
	"context"
	"time"

	"github.com/patrickmn/go-cache"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DB provides type for datebase
type DB struct {
	DbURL string
	Cache *cache.Cache
}

// New returns a new DB object
func New(url string) DB {
	return DB{
		DbURL: url,
		Cache: cache.New(5*time.Minute, 10*time.Minute),
	}
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
