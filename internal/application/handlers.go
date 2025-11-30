package application

import (
	"fitness-bot/internal/adapter"
	"fitness-bot/internal/adapter/repositories"
	"fitness-bot/internal/domain"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HandleMessage(bot *tgbotapi.BotAPI, msg *tgbotapi.Message, actRepo *repositories.ActivityRepo,
	userRepo *repositories.UserRepo, mealRepo *repositories.MealRepo, weightRepo *repositories.WeightChangesRepo) {
	text := msg.Text

	if isProcessing := HandleRegistration(bot, msg); isProcessing {
		return
	}

	if text == "/start" {
		StartHandler(bot, msg, userRepo)
		return
	}

	user, err := userRepo.GetUserByTelegramID(msg.From.ID)
	if err != nil {
		Reply(bot, msg, "Ошибка базы данных. Попробуйте позже.")
		return
	}

	if user.ID == 0 {
		msg := tgbotapi.NewMessage(msg.Chat.ID, "Вы не зарегистрированы. Нажмите кнопку ниже для начала:")
		msg.ReplyMarkup = adapter.StartKeyboard()
		bot.Send(msg)
		return
	}

	if IsAddingActivity(msg.Chat.ID) {
		HandleActivityDuration(bot, msg, user, actRepo, userRepo)
		return
	}

	if IsAddingFood(msg.Chat.ID) {
		HandleFoodInput(bot, msg, user, mealRepo, userRepo)
		return
	}

	switch {
	case text == "/start" || text == "🏠 Главное меню":
		ShowMainMenu(bot, msg, user)

	case text == "📊 Статистика" || strings.HasPrefix(text, "/stats"):
		StatsHandler(bot, msg, user, weightRepo, actRepo)

	case text == "🍎 Добавить еду" || strings.HasPrefix(text, "/addfood"):
		AddFoodHandler(bot, msg, user)

	case text == "💧 Вода" || strings.HasPrefix(text, "/water"):
		WaterHandler(bot, msg, user)

	case text == "🏃 Активность" || strings.HasPrefix(text, "/addactivity"):
		ActivityHandler(bot, msg, user)

	case text == "✏️ Редактировать данные" || strings.HasPrefix(text, "/edit"):
		EditHandler(bot, msg, user, userRepo)

	case text == "📋 Проверить питание" || strings.HasPrefix(text, "/checkfood"):
		CheckFoodHandler(bot, msg, user, userRepo, mealRepo)

	default:
		Reply(bot, msg, "Не понял команду. Используйте кнопки меню:")
		ShowMainMenu(bot, msg, user)
	}
}

func ShowMainMenu(bot *tgbotapi.BotAPI, msg *tgbotapi.Message, u domain.User) {
	text := `🏠 *Главное меню*

Выберите действие:`

	msgOut := tgbotapi.NewMessage(msg.Chat.ID, text)
	msgOut.ReplyMarkup = adapter.MainMenuKeyboard()
	msgOut.ParseMode = "Markdown"
	bot.Send(msgOut)
}

func Reply(bot *tgbotapi.BotAPI, msg *tgbotapi.Message, text string) {
	message := tgbotapi.NewMessage(msg.Chat.ID, text)
	bot.Send(message)
}
