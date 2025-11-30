package application

import (
	"fitness-bot/internal/adapter"
	"fitness-bot/internal/adapter/repositories"
	"fmt"
	"log"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var activityStates = map[int64]string{}

func HandleCallback(bot *tgbotapi.BotAPI, cb *tgbotapi.CallbackQuery, uRepo *repositories.UserRepo) {
	data := cb.Data
	chatID := cb.Message.Chat.ID
	userID := cb.From.ID

	log.Printf("Callback received: %s from user %d", data, userID)

	if handleRegistrationCallbacks(bot, cb, uRepo) {
		return
	}

	if strings.HasPrefix(data, "water:") {
		HandleWaterCallback(bot, cb, uRepo)
		return
	}

	if strings.HasPrefix(data, "activity:") {
		HandleActivityCallback(bot, cb, uRepo)
		return
	}

	if strings.HasPrefix(data, "food:") {
		return
	}

	Ack(bot, cb)
	Send(bot, chatID, "Неизвестная команда")
}

func handleRegistrationCallbacks(bot *tgbotapi.BotAPI, cb *tgbotapi.CallbackQuery,
	userRepo *repositories.UserRepo) bool {
	data := cb.Data
	chatID := cb.Message.Chat.ID
	state := regStates[chatID]

	if strings.HasPrefix(data, "goal:") && state == RegGoal {
		goal := strings.TrimPrefix(data, "goal:")
		regData[chatID]["goal"] = goal
		regStates[chatID] = RegActivity

		msg := tgbotapi.NewMessage(chatID, "Выберите вашу активность:")
		msg.ReplyMarkup = adapter.ActivityButtons()
		bot.Send(msg)

		Ack(bot, cb)
		return true
	}

	if strings.HasPrefix(data, "activity:") && state == RegActivity {
		activity := strings.TrimPrefix(data, "activity:")
		regData[chatID]["activity"] = activity
		regStates[chatID] = RegCompleted

		Ack(bot, cb)
		FinalizeRegistration(bot, chatID, userRepo)
		return true
	}

	return false
}

func HandleWaterCallback(bot *tgbotapi.BotAPI, cb *tgbotapi.CallbackQuery, uRepo *repositories.UserRepo) {
	data := strings.TrimPrefix(cb.Data, "water:")
	chatID := cb.Message.Chat.ID

	user, err := uRepo.GetUserByTelegramID(cb.From.ID)
	if err != nil || user.ID == 0 {
		Send(bot, chatID, "Сначала зарегистрируйтесь через /start")
		Ack(bot, cb)
		return
	}

	switch data {
	case "250", "500":
		ml, _ := strconv.Atoi(data)
		uRepo.UpdateWaterToday(ml, user.ID)
		//_, _ = uRepo.Db.Exec("UPDATE users SET water_today = water_today + ? WHERE id = ?", ml, user.ID)
		Send(bot, chatID, fmt.Sprintf("💧 Добавлено %d мл воды сегодня!", ml))

	case "60", "120", "240":
		mins, _ := strconv.Atoi(data)
		uRepo.UpdateWaterIntervalMinutes(mins, user.ID)
		//_, _ = uRepo.Db.Exec("UPDATE users SET water_interval_minutes = ? WHERE id = ?", mins, user.ID)
		StartWaterReminderForUser(bot, user.TgID, mins)
		Send(bot, chatID, fmt.Sprintf("⏰ Напоминания установлены каждые %d минут", mins))

	case "off":
		uRepo.UpdateWaterIntervalMinutesOff(user.ID)
		//_, _ = uRepo.Db.Exec("UPDATE users SET water_interval_minutes = 0 WHERE id = ?", user.ID)
		Send(bot, chatID, "🔕 Напоминания о воде отключены")
	}

	Ack(bot, cb)
}

func HandleActivityCallback(bot *tgbotapi.BotAPI, cb *tgbotapi.CallbackQuery, uRepo *repositories.UserRepo) {
	data := strings.TrimPrefix(cb.Data, "activity:")
	chatID := cb.Message.Chat.ID

	user, err := uRepo.GetUserByTelegramID(cb.From.ID)
	if err != nil || user.ID == 0 {
		Send(bot, chatID, "Сначала зарегистрируйтесь через /start")
		Ack(bot, cb)
		return
	}

	activityStates[chatID] = data
	Send(bot, chatID, fmt.Sprintf("Выбрана активность: %s\nТеперь введите длительность в минутах:", data))
	Ack(bot, cb)
}

func Send(bot *tgbotapi.BotAPI, chatID int64, text string) {
	m := tgbotapi.NewMessage(chatID, text)
	m.ParseMode = "Markdown"
	bot.Send(m)
}

func Ack(bot *tgbotapi.BotAPI, cb *tgbotapi.CallbackQuery) {
	bot.Request(tgbotapi.NewCallback(cb.ID, ""))
}
