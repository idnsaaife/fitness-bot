package application

import (
	"fitness-bot/internal/adapter/repositories"
	"fitness-bot/internal/domain"
	"fmt"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Хранилище состояний и временных данных для добавления еды

var (
	foodStates   = map[int64]string{} // chatID -> "waiting_calories" или "waiting_description"
	foodTempData = map[int64]int{}    // chatID -> calories
)

type FoodHandler struct {
	bot *tgbotapi.BotAPI
}

func NewFoodHandler(Bot *tgbotapi.BotAPI) *FoodHandler {
	return &FoodHandler{bot: Bot}
}

func (foodHandler *FoodHandler) AddFoodHandler(msg *tgbotapi.Message) {
	foodStates[msg.Chat.ID] = "waiting_calories"

	text := `🍎 *Добавление еды*

Введите количество калорий:
Пример: *250*`

	msgOut := tgbotapi.NewMessage(msg.Chat.ID, text)
	msgOut.ParseMode = "Markdown"
	msgOut.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true) // Убираем клавиатуру для ввода
	foodHandler.bot.Send(msgOut)
}

func (foodHandler *FoodHandler) HandleFoodInput(msg *tgbotapi.Message, u domain.User,
	mealRepo *repositories.MealRepo, userRepo *repositories.UserRepo, appHandler *AppHandler) {
	state, exists := foodStates[msg.Chat.ID]
	if !exists {
		return
	}

	text := strings.TrimSpace(msg.Text)

	switch state {
	case "waiting_calories":
		kcal, err := strconv.Atoi(text)
		if err != nil || kcal <= 0 || kcal > 5000 {
			appHandler.Reply(foodHandler.bot, msg, "Неверный формат калорий. Введите число от 1 до 5000:")
			return
		}

		foodStates[msg.Chat.ID] = "waiting_description"
		foodTempData[msg.Chat.ID] = kcal
		appHandler.Reply(foodHandler.bot, msg, "Теперь введите описание еды:\nПример: *Яблоко* или *Овсяная каша*")
		return

	case "waiting_description":

		kcal, exists := foodTempData[msg.Chat.ID]
		if !exists {
			delete(foodStates, msg.Chat.ID)
			appHandler.ShowMainMenu(foodHandler.bot, msg)
			return
		}

		desc := text

		err := mealRepo.SaveFoodWithCalories(*u.GetId(), desc, kcal)
		if err != nil {
			appHandler.Reply(foodHandler.bot, msg, "Ошибка сохранения еды")
			delete(foodStates, msg.Chat.ID)
			delete(foodTempData, msg.Chat.ID)
			appHandler.ShowMainMenu(foodHandler.bot, msg)
			return
		}

		userRepo.UpdateCalories(kcal, *u.GetId())
		appHandler.Reply(foodHandler.bot, msg, fmt.Sprintf("✅ Добавлено: *%s* — *%d ккал*", desc, kcal))

		delete(foodStates, msg.Chat.ID)
		delete(foodTempData, msg.Chat.ID)

		appHandler.ShowMainMenu(foodHandler.bot, msg)
		return
	}
}

func (foodHandler *FoodHandler) IsAddingFood(chatID int64) bool {
	state, exists := foodStates[chatID]
	return exists && (state == "waiting_calories" || state == "waiting_description")
}

func (foodHandler *FoodHandler) CheckFoodHandler(msg *tgbotapi.Message, u domain.User,
	userRepo *repositories.UserRepo, mealRepo *repositories.MealRepo, appHandler *AppHandler) {
	startOfDay := time.Now().Format("2006-01-02") + " 00:00:00"
	rows, err := mealRepo.GetAllFoodByDay(*u.GetId(), startOfDay)
	if err != nil {
		appHandler.Reply(foodHandler.bot, msg, "Ошибка чтения базы")
		return
	}
	defer rows.Close()

	var total int
	text := "Еда сегодня:\n"
	for rows.Next() {
		var desc string
		var kcal int
		var createdAt string
		rows.Scan(&desc, &kcal, &createdAt)
		text += fmt.Sprintf("- %s: %d ккал\n", desc, kcal)
		total += kcal
	}

	if *u.GetCaloriesGoal() == 0 {
		u.SetCaloriesGoal(1000)
		userRepo.UpdateGoalCalories(*u.GetCaloriesGoal(), *u.GetId())
	}

	remaining := *u.GetCaloriesGoal() - total
	if remaining < 0 {
		remaining = 0
	}

	text += fmt.Sprintf("\nВсего: %d ккал\nОсталось до дневной нормы (%d): %d ккал", total, *u.GetCaloriesGoal(), remaining)
	appHandler.Reply(foodHandler.bot, msg, text)
}
