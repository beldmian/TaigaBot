package bot

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
)

// OnMessage provide handler for MessageCreate event
func OnMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	content := strings.ToLower(m.Content)
	if strings.HasPrefix(content, "!") {
		logger.Info("Execute command", zap.String("command", content))
		if strings.HasPrefix(content, "!help") {
			go Help(s, m)
		} else if content == "!colors" {
			go ColorsList(s, m)
		} else if strings.HasPrefix(content, "!color ") {
			go PickColor(s, m)
		} else if strings.HasPrefix(content, "!delete ") {
			go BulkDelete(s, m)
		} else if strings.HasPrefix(content, "!massrole ") {
			go MassRole(s, m)
		} else if strings.HasPrefix(content, "!anime ") {
			go GetAnime(s, m)
		} else if content == "!tasks" {
			go Tasks(s, m)
		} else if strings.HasPrefix(content, "!task add ") {
			go TaskAdd(s, m)
		} else if strings.HasPrefix(content, "!task done ") {
			go TaskDone(s, m)
		}
	}
}

// OnBan provide handler for GuildBanAdd event
func OnBan(s *discordgo.Session, m *discordgo.GuildBanAdd) {
	logger.Info("User banned", zap.String("user_id", m.User.ID))
	if _, err := s.ChannelMessageSendEmbed(logsID, &discordgo.MessageEmbed{
		Title: m.User.Username + " был забанен на сервере",
		Color: 2343740,
	}); err != nil {
		SendErrorMessage(s, err)
	}
}
