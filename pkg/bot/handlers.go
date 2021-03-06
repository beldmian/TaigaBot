package bot

import (
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
)

// OnMessage provide handler for MessageCreate event
func (bot *Bot) OnMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	content := strings.ToLower(m.Content)
	if m.Author.Bot {
		return
	}
	if m.GuildID == "" {
		s.ChannelMessageSend(m.ChannelID, "Не работаю в личных сообщениях")
		return
	}
	s.State.MemberAdd(m.Member)
	s.State.MessageAdd(m.Message)
	prefix, err := bot.DB.GetPrefix(m.GuildID)
	if err != nil {
		bot.SendErrorMessage(s, err)
		return
	}
	if !strings.HasPrefix(content, prefix) {
		return
	}
	content = strings.TrimPrefix(content, prefix)
	log.Print(content)

	for _, command := range bot.Commands {
		if strings.HasPrefix(content, command.Command) {
			if command.Category == moderationCategory {
				permissions, err := s.State.UserChannelPermissions(m.Author.ID, m.ChannelID)
				if err != nil && err != discordgo.ErrStateNotFound {
					bot.SendErrorMessage(s, err)
				}
				if command.Permissions&permissions != command.Permissions {
					s.ChannelMessageSend(m.ChannelID, "У вас недостаточно прав")
					return
				}
			}
			if command.Category == nsfwCategory {
				ch, err := s.Channel(m.ChannelID)
				if err != nil {
					bot.SendErrorMessage(s, err)
				}
				if !ch.NSFW {
					s.ChannelMessageSend(m.ChannelID, "Только nsfw каналы")
					return
				}
			}
			bot.Logger.Info("Execute command", zap.String("command", content), zap.String("guild_id", m.GuildID))
			locale, err := bot.GetServerLocal(s, m.GuildID)
			if err != nil {
				bot.SendErrorMessage(s, err)
				return
			}
			command.Handler(s, m, locale)
		}
	}
}

// OnBan provide handler for GuildBanAdd event
func (bot *Bot) OnBan(s *discordgo.Session, m *discordgo.GuildBanAdd) {
	g, err := s.Guild(m.GuildID)
	if err != nil {
		bot.SendErrorMessage(s, err)
	}
	logChannel := g.SystemChannelID
	bot.Logger.Info("User banned", zap.String("user_id", m.User.ID))
	if _, err := s.ChannelMessageSendEmbed(logChannel, &discordgo.MessageEmbed{
		Title: m.User.Username + " был забанен на сервере",
		Color: 2343740,
	}); err != nil {
		return
	}
}
