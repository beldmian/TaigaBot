package main

import (
	"github.com/bwmarrin/discordgo"
	"log"
	"strings"
)

func OnMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Content == "!colors" {
		go ColorsList(s, m)
	} else if strings.HasPrefix(m.Content, "!color ") {
		go PickColor(s, m)
	} else if strings.HasPrefix(m.Content, "!delete ") {
		go BulkDelete(s, m)
	} else if strings.HasPrefix(m.Content, "!massrole ") {
		go MassRole(s, m)
	}
}
func OnBan(s *discordgo.Session, m *discordgo.GuildBanAdd) {
	if _, err := s.ChannelMessageSendEmbed(LogsId, &discordgo.MessageEmbed{
		Title: m.User.Username + " был забанен на сервере",
		Color: 2343740,
	}); err != nil {
		log.Fatal(err)
	}
}
func OnMemberRemove(s *discordgo.Session, m *discordgo.GuildMemberRemove) {
	if _, err := s.ChannelMessageSendEmbed(LogsId, &discordgo.MessageEmbed{
		Title: m.User.Username + " больше не на сервере",
		Color: 2343740,
	}); err != nil {
		log.Fatal(err)
	}
}
