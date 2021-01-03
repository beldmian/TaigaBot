package bot

import "github.com/bwmarrin/discordgo"

// Command provide struct for commands
type Command struct {
	Name        string
	Description string
	Command     string
	Moderation  bool
	Permissions int
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
			Permissions: 8192,
		},
		{
			Name:        "`!massrole @<роль>`",
			Description: "Выдает или забирает роль у всех на сервере",
			Command:     "!massrole ",
			Moderation:  true,
			Handler:     bot.MassRole,
			Permissions: 268435456,
		},
		{
			Name:        "`!poll <вариант 1> | <вариант 2> ...`",
			Description: "Создает опрос с несколькими вариантами ответа",
			Command:     "!poll ",
			Moderation:  false,
			Handler:     bot.Poll,
		},
		{
			Name:        "`!vote",
			Description: "Вы можете поддержать бота проголосовав за него на top.gg",
			Command:     "!vote",
			Moderation:  false,
			Handler:     bot.Vote,
		},
	}

	bot.Commands = commands
}
