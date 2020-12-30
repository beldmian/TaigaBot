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
	Session  *discordgo.Session
	DB       *db.DB
	LogsID   string
	Logger   *zap.Logger
	Commands []Command
}

// Command provide struct for commands
type Command struct {
	Name        string
	Description string
	Command     string
	Moderation  bool
	Handler     func(s *discordgo.Session, m *discordgo.MessageCreate)
}

func (bot *Bot) initCommands() {
	commands := []Command{
		{
			Name:        "`!help (moderation)`",
			Description: "Список команд бота",
			Command:     "!help",
			Moderation:  false,
			Handler:     bot.Help,
		},
		{
			Name:        "`!colors`",
			Description: "Список доступниых цветов",
			Command:     "!colors",
			Moderation:  false,
			Handler:     bot.ColorsList,
		},
		{
			Name:        "`!color <номер цвета>`",
			Description: "Выдает вам этот цвет",
			Command:     "!color ",
			Moderation:  false,
			Handler:     bot.PickColor,
		},
		{
			Name:        "`!anime <название>`",
			Description: "Ищет аниме по его названию",
			Command:     "!anime ",
			Moderation:  false,
			Handler:     bot.GetAnime,
		},
		{
			Name:        "`!tasks`",
			Description: "Выдает список заданий",
			Command:     "!tasks",
			Moderation:  false,
			Handler:     bot.Tasks,
		},
		{
			Name:        "`!task add <дата 01.02.2020> <текст задания>`",
			Description: "Добавляет вам задание",
			Command:     "!task add ",
			Moderation:  false,
			Handler:     bot.TaskAdd,
		},
		{
			Name:        "`!task done <дата 01.02.2020>`",
			Description: "Отмечает сделанными все задания на данную дату",
			Command:     "!task done ",
			Moderation:  false,
			Handler:     bot.TaskDone,
		},
		{
			Name:        "`!delete <число сообщений>`",
			Description: "Удаляет сообщения",
			Command:     "!delete ",
			Moderation:  true,
			Handler:     bot.BulkDelete,
		},
		{
			Name:        "`!massrole @<роль>`",
			Description: "Выдает или забирает роль у всех на сервере",
			Command:     "!massrole ",
			Moderation:  true,
			Handler:     bot.MassRole,
		},
	}

	bot.Commands = commands
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

	bot.initCommands()

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
