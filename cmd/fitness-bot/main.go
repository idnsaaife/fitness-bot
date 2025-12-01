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

	actRepo := repositories.NewActivityRepo(db.GetDB())
	userRepo := repositories.NewUserRepo(db.GetDB())
	weightRepo := repositories.NewWeightChangesRepo(db.GetDB())
	mealRepo := repositories.NewMealRepo(db.GetDB())

	tgBot := adapter.NewBotApi()

	appHandler := application.NewAppHandler(tgBot.GetTgBotApi())

	foodHandler := application.NewFoodHandler(tgBot.GetTgBotApi())
	actHandler := application.NewActHandler(tgBot.GetTgBotApi())
	waterHandler := application.NewWaterHandler(tgBot.GetTgBotApi())

	tgBot.GetTgBotApi().Debug = false
	log.Println("Bot started")

	//wNot := repositories.NewWaterNotificationRepo(db.SQL)

	waterHandler.StartWaterReminders(tgBot.GetTgBotApi(), userRepo)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 30
	updates := tgBot.GetTgBotApi().GetUpdatesChan(u)

	callbackHandler := application.NewCallbackHandler(tgBot.GetTgBotApi())

	for upd := range updates {
		if upd.Message != nil {
			go func() {
				appHandler.HandleMessage(tgBot.GetTgBotApi(), upd.Message, actRepo, userRepo,
					mealRepo, weightRepo, foodHandler, actHandler, waterHandler)
			}()
		}

		if upd.CallbackQuery != nil {
			go func() {
				callbackHandler.HandleCallback(tgBot.GetTgBotApi(), upd.CallbackQuery, userRepo, waterHandler, actHandler, appHandler)
			}()
		}
	}
}
