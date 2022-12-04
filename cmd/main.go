package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"stabilizationTimeBot/pkg/core/bot"
	"stabilizationTimeBot/pkg/core/puller"
)

func main() {
	token := os.Getenv("token")
	chatID, err := strconv.ParseInt(os.Getenv("chatID"), 10, 64)
	if err != nil {
		log.Panic(err)
	}

	b := bot.NewBot(token, chatID)
	p := puller.NewPuller(b)

	quitCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	p.Run(quitCtx)
}
