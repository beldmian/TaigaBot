package main

import (
	"github.com/bwmarrin/discordgo"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var LogsId string

func main() {
	token, exists := os.LookupEnv("TOKEN")
	LogsId, exists = os.LookupEnv("LOGS_ID")
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

	log.Print("Started")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	if err := discord.Close(); err != nil {
		log.Fatal(err)
	}
}
