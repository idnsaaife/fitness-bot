package application

import (
	"fitness-bot/internal/adapter"
	"fitness-bot/internal/adapter/repositories"
	"fitness-bot/internal/domain"
	"fmt"
	"strconv"
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

	if *user.GetId() == 0 {
		msg := tgbotapi.NewMessage(msg.Chat.ID, "Вы не зарегистрированы. Нажмите кнопку ниже для начала:")
		msg.ReplyMarkup = keyboardHandler.StartKeyboard()
		bot.Send(msg)
		return
	}

	if actHandler.IsAddingActivity(msg.Chat.ID) {
		actHandler.HandleActivityDuration(msg, user, actRepo, userRepo, appHandler)
		return
	}

	if foodHandler.IsAddingFood(msg.Chat.ID) {
		foodHandler.HandleFoodInput(msg, user, mealRepo, userRepo, appHandler)
		return
	}

	switch {
	case text == "/start" || text == "🏠 Главное меню":
		appHandler.ShowMainMenu(bot, msg)

	case text == "📊 Статистика" || strings.HasPrefix(text, "/stats"):
		appHandler.StatsHandler(bot, msg, user, weightRepo, actRepo)

	case text == "🍎 Добавить еду" || strings.HasPrefix(text, "/addfood"):
		foodHandler.AddFoodHandler(msg)

	case text == "💧 Вода" || strings.HasPrefix(text, "/water"):
		waterHandler.HandlerWater(msg)

	case text == "🏃 Активность" || strings.HasPrefix(text, "/addactivity"):
		actHandler.ActivityHandler(msg)

	case text == "✏️ Редактировать данные" || strings.HasPrefix(text, "/edit"):
		appHandler.EditHandler(bot, msg, user, userRepo, actHandler)

	case text == "📋 Проверить питание" || strings.HasPrefix(text, "/checkfood"):
		foodHandler.CheckFoodHandler(msg, user, userRepo, mealRepo, appHandler)

	default:
		appHandler.Reply(bot, msg, "Не понял команду. Используйте кнопки меню:")
		appHandler.ShowMainMenu(bot, msg)
	}
}

func (appHandler *AppHandler) ShowMainMenu(bot *tgbotapi.BotAPI, msg *tgbotapi.Message) {
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

func (appHandler *AppHandler) EditHandler(bot *tgbotapi.BotAPI, msg *tgbotapi.Message, u domain.User,
	uRepo *repositories.UserRepo, actHandler *ActHandler) {
	args := strings.Fields(msg.Text)
	if len(args) < 6 {
		appHandler.Reply(bot, msg, "Использование:\n/edit <рост_cm> <вес_kg> <возраст> <цель> <подвижность>\nЦель: похудеть|набрать|оставить\nПодвижность: низкая|средняя|приемлемая|высокая\nПример: /edit 170 65 28 похудеть средняя")
		return
	}

	height, err := strconv.Atoi(args[1])
	if err != nil {
		appHandler.Reply(bot, msg, "Неверный рост")
		return
	}
	weight, err := strconv.ParseFloat(args[2], 64)
	if err != nil {
		appHandler.Reply(bot, msg, "Неверный вес")
		return
	}
	age, err := strconv.Atoi(args[3])
	if err != nil {
		appHandler.Reply(bot, msg, "Неверный возраст")
		return
	}

	goalRaw := strings.ToLower(args[4])
	var goal domain.Goal
	switch goalRaw {
	case "похудеть", "lose":
		goal = domain.GoalLose
	case "набрать", "gain":
		goal = domain.GoalGain
	default:
		goal = domain.GoalMaintain
	}

	actRaw := strings.ToLower(args[5])
	var act domain.ActivityLevel
	switch actRaw {
	case "низкая":
		act = domain.ActivityLow
	case "средняя":
		act = domain.ActivityMedium
	case "приемлемая":
		act = domain.ActivityGood
	case "высокая":
		act = domain.ActivityHigh
	default:
		act = domain.ActivityMedium
	}

	err = uRepo.UpdateUserParams(height, weight, age, string(goal), string(act), *u.GetId())
	if err != nil {
		appHandler.Reply(bot, msg, "Ошибка обновления данных")
		return
	}

	u.SetHeightCm(height)
	u.SetWeightKg(weight)
	u.SetAge(age)
	u.SetGoal(goal)
	u.SetActivityLevel(act)

	newCal := actHandler.CalcDailyCalories(u)

	uRepo.UpdateGoalCalories(newCal, *u.GetId())

	appHandler.Reply(bot, msg, fmt.Sprintf("Данные обновлены. Новая дневная норма: %d ккал", newCal))
}
