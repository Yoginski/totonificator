package main

import (
	"log"
	"totonificator/bindata"
	"totonificator/config"
	"totonificator/face"
	"totonificator/huificate"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {
	config, err := config.Get()
	if err != nil {
		log.Fatalf("Failed to get config: %s\n", err)
	}

	imageName := "totonia.png"
	imageBytes, err := bindata.Asset(imageName)
	if err != nil {
		log.Fatalf("Failed to load image (%q): %s\n", imageName, err)
	}
	maker, err := face.NewFaceMaker(imageBytes)
	if err != nil {
		log.Fatalf("Failed to initialize facemaker: %s\n", err)
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

		if update.Message.From.ID == 151262755 {
			log.Println("Работаю, работаю!")
		}

		if update.Message.From.ID != 144797944 {
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		response := huificate.Huify(update.Message.Text, 1)
		if response == "" {
			log.Printf("can't huify: %q\n", update.Message.Text)
			//response = "Отъебись, пидор!"
			continue
		}
		image, err := maker.Make(response, "Roboto-Regular.ttf", "black", 36)
		if err != nil {
			log.Printf("can't create picture: %s\n", err)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Смотри логи, криворукий!")
			msg.ReplyToMessageID = update.Message.MessageID
			bot.Send(msg)
			continue
		}
		photoUpload := tgbotapi.NewPhotoUpload(
			update.Message.Chat.ID,
			tgbotapi.FileBytes{Name: "totonia", Bytes: image},
		)
		photoUpload.ReplyToMessageID = update.Message.MessageID
		bot.Send(photoUpload)
	}
}
