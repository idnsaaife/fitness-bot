package adapter

import (
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TgBot struct {
	botApi *tgbotapi.BotAPI
}

func (tgBot *TgBot) GetTgBotApi() *tgbotapi.BotAPI {
	return tgBot.botApi
}

func NewBotApi() *TgBot {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("BOT_TOKEN"))
	if err != nil {
		log.Fatal(err)
	}
	return &TgBot{botApi: bot}
}
