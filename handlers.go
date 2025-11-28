package main

import (
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HandleMessage(bot *tgbotapi.BotAPI, msg *tgbotapi.Message) {
	text := msg.Text
	user, err := GetUserByTelegramID(msg.From.ID)
	
	if HandleRegistration(bot, msg) {
		return
	}
	
	if err != nil || user.ID == 0 {
		if text != "/start" {
			reply(bot, msg, "Вы не зарегистрированы. Напишите /start для регистрации.")
			return
		}
	}

	switch {
	case text == "/start":
		StartHandler(bot, msg)
	case strings.HasPrefix(text, "/edit"):
		EditHandler(bot, msg, user)
	case strings.HasPrefix(text, "/addfood"):
		AddFoodHandler(bot, msg, user)
	case strings.HasPrefix(text, "/checkfood"):
		CheckFoodHandler(bot, msg, user)
	case strings.HasPrefix(text, "/addactivity"):
		AddActivityHandler(bot, msg, user)
	case strings.HasPrefix(text, "/water"):
		WaterCommandHandler(bot, msg, user)
	case strings.HasPrefix(text, "/stats"):
		StatsHandler(bot, msg, user)
	default:
		reply(bot, msg, "Не понял команду. Напиши /help")
	}
}


func reply(bot *tgbotapi.BotAPI, msg *tgbotapi.Message, text string) {
	bot.Send(tgbotapi.NewMessage(msg.Chat.ID, text))
}
