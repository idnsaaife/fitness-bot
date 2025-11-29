package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func AddFoodHandler(bot *tgbotapi.BotAPI, msg *tgbotapi.Message, u User) {
	args := strings.SplitN(msg.Text, " ", 3)
	if len(args) < 3 {
		reply(bot, msg, "Использование:\n/addfood <ккал> <описание>\nПример: /addfood 250 Яблоко")
		return
	}
	kcal, err := strconv.Atoi(args[1])
	if err != nil {
		reply(bot, msg, "Неверный формат калорий")
		return
	}
	desc := args[2]

	_, err = DB.Exec("INSERT INTO meals (user_id, description, calories) VALUES (?, ?, ?)", u.ID, desc, kcal)
	if err != nil {
		reply(bot, msg, "Ошибка сохранения еды")
		return
	}

	_, _ = DB.Exec("UPDATE users SET calories_today = calories_today + ? WHERE id = ?", kcal, u.ID)

	reply(bot, msg, fmt.Sprintf("Добавлено: %s — %d ккал", desc, kcal))
}

func CheckFoodHandler(bot *tgbotapi.BotAPI, msg *tgbotapi.Message, u User) {
	// показаьть еду за сегодня и оставшийся лимит
	startOfDay := time.Now().Format("2006-01-02") + " 00:00:00"
	rows, err := DB.Query("SELECT description, calories, created_at FROM meals WHERE user_id = ? AND created_at >= ?", u.ID, startOfDay)
	if err != nil {
		reply(bot, msg, "Ошибка чтения базы")
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
		_, _ = DB.Exec("UPDATE users SET calories_goal = ? WHERE id = ?", u.CaloriesGoal, u.ID)
	}

	remaining := u.CaloriesGoal - total
	if remaining < 0 {
		remaining = 0
	}

	text += fmt.Sprintf("\nВсего: %d ккал\nОсталось до дневной нормы (%d): %d ккал", total, u.CaloriesGoal, remaining)
	reply(bot, msg, text)
}

