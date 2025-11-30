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

// Обработчик кнопки "Добавить еду"
func AddFoodHandler(bot *tgbotapi.BotAPI, msg *tgbotapi.Message, u domain.User) {
	foodStates[msg.Chat.ID] = "waiting_calories"

	text := `🍎 *Добавление еды*

Введите количество калорий:
Пример: *250*`

	msgOut := tgbotapi.NewMessage(msg.Chat.ID, text)
	msgOut.ParseMode = "Markdown"
	msgOut.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true) // Убираем клавиатуру для ввода
	bot.Send(msgOut)
}

// Обработчик ввода данных для еды
func HandleFoodInput(bot *tgbotapi.BotAPI, msg *tgbotapi.Message, u domain.User,
	mealRepo *repositories.MealRepo, userRepo *repositories.UserRepo) {
	state, exists := foodStates[msg.Chat.ID]
	if !exists {
		return
	}

	text := strings.TrimSpace(msg.Text)

	switch state {
	case "waiting_calories":
		kcal, err := strconv.Atoi(text)
		if err != nil || kcal <= 0 || kcal > 5000 {
			Reply(bot, msg, "Неверный формат калорий. Введите число от 1 до 5000:")
			return
		}

		foodStates[msg.Chat.ID] = "waiting_description"
		foodTempData[msg.Chat.ID] = kcal
		Reply(bot, msg, "Теперь введите описание еды:\nПример: *Яблоко* или *Овсяная каша*")
		return

	case "waiting_description":

		kcal, exists := foodTempData[msg.Chat.ID]
		if !exists {
			delete(foodStates, msg.Chat.ID)
			ShowMainMenu(bot, msg, u)
			return
		}

		desc := text

		err := mealRepo.SaveFoodWithCalories(u.ID, desc, kcal)
		//_, err := adapter.DB.Exec("INSERT INTO meals (user_id, description, calories) VALUES (?, ?, ?)", u.ID, desc, kcal)
		if err != nil {
			Reply(bot, msg, "Ошибка сохранения еды")
			delete(foodStates, msg.Chat.ID)
			delete(foodTempData, msg.Chat.ID)
			ShowMainMenu(bot, msg, u)
			return
		}

		userRepo.UpdateCalories(kcal, u.ID)
		//_, _ = adapter.DB.Exec("UPDATE users SET calories_today = calories_today + ? WHERE id = ?", kcal, u.ID)

		Reply(bot, msg, fmt.Sprintf("✅ Добавлено: *%s* — *%d ккал*", desc, kcal))

		delete(foodStates, msg.Chat.ID)
		delete(foodTempData, msg.Chat.ID)

		ShowMainMenu(bot, msg, u)
		return
	}
}

func IsAddingFood(chatID int64) bool {
	state, exists := foodStates[chatID]
	return exists && (state == "waiting_calories" || state == "waiting_description")
}

// //func no usages
func AddFoodCommandHandler(bot *tgbotapi.BotAPI, msg *tgbotapi.Message, u domain.User) {
	args := strings.SplitN(msg.Text, " ", 3)
	if len(args) < 3 {
		AddFoodHandler(bot, msg, u)
		return
	}

	kcal, err := strconv.Atoi(args[1])
	if err != nil {
		Reply(bot, msg, "Неверный формат калорий")
		return
	}
	desc := args[2]

	//_, err = adapter.DB.Exec("INSERT INTO meals (user_id, description, calories) VALUES (?, ?, ?)", u.ID, desc, kcal)
	if err != nil {
		Reply(bot, msg, "Ошибка сохранения еды")
		return
	}

	//_, _ = adapter.DB.Exec("UPDATE users SET calories_today = calories_today + ? WHERE id = ?", kcal, u.ID)

	Reply(bot, msg, fmt.Sprintf("Добавлено: %s — %d ккал", desc, kcal))
}

func CheckFoodHandler(bot *tgbotapi.BotAPI, msg *tgbotapi.Message, u domain.User,
	userRepo *repositories.UserRepo, mealRepo *repositories.MealRepo) {
	startOfDay := time.Now().Format("2006-01-02") + " 00:00:00"
	rows, err := mealRepo.GetAllFoodByDay(u.ID, startOfDay)
	//rows, err := adapter.DB.Query("SELECT description, calories, created_at FROM meals WHERE user_id = ? AND created_at >= ?", u.ID, startOfDay)
	if err != nil {
		Reply(bot, msg, "Ошибка чтения базы")
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

	if u.CaloriesGoal == 0 {
		u.CaloriesGoal = 1000
		userRepo.UpdateGoalCalories(u.CaloriesGoal, u.ID)
		//_, _ = userRepo.Db.Exec("UPDATE users SET calories_goal = ? WHERE id = ?", u.CaloriesGoal, u.ID)
	}

	remaining := u.CaloriesGoal - total
	if remaining < 0 {
		remaining = 0
	}

	text += fmt.Sprintf("\nВсего: %d ккал\nОсталось до дневной нормы (%d): %d ккал", total, u.CaloriesGoal, remaining)
	Reply(bot, msg, text)
}
