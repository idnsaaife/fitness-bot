package main

import (
	"fmt"
	"strconv"
	"time"
	"strings"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func WaterCommandHandler(bot *tgbotapi.BotAPI, msg *tgbotapi.Message, u User) {
	args := strings.Fields(msg.Text)
	if len(args) < 2 {
		reply(bot, msg, "Использование:\n/water off\n/water 1\n/water 2\n/water 4\n(устанавливает частоту напоминаний в часах)")
		return
	}
	if args[1] == "off" {
		_, _ = DB.Exec("UPDATE users SET water_interval_minutes = 0 WHERE id = ?", u.ID)
		reply(bot, msg, "Напоминания о воде отключены.")
		return
	}
	hours, err := strconv.Atoi(args[1])
	if err != nil || !(hours == 1 || hours == 2 || hours == 4) {
		reply(bot, msg, "Неверный аргумент. Разрешены: 1,2,4 или off")
		return
	}
	mins := hours * 60
	_, _ = DB.Exec("UPDATE users SET water_interval_minutes = ? WHERE id = ?", mins, u.ID)
	reply(bot, msg, fmt.Sprintf("Напоминания установлены каждые %d часов.", hours))

	StartWaterReminderForUser(bot, u.TgID, mins)
}

var waterReminders = map[int64]chan bool{} // tgID -> stop channel

func StartWaterRemindersOnBoot() {
	rows, err := DB.Query("SELECT tg_id, water_interval_minutes FROM users WHERE water_interval_minutes > 0")
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var tgID int64
		var mins int
		rows.Scan(&tgID, &mins)
		_ = tgID
		_ = mins
	}
}

func StartWaterReminders(bot *tgbotapi.BotAPI) {
	rows, err := DB.Query("SELECT tg_id, water_interval_minutes FROM users WHERE water_interval_minutes > 0")
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var tgID int64
		var mins int
		rows.Scan(&tgID, &mins)
		StartWaterReminderForUser(bot, tgID, mins)
	}
}

func StartWaterReminderForUser(bot *tgbotapi.BotAPI, tgID int64, mins int) {
	if ch, ok := waterReminders[tgID]; ok {
		ch <- true
		delete(waterReminders, tgID)
	}

	stop := make(chan bool)
	waterReminders[tgID] = stop

	go func() {
		ticker := time.NewTicker(time.Duration(mins) * time.Minute)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				msg := tgbotapi.NewMessage(tgID, "⏰ Пора выпить воды! 💧 Отметь сколько мл с помощью /water 250")
				bot.Send(msg)
			case <-stop:
				return
			}
		}
	}()
}
