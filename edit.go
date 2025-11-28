package main

import (
	"fmt"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func EditHandler(bot *tgbotapi.BotAPI, msg *tgbotapi.Message, u User) {
	args := strings.Fields(msg.Text)
	if len(args) < 6 {
		reply(bot, msg, "Использование:\n/edit <рост_cm> <вес_kg> <возраст> <цель> <подвижность>\nЦель: похудеть|набрать|оставить\nПодвижность: низкая|средняя|приемлемая|высокая\nПример: /edit 170 65 28 похудеть средняя")
		return
	}

	height, err := strconv.Atoi(args[1])
	if err != nil {
		reply(bot, msg, "Неверный рост")
		return
	}
	weight, err := strconv.ParseFloat(args[2], 64)
	if err != nil {
		reply(bot, msg, "Неверный вес")
		return
	}
	age, err := strconv.Atoi(args[3])
	if err != nil {
		reply(bot, msg, "Неверный возраст")
		return
	}

	goalRaw := strings.ToLower(args[4])
	var goal Goal
	switch goalRaw {
	case "похудеть", "lose":
		goal = GoalLose
	case "набрать", "gain":
		goal = GoalGain
	default:
		goal = GoalMaintain
	}

	actRaw := strings.ToLower(args[5])
	var act ActivityLevel
	switch actRaw {
	case "низкая":
		act = ActivityLow
	case "средняя":
		act = ActivityMedium
	case "приемлемая":
		act = ActivityGood
	case "высокая":
		act = ActivityHigh
	default:
		act = ActivityMedium
	}

	_, err = DB.Exec("UPDATE users SET height_cm = ?, weight_kg = ?, age = ?, goal = ?, activity_level = ? WHERE id = ?", height, weight, age, string(goal), string(act), u.ID)
	if err != nil {
		reply(bot, msg, "Ошибка обновления данных")
		return
	}

	// пересчитать дневную норму
	u.HeightCm = height
	u.WeightKg = weight
	u.Age = age
	u.Goal = goal
	u.ActivityLevel = act

	newCal := CalcDailyCalories(u)
	_, _ = DB.Exec("UPDATE users SET calories_goal = ? WHERE id = ?", newCal, u.ID)

	reply(bot, msg, fmt.Sprintf("Данные обновлены. Новая дневная норма: %d ккал", newCal))
}
