package application

import (
	"fitness-bot/internal/adapter/repositories"
	"fitness-bot/internal/domain"
	"fmt"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

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

	err = uRepo.UpdateUserParams(height, weight, age, string(goal), string(act), u.ID)
	//_, err = uRepo.Db.Exec("UPDATE users SET height_cm = ?, weight_kg = ?, age = ?, goal = ?, activity_level = ? WHERE id = ?", height, weight, age, string(goal), string(act), u.ID)
	if err != nil {
		appHandler.Reply(bot, msg, "Ошибка обновления данных")
		return
	}

	u.HeightCm = height
	u.WeightKg = weight
	u.Age = age
	u.Goal = goal
	u.ActivityLevel = act

	newCal := actHandler.CalcDailyCalories(u)

	uRepo.UpdateGoalCalories(newCal, u.ID)
	//_, _ = uRepo.Db.Exec("UPDATE users SET calories_goal = ? WHERE id = ?", newCal, u.ID)

	appHandler.Reply(bot, msg, fmt.Sprintf("Данные обновлены. Новая дневная норма: %d ккал", newCal))
}
