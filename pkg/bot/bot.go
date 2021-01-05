package bot

import (
	"log"
	"os"
	"time"

	"github.com/beldmian/TaigaBot/pkg/db"
	"github.com/beldmian/TaigaBot/pkg/types"
	"github.com/bwmarrin/discordgo"
	"github.com/top-gg/go-dbl"
	"go.uber.org/zap"
)

// Bot provide struct for bot
type Bot struct {
	Session   *discordgo.Session
	DB        *db.DB
	LogsID    string
	Logger    *zap.Logger
	Commands  []Command
	DBLclient *dbl.Client
}

// InitBot initializes bot process
func InitBot(config types.Config) *Bot {
	var token string
	var dblToken string
	var logsID string
	var datebase db.DB
	if config.Production {
		dbURI, exists := os.LookupEnv("DB_URI")
		token, exists = os.LookupEnv("TOKEN")
		logsID, exists = os.LookupEnv("LOGS_ID")
		dblToken, exists = os.LookupEnv("DBL_TOKEN")
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
		dblToken = config.Bot.DBLToken
	}
	discord, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal(err)
	}
	loggerConfig := zap.NewProductionConfig()
	loggerConfig.OutputPaths = []string{"stdout", "./bot.log"}
	logger, _ := loggerConfig.Build()

	dblClient, err := dbl.NewClient(dblToken)
	if err != nil {
		log.Fatal(err)
	}

	bot := Bot{
		Session:   discord,
		DB:        &datebase,
		LogsID:    logsID,
		Logger:    logger,
		DBLclient: dblClient,
	}

	bot.initCommands()

	bot.Session.AddHandler(bot.OnMessage)
	bot.Session.AddHandler(bot.OnBan)

	bot.Session.State = discordgo.NewState()
	bot.Session.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsAll)

	logger.Info("Bot started")
	if err := bot.Session.Open(); err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			bot.Logger.Info("Data sended")
			if err := bot.PostData(); err != nil {
				bot.Logger.Warn("Data send error", zap.Error(err))
			}
			time.Sleep(60 * time.Second)
		}
	}()

	return &bot
}

// StopBot stops the bot session
func (bot *Bot) StopBot() {
	if err := bot.Session.Close(); err != nil {
		bot.Logger.Fatal("Error on closeing session", zap.Error(err))
	}
}

// PostData ...
func (bot *Bot) PostData() error {
	err := bot.DBLclient.PostBotStats(bot.Session.State.User.ID, &dbl.BotStatsPayload{
		Shards: []int{len(bot.Session.State.Guilds)},
	})
	return err
}
