package bot

import "github.com/bwmarrin/discordgo"

const (
	basicCategory = iota
	moderationCategory
	gifCategory
)

// Command provide struct for commands
type Command struct {
	Translation Translation
	Command     string
	Category    int
	Permissions int
	Handler     func(s *discordgo.Session, m *discordgo.MessageCreate, locale string)
}

// Translation provide struct for commands translation
type Translation struct {
	RussianName        string
	EnglishName        string
	RussianDescription string
	EnglishDescription string
}

func (bot *Bot) initCommands() {
	commands := []Command{
		{
			Translation: Translation{
				RussianName:        "`!help (moderation) (gif)`",
				RussianDescription: "Список команд бота",
				EnglishName:        "`!help (moderation) (gif)`",
				EnglishDescription: "List of bot commands",
			},
			Command:  "!help",
			Category: basicCategory,
			Handler:  bot.Help,
		},
		{
			Translation: Translation{
				RussianName:        "`!colors`",
				RussianDescription: "Список доступниых цветов",
				EnglishName:        "`!colors`",
				EnglishDescription: "List of avaivable colors",
			},
			Command:  "!colors",
			Category: basicCategory,
			Handler:  bot.ColorsList,
		},
		{
			Translation: Translation{
				RussianName:        "`!color <номер цвета>`",
				RussianDescription: "Выдает вам этот цвет",
				EnglishName:        "`!color <color number>`",
				EnglishDescription: "Give you that color",
			},
			Command:  "!color ",
			Category: basicCategory,
			Handler:  bot.PickColor,
		},
		{
			Translation: Translation{
				RussianName:        "`!anime <название>`",
				RussianDescription: "Ищет аниме по его названию",
				EnglishName:        "`!anime <name>`",
				EnglishDescription: "Search anime by name (russian description only)",
			},
			Command:  "!anime ",
			Category: basicCategory,
			Handler:  bot.GetAnime,
		},
		{
			Translation: Translation{
				RussianName:        "`!tasks`",
				RussianDescription: "Список ваших заданий",
				EnglishName:        "`!tasks`",
				EnglishDescription: "Your task list",
			},
			Command:  "!tasks",
			Category: basicCategory,
			Handler:  bot.Tasks,
		},
		{
			Translation: Translation{
				RussianName:        "`!task add <дата 01.02.2020> <текст задания>`",
				RussianDescription: "Добавляет вам задаие",
				EnglishName:        "`!task add <date 01.02.2020> <task desc>`",
				EnglishDescription: "Add task to your list",
			},
			Command:  "!task add ",
			Category: basicCategory,
			Handler:  bot.TaskAdd,
		},
		{
			Translation: Translation{
				RussianName:        "`!task done <дата 01.02.2020>`",
				RussianDescription: "Отмечает задание сделанным",
				EnglishName:        "`!task done <date 01.02.2020>`",
				EnglishDescription: "Mark task as done",
			},
			Command:  "!task done ",
			Category: basicCategory,
			Handler:  bot.TaskDone,
		},
		{
			Translation: Translation{
				RussianName:        "`!delete <число сообщений>`",
				RussianDescription: "Удаляет сообщения",
				EnglishName:        "`!delete <message count>`",
				EnglishDescription: "Delete messages",
			},
			Command:     "!delete ",
			Category:    moderationCategory,
			Handler:     bot.BulkDelete,
			Permissions: 8192,
		},
		{
			Translation: Translation{
				RussianName:        "`!massrole @<роль>`",
				RussianDescription: "Выдает или забирает роль у всех на сервере",
				EnglishName:        "`!massrole @<role>`",
				EnglishDescription: "Give role for all server members",
			},
			Command:     "!massrole ",
			Category:    moderationCategory,
			Handler:     bot.MassRole,
			Permissions: 268435456,
		},
		{
			Translation: Translation{
				RussianName:        "`!poll <вариант 1> | <вариант 2> ...`",
				RussianDescription: "Создает опрос с несколькими вариантами ответа",
				EnglishName:        "`!poll <option 1> | <option 2> ...`",
				EnglishDescription: "Create poll with some answers",
			},
			Command:  "!poll ",
			Category: basicCategory,
			Handler:  bot.Poll,
		},
		{
			Translation: Translation{
				RussianName:        "`!vote`",
				RussianDescription: "Вы можете поддержать бота проголосовав за него на top.gg",
				EnglishName:        "`!vote`",
				EnglishDescription: "You can support bot on top.gg",
			},
			Command:  "!vote",
			Category: basicCategory,
			Handler:  bot.Vote,
		},
		{
			Translation: Translation{
				RussianName:        "`!kiss`",
				RussianDescription: "Отправляет гифку или каринку с поцелуем",
				EnglishName:        "`!kiss`",
				EnglishDescription: "Send gif or image with kiss",
			},
			Command:  "!kiss",
			Category: gifCategory,
			Handler:  bot.Kiss,
		},
		{
			Translation: Translation{
				RussianName:        "`!hug`",
				RussianDescription: "Отправляет гифку обнимашек",
				EnglishName:        "`!hug`",
				EnglishDescription: "Send gif or image with hug",
			},
			Command:  "!hug",
			Category: gifCategory,
			Handler:  bot.Hug,
		},
		{
			Translation: Translation{
				RussianName:        "`!pat`",
				RussianDescription: "Отправляет гифку с поглаживанием по голове",
				EnglishName:        "`!pat`",
				EnglishDescription: "Send gif with pat",
			},
			Command:  "!pat",
			Category: gifCategory,
			Handler:  bot.Pat,
		},
		{
			Translation: Translation{
				RussianName:        "`!pout`",
				RussianDescription: "Отправляет гифку с дующимся персонажем",
				EnglishName:        "`!pout`",
				EnglishDescription: "Send gif with pout",
			},
			Command:  "!pout",
			Category: gifCategory,
			Handler:  bot.Pout,
		},
	}

	bot.Commands = commands
}
