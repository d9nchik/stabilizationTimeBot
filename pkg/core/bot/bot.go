package bot

import (
	"fmt"
	"log"
	"os"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	bot    *tgbotapi.BotAPI
	chatID int64
	tz     *time.Location
}

func NewBot(token string, chatID int64) *Bot {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	tz, err := time.LoadLocation("Europe/Kiev")
	if err != nil {
		log.Panic(err)
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	return &Bot{bot: bot, chatID: chatID, tz: tz}
}

func (b *Bot) SendFile(filename string) bool {
	file, err := os.Open(filename)
	if err != nil {
		log.Printf("Error opening file: %v", err)
		return false
	}

	f := tgbotapi.FileReader{
		Name:   fmt.Sprintf("Графіки погодиних відключень%s.pdf", time.Now().In(b.tz).Format("2006-01-02")),
		Reader: file,
	}

	doc := tgbotapi.NewDocument(b.chatID, f)

	if _, err := b.bot.Send(doc); err != nil {
		log.Printf("Error sending file: %v", err)
		return false
	}
	return true
}
