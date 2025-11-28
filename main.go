package main

import (
	"log"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	InitDB()

	bot, err := tgbotapi.NewBotAPI("8553264931:AAF_-yplXOLUAAIkF39rgP98oEc_34TD7bw")
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = false
	log.Println("Bot started")

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 30
	StartWaterReminders(bot)
	updates := bot.GetUpdatesChan(u)

	for upd := range updates {
    if upd.Message != nil {
        go HandleMessage(bot, upd.Message)
    }

    if upd.CallbackQuery != nil {
        go HandleCallback(bot, upd.CallbackQuery)
    }
	}
}
