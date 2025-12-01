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

type CallbackHandler struct {
	bot *tgbotapi.BotAPI
}

func NewCallbackHandler(Bot *tgbotapi.BotAPI) *CallbackHandler {
	return &CallbackHandler{bot: Bot}
}

var activityStates = map[int64]string{}

func (callbackHandler *CallbackHandler) HandleCallback(bot *tgbotapi.BotAPI, cb *tgbotapi.CallbackQuery,
	uRepo *repositories.UserRepo, waterHandler *WaterHandler, actHandler *ActHandler, appHandler *AppHandler) {
	data := cb.Data
	chatID := cb.Message.Chat.ID
	userID := cb.From.ID

	log.Printf("Callback received: %s from user %d", data, userID)

	if callbackHandler.handleRegistrationCallbacks(callbackHandler.bot, cb, uRepo, actHandler, appHandler) {
		return
	}

	if strings.HasPrefix(data, "water:") {
		callbackHandler.HandleWaterCallback(callbackHandler.bot, cb, uRepo, waterHandler)
		return
	}

	if strings.HasPrefix(data, "activity:") {
		callbackHandler.HandleActivityCallback(callbackHandler.bot, cb, uRepo)
		return
	}

	if strings.HasPrefix(data, "food:") {
		return
	}

	callbackHandler.Ack(callbackHandler.bot, cb)
	callbackHandler.Send(callbackHandler.bot, chatID, "Неизвестная команда")
}

func (callbackHandler *CallbackHandler) handleRegistrationCallbacks(bot *tgbotapi.BotAPI, cb *tgbotapi.CallbackQuery,
	userRepo *repositories.UserRepo, actHandler *ActHandler, appHandler *AppHandler) bool {
	data := cb.Data
	chatID := cb.Message.Chat.ID
	state := regStates[chatID]

	if strings.HasPrefix(data, "goal:") && state == RegGoal {
		keyboardHandler := adapter.NewKeyboardHandler()
		goal := strings.TrimPrefix(data, "goal:")
		regData[chatID]["goal"] = goal
		regStates[chatID] = RegActivity

		msg := tgbotapi.NewMessage(chatID, "Выберите вашу активность:")
		msg.ReplyMarkup = keyboardHandler.ActivityButtons()
		bot.Send(msg)

		callbackHandler.Ack(bot, cb)
		return true
	}

	if strings.HasPrefix(data, "activity:") && state == RegActivity {
		activity := strings.TrimPrefix(data, "activity:")
		regData[chatID]["activity"] = activity
		regStates[chatID] = RegCompleted

		callbackHandler.Ack(bot, cb)
		appHandler.FinalizeRegistration(bot, chatID, userRepo, actHandler, callbackHandler)
		return true
	}

	return false
}

func (callbackHandler *CallbackHandler) HandleWaterCallback(bot *tgbotapi.BotAPI, cb *tgbotapi.CallbackQuery,
	uRepo *repositories.UserRepo, waterHandler *WaterHandler) {
	data := strings.TrimPrefix(cb.Data, "water:")
	chatID := cb.Message.Chat.ID

	user, err := uRepo.GetUserByTelegramID(cb.From.ID)
	if err != nil || *user.GetId() == 0 {
		callbackHandler.Send(bot, chatID, "Сначала зарегистрируйтесь через /start")
		callbackHandler.Ack(bot, cb)
		return
	}

	switch data {
	case "250", "500":
		ml, _ := strconv.Atoi(data)
		uRepo.UpdateWaterToday(ml, *user.GetId())
		//_, _ = uRepo.Db.Exec("UPDATE users SET water_today = water_today + ? WHERE id = ?", ml, user.ID)
		callbackHandler.Send(bot, chatID, fmt.Sprintf("💧 Добавлено %d мл воды сегодня!", ml))

	case "60", "120", "240":
		mins, _ := strconv.Atoi(data)
		uRepo.UpdateWaterIntervalMinutes(mins, *user.GetId())
		//_, _ = uRepo.Db.Exec("UPDATE users SET water_interval_minutes = ? WHERE id = ?", mins, user.ID)
		waterHandler.StartWaterReminderForUser(bot, *user.GetTgID(), mins)
		callbackHandler.Send(bot, chatID, fmt.Sprintf("⏰ Напоминания установлены каждые %d минут", mins))

	case "off":
		uRepo.UpdateWaterIntervalMinutesOff(*user.GetId())
		//_, _ = uRepo.Db.Exec("UPDATE users SET water_interval_minutes = 0 WHERE id = ?", user.ID)
		callbackHandler.Send(bot, chatID, "🔕 Напоминания о воде отключены")
	}

	callbackHandler.Ack(bot, cb)
}

func (callbackHandler *CallbackHandler) HandleActivityCallback(bot *tgbotapi.BotAPI, cb *tgbotapi.CallbackQuery, uRepo *repositories.UserRepo) {
	data := strings.TrimPrefix(cb.Data, "activity:")
	chatID := cb.Message.Chat.ID

	user, err := uRepo.GetUserByTelegramID(cb.From.ID)
	if err != nil || *user.GetId() == 0 {
		callbackHandler.Send(bot, chatID, "Сначала зарегистрируйтесь через /start")
		callbackHandler.Ack(bot, cb)
		return
	}

	activityStates[chatID] = data
	callbackHandler.Send(bot, chatID, fmt.Sprintf("Выбрана активность: %s\nТеперь введите длительность в минутах:", data))
	callbackHandler.Ack(bot, cb)
}

func (callbackHandler *CallbackHandler) Send(bot *tgbotapi.BotAPI, chatID int64, text string) {
	m := tgbotapi.NewMessage(chatID, text)
	m.ParseMode = "Markdown"
	bot.Send(m)
}

func (callbackHandler *CallbackHandler) Ack(bot *tgbotapi.BotAPI, cb *tgbotapi.CallbackQuery) {
	bot.Request(tgbotapi.NewCallback(cb.ID, ""))
}
