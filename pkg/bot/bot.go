package bot

import (
	"log"
	"os"

	"github.com/beldmian/TaigaBot/pkg/db"
	"github.com/beldmian/TaigaBot/pkg/types"
	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
)

// Bot provide struct for bot
type Bot struct {
	Session *discordgo.Session
	DB      *db.DB
	LogsID  string
	Logger  *zap.Logger
}

// InitBot initializes bot process
func InitBot(config types.Config) *Bot {
	var token string
	var logsID string
	var datebase db.DB
	if config.Production {
		dbURI, exists := os.LookupEnv("DB_URI")
		token, exists = os.LookupEnv("TOKEN")
		logsID, exists = os.LookupEnv("LOGS_ID")
		if !exists {
			log.Print("No token or logs channel ID provided")
			return nil
		}
		datebase = db.DB{
			DbURL: dbURI,
		}
	} else {
		token = config.Bot.Token
		logsID = config.Bot.LogsID
		datebase = db.DB{
			DbURL: config.Bot.DBURI,
		}
	}
	discord, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal(err)
	}
	loggerConfig := zap.NewProductionConfig()
	loggerConfig.OutputPaths = []string{"stdout", "/tmp/logs"}
	logger, _ := loggerConfig.Build()

	bot := Bot{
		Session: discord,
		DB:      &datebase,
		LogsID:  logsID,
		Logger:  logger,
	}

	bot.Session.AddHandler(bot.OnMessage)
	bot.Session.AddHandler(bot.OnBan)

	logger.Info("Bot started")
	if err := bot.Session.Open(); err != nil {
		log.Fatal(err)
	}

	return &bot
}

// StopBot stops the bot session
func (bot *Bot) StopBot() {
	if err := bot.Session.Close(); err != nil {
		bot.Logger.Fatal("Error on closeing session", zap.Error(err))
	}
}
