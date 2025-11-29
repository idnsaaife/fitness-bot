package main

import (
	"fmt"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func AddActivityHandler(bot *tgbotapi.BotAPI, msg *tgbotapi.Message, u User) {

	args := strings.SplitN(msg.Text, " ", 3)
	if len(args) < 3 {
		reply(bot, msg, "Использование:\n/addactivity <мин> <тип>\nТипы: бег, эллипс, велик, силовая, ходьба\nПример: /addactivity 30 бег")
		return
	}
	mins, err := strconv.Atoi(args[1])
	if err != nil {
		reply(bot, msg, "Неверный формат минут")
		return
	}
	atype := strings.ToLower(strings.TrimSpace(args[2]))
	cal := CaloriesForActivity(atype, mins, u.WeightKg)

	_, err = DB.Exec("INSERT INTO activities (user_id, atype, duration_min, calories) VALUES (?, ?, ?, ?)", u.ID, atype, mins, cal)
	if err != nil {
		reply(bot, msg, "Ошибка сохранения активности")
		return
	}

	_, _ = DB.Exec("UPDATE users SET calories_today = calories_today - ? WHERE id = ?", cal, u.ID)
	reply(bot, msg, fmt.Sprintf("Занятие: %s, %d минут — ~%d ккал сожжено", atype, mins, cal))
}
