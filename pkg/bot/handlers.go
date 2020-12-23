package bot

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
)

// OnMessage provide handler for MessageCreate event
func (bot *Bot) OnMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	content := strings.ToLower(m.Content)
	if strings.HasPrefix(content, "!") {
		bot.Logger.Info("Execute command", zap.String("command", content))
		if strings.HasPrefix(content, "!help") {
			go bot.Help(s, m)
		} else if content == "!colors" {
			go bot.ColorsList(s, m)
		} else if strings.HasPrefix(content, "!color ") {
			go bot.PickColor(s, m)
		} else if strings.HasPrefix(content, "!delete ") {
			go bot.BulkDelete(s, m)
		} else if strings.HasPrefix(content, "!massrole ") {
			go bot.MassRole(s, m)
		} else if strings.HasPrefix(content, "!anime ") {
			go bot.GetAnime(s, m)
		} else if content == "!tasks" {
			go bot.Tasks(s, m)
		} else if strings.HasPrefix(content, "!task add ") {
			go bot.TaskAdd(s, m)
		} else if strings.HasPrefix(content, "!task done ") {
			go bot.TaskDone(s, m)
		}
	}
}

// OnBan provide handler for GuildBanAdd event
func (bot *Bot) OnBan(s *discordgo.Session, m *discordgo.GuildBanAdd) {
	bot.Logger.Info("User banned", zap.String("user_id", m.User.ID))
	if _, err := s.ChannelMessageSendEmbed(bot.LogsID, &discordgo.MessageEmbed{
		Title: m.User.Username + " был забанен на сервере",
		Color: 2343740,
	}); err != nil {
		bot.SendErrorMessage(s, err)
	}
}
