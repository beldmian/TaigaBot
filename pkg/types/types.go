package types

import "time"

// Task provide struct for db type task
type Task struct {
	Date   time.Time `json:"date,omitempty"`
	Title  string    `json:"title,omitempty"`
	Done   bool      `json:"done,omitempty"`
	UserID string    `json:"user_id,omitempty" bson:"user_id"`
}
