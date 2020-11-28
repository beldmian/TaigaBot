package bot

import (
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
)

var logsID string
var discord *discordgo.Session

// InitBot initializes bot process
func InitBot() {
	token, exists := os.LookupEnv("TOKEN")
	logsID, exists = os.LookupEnv("LOGS_ID")
	if !exists {
		log.Print("No token or logs channel ID provided")
		return
	}
	discord, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal(err)
	}
	discord.AddHandler(OnMessage)
	discord.AddHandler(OnBan)
	discord.AddHandler(OnMemberRemove)

	if err := discord.Open(); err != nil {
		log.Fatal(err)
	}
}

// StopBot stops the bot session
func StopBot() {
	if err := discord.Close(); err != nil {
		log.Fatal(err)
	}
}
