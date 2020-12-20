package bot

import (
	"log"
	"os"

	"github.com/beldmian/TaigaBot/pkg/db"
	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
)

var logsID string
var discord *discordgo.Session
var datebase db.DB
var logger *zap.Logger

// InitBot initializes bot process
func InitBot() {
	token, exists := os.LookupEnv("TOKEN")
	logsID, exists = os.LookupEnv("LOGS_ID")
	dbURI, exists := os.LookupEnv("DB_URI")
	if !exists {
		log.Print("No token or logs channel ID provided")
		return
	}
	datebase = db.DB{
		DbURL: dbURI,
	}
	discord, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal(err)
	}
	loggerConfig := zap.NewProductionConfig()
	loggerConfig.OutputPaths = []string{"stdout", "/tmp/logs"}
	logger, _ = loggerConfig.Build()
	discord.AddHandler(OnMessage)
	discord.AddHandler(OnBan)

	logger.Info("Bot started")
	if err := discord.Open(); err != nil {
		log.Fatal(err)
	}
}

// StopBot stops the bot session
func StopBot() {
	if err := discord.Close(); err != nil {
		logger.Fatal("Error on closeing session", zap.Error(err))
	}
}
