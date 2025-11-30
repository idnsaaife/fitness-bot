package main

import (
	"fitness-bot/internal/adapter"
	"fitness-bot/internal/adapter/repositories"
	"fitness-bot/internal/application"
	"log"

	"github.com/joho/godotenv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	_ = godotenv.Load(".env")
	//adapter.InitDB()

	db, err := adapter.InitDB()
	if err != nil {
		log.Fatal(err)
	}

	actRepo := repositories.NewActivityRepo(db.SQL)
	userRepo := repositories.NewUserRepo(db.SQL)
	weightRepo := repositories.NewWeightChangesRepo(db.SQL)
	mealRepo := repositories.NewMealRepo(db.SQL)

	tgBot := adapter.NewBotApi()

	tgBot.BotApi.Debug = false
	log.Println("Bot started")

	//wNot := repositories.NewWaterNotificationRepo(db.SQL)

	application.StartWaterReminders(tgBot.BotApi, userRepo)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 30
	updates := tgBot.BotApi.GetUpdatesChan(u)

	for upd := range updates {
		if upd.Message != nil {
			go func() {
				application.HandleMessage(tgBot.BotApi, upd.Message, actRepo, userRepo, mealRepo, weightRepo)
			}()
		}

		if upd.CallbackQuery != nil {
			go func() {
				application.HandleCallback(tgBot.BotApi, upd.CallbackQuery, userRepo)
			}()
		}
	}
}
