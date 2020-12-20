package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/beldmian/TaigaBot/pkg/bot"
)

func main() {
	bot.InitBot()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	bot.StopBot()
}
