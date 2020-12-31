package db

import (
	"context"
	"time"

	"github.com/beldmian/TaigaBot/pkg/types"
	"go.mongodb.org/mongo-driver/bson"
)

// GetTasks provide function for get tasks
func (db DB) GetTasks(userID string) ([]types.Task, error) {
	var tasks []types.Task
	client, err := db.Connect()
	if err != nil {
		return tasks, err
	}
	filter := bson.M{"user_id": userID}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cursor, err := client.Database("tasker").Collection("tasks").Find(ctx, filter)
	if err != nil {
		return tasks, err
	}
	for cursor.Next(ctx) {
		var task types.Task
		if err := cursor.Decode(&task); err != nil {
			return tasks, err
		}
		if task.Done && time.Now().Sub(task.Date) > 720*time.Hour {
			if _, err := client.Database("tasker").Collection("tasks").DeleteOne(ctx, bson.M{"title": task.Title}); err != nil {
				return tasks, err
			}
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

// AddTask provide function for get tasks
func (db DB) AddTask(task types.Task) error {
	client, err := db.Connect()
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if _, err := client.Database("tasker").Collection("tasks").InsertOne(ctx, task); err != nil {
		return err
	}
	return nil
}

// DoneTask provide function for make task done
func (db DB) DoneTask(date time.Time) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := db.Connect()
	if err != nil {
		return err
	}
	if _, err := client.Database("tasker").Collection("tasks").UpdateMany(ctx, bson.M{"date": date}, bson.M{"$set": bson.M{"done": true}}); err != nil {
		return err
	}
	return nil
}
