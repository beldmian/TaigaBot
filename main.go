package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/BurntSushi/toml"
	"github.com/beldmian/TaigaBot/pkg/bot"
	"github.com/beldmian/TaigaBot/pkg/types"
)

func main() {
	var config types.Config
	if _, err := toml.DecodeFile("./config.toml", &config); err != nil {
		log.Fatal(err)
	}
	myBot := bot.InitBot(config)

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	myBot.StopBot()
}
