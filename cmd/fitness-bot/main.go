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

	appHandler := application.NewAppHandler(tgBot.BotApi)

	foodHandler := application.NewFoodHandler(tgBot.BotApi)
	actHandler := application.NewActHandler(tgBot.BotApi)
	waterHandler := application.NewWaterHandler(tgBot.BotApi)

	tgBot.BotApi.Debug = false
	log.Println("Bot started")

	//wNot := repositories.NewWaterNotificationRepo(db.SQL)

	waterHandler.StartWaterReminders(tgBot.BotApi, userRepo)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 30
	updates := tgBot.BotApi.GetUpdatesChan(u)

	callbackHandler := application.NewCallbackHandler(tgBot.BotApi)

	for upd := range updates {
		if upd.Message != nil {
			go func() {
				appHandler.HandleMessage(tgBot.BotApi, upd.Message, actRepo, userRepo,
					mealRepo, weightRepo, foodHandler, actHandler, waterHandler)
			}()
		}

		if upd.CallbackQuery != nil {
			go func() {
				callbackHandler.HandleCallback(tgBot.BotApi, upd.CallbackQuery, userRepo, waterHandler, actHandler, appHandler)
			}()
		}
	}
}
