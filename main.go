package main

import (
	"log"
	"totonificator/config"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {
	config, err := config.Get()
	if err != nil {
		log.Fatalf("Failed to get config: %s\n", err)
	}

	bot, err := tgbotapi.NewBotAPI(config.Token)
	if err != nil {
		log.Fatalf("Failed to create bot: %s\n", err)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		msg.ReplyToMessageID = update.Message.MessageID

		bot.Send(msg)
	}
}
