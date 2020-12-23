package types

import (
	"time"
)

// Task provide struct for db type task
type Task struct {
	Date   time.Time `json:"date,omitempty"`
	Title  string    `json:"title,omitempty"`
	Done   bool      `json:"done,omitempty"`
	UserID string    `json:"user_id,omitempty" bson:"user_id"`
}

// Config provide struct for config of bot
type Config struct {
	Production bool `toml:"prod"`
	Bot        BotConfig
}

// BotConfig provide struct for bot config
type BotConfig struct {
	Token  string
	DBURI  string `toml:"db_uri"`
	LogsID string `toml:"logs_id"`
}
