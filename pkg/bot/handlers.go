package bot

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
)

// OnMessage provide handler for MessageCreate event
func (bot *Bot) OnMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	content := strings.ToLower(m.Content)
	for _, command := range bot.Commands {
		if strings.HasPrefix(content, command.Command) {
			bot.Logger.Info("Execute command", zap.String("command", content))
			command.Handler(s, m)
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
