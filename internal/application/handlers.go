package application

import (
	"fitness-bot/internal/adapter"
	"fitness-bot/internal/adapter/repositories"
	"fitness-bot/internal/domain"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type AppHandler struct {
	bot *tgbotapi.BotAPI
}

func NewAppHandler(Bot *tgbotapi.BotAPI) *AppHandler {
	return &AppHandler{bot: Bot}
}

func (appHandler *AppHandler) HandleMessage(bot *tgbotapi.BotAPI, msg *tgbotapi.Message, actRepo *repositories.ActivityRepo,
	userRepo *repositories.UserRepo, mealRepo *repositories.MealRepo,
	weightRepo *repositories.WeightChangesRepo, foodHandler *FoodHandler,
	actHandler *ActHandler, waterHandler *WaterHandler) {
	keyboardHandler := adapter.NewKeyboardHandler()
	text := msg.Text

	if isProcessing := appHandler.HandleRegistration(bot, msg); isProcessing {
		return
	}

	if text == "/start" {
		appHandler.StartHandler(bot, msg, userRepo)
		return
	}

	user, err := userRepo.GetUserByTelegramID(msg.From.ID)
	if err != nil {
		appHandler.Reply(bot, msg, "Ошибка базы данных. Попробуйте позже.")
		return
	}

	if user.ID == 0 {
		msg := tgbotapi.NewMessage(msg.Chat.ID, "Вы не зарегистрированы. Нажмите кнопку ниже для начала:")
		msg.ReplyMarkup = keyboardHandler.StartKeyboard()
		bot.Send(msg)
		return
	}

	if actHandler.IsAddingActivity(msg.Chat.ID) {
		actHandler.HandleActivityDuration(bot, msg, user, actRepo, userRepo, appHandler)
		return
	}

	if foodHandler.IsAddingFood(msg.Chat.ID) {
		foodHandler.HandleFoodInput(bot, msg, user, mealRepo, userRepo, appHandler)
		return
	}

	switch {
	case text == "/start" || text == "🏠 Главное меню":
		appHandler.ShowMainMenu(bot, msg, user)

	case text == "📊 Статистика" || strings.HasPrefix(text, "/stats"):
		appHandler.StatsHandler(bot, msg, user, weightRepo, actRepo)

	case text == "🍎 Добавить еду" || strings.HasPrefix(text, "/addfood"):
		foodHandler.AddFoodHandler(bot, msg, user)

	case text == "💧 Вода" || strings.HasPrefix(text, "/water"):
		waterHandler.HandlerWater(bot, msg, user)

	case text == "🏃 Активность" || strings.HasPrefix(text, "/addactivity"):
		actHandler.ActivityHandler(bot, msg, user)

	case text == "✏️ Редактировать данные" || strings.HasPrefix(text, "/edit"):
		appHandler.EditHandler(bot, msg, user, userRepo, actHandler)

	case text == "📋 Проверить питание" || strings.HasPrefix(text, "/checkfood"):
		foodHandler.CheckFoodHandler(bot, msg, user, userRepo, mealRepo, appHandler)

	default:
		appHandler.Reply(bot, msg, "Не понял команду. Используйте кнопки меню:")
		appHandler.ShowMainMenu(bot, msg, user)
	}
}

func (appHandler *AppHandler) ShowMainMenu(bot *tgbotapi.BotAPI, msg *tgbotapi.Message, u domain.User) {
	text := `🏠 *Главное меню*

Выберите действие:`

	keyboardHandler := adapter.NewKeyboardHandler()
	msgOut := tgbotapi.NewMessage(msg.Chat.ID, text)
	msgOut.ReplyMarkup = keyboardHandler.MainMenuKeyboard()
	msgOut.ParseMode = "Markdown"
	bot.Send(msgOut)
}

func (appHandler *AppHandler) Reply(bot *tgbotapi.BotAPI, msg *tgbotapi.Message, text string) {
	message := tgbotapi.NewMessage(msg.Chat.ID, text)
	bot.Send(message)
}
