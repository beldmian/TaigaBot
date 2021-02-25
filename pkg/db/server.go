package db

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

// ServerInfo provide db struct to store server info
type ServerInfo struct {
	ServerID string `bson:"server_id,omitempty"`
	Prefix   string `bson:"prefix"`
}

// GetPrefix returns prefix of server by id
func (db DB) GetPrefix(id string) (string, error) {
	prefix, found := db.Cache.Get(id)
	if found {
		return prefix.(ServerInfo).Prefix, nil
	}
	client, err := db.Connect()
	if err != nil {
		return "", err
	}
	filter := bson.M{"server_id": id}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var server ServerInfo
	if err := client.Database("info").Collection("servers").FindOne(ctx, filter).Decode(&server); err != nil {
		return "", err
	}
	if server.Prefix != "" {
		db.Cache.Set(id, server, 5*time.Minute)
		return server.Prefix, nil
	}
	return "!", nil
}

// SetPrefix returns prefix of server by id
func (db DB) SetPrefix(id string, prefix string) error {
	client, err := db.Connect()
	if err != nil {
		return err
	}
	filter := bson.M{"server_id": id}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var server ServerInfo
	if err := client.Database("info").Collection("servers").FindOne(ctx, filter).Decode(&server); err != nil {
		return err
	}
	server.Prefix = prefix
	db.Cache.Set(id, server, 5*time.Minute)

	if _, err := client.Database("info").Collection("servers").ReplaceOne(ctx, filter, server); err != nil {
		return err
	}

	return nil
}
