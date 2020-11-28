package bot

import (
	"github.com/bwmarrin/discordgo"
	"strings"
)

// OnMessage provide handler for MessageCreate event
func OnMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	s.State.MessageAdd(m.Message)
	if m.Content == "!help" {
		go Help(s, m)
	} else if m.Content == "!colors" {
		go ColorsList(s, m)
	} else if strings.HasPrefix(m.Content, "!color ") {
		go PickColor(s, m)
	} else if strings.HasPrefix(m.Content, "!delete ") {
		go BulkDelete(s, m)
	} else if strings.HasPrefix(m.Content, "!massrole ") {
		go MassRole(s, m)
	} else if strings.HasPrefix(m.Content, "!anime ") {
		go GetAnime(s, m)
	}
}

// OnBan provide handler for GuildBanAdd event
func OnBan(s *discordgo.Session, m *discordgo.GuildBanAdd) {
	if _, err := s.ChannelMessageSendEmbed(logsID, &discordgo.MessageEmbed{
		Title: m.User.Username + " был забанен на сервере",
		Color: 2343740,
	}); err != nil {
		SendErrorMessage(s, err)
	}
}

// OnMemberRemove provide handler for GuildMemberRemove event
func OnMemberRemove(s *discordgo.Session, m *discordgo.GuildMemberRemove) {
	if _, err := s.ChannelMessageSendEmbed(logsID, &discordgo.MessageEmbed{
		Title: m.User.Username + " больше не на сервере",
		Color: 2343740,
	}); err != nil {
		SendErrorMessage(s, err)
	}
}

// OnEdit provide handler for MessageEdit event
func OnEdit(s *discordgo.Session, m *discordgo.MessageUpdate) {
	before , err := s.State.Message(m.ChannelID, m.ID)
	if err != nil {
		SendErrorMessage(s, err)
	}
	if before == nil {return}
	s.ChannelMessageSendEmbed(logsID, &discordgo.MessageEmbed{
		Title: "Изменено сообщение",
		Fields: []*discordgo.MessageEmbedField{
			{
				Name: "Было",
				Value: before.Content,
			},
			{
				Name: "Стало",
				Value: m.Content,
			},
		},
	})
}