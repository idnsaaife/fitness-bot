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

type RegState string

const (
	RegNone      RegState = ""
	RegHeight    RegState = "height"
	RegWeight    RegState = "weight"
	RegAge       RegState = "age"
	RegGoal      RegState = "goal"
	RegActivity  RegState = "activity"
	RegCompleted RegState = "done"
)

var (
	regStates = map[int64]RegState{}
	regData   = map[int64]map[string]string{}
)

func (appHandler *AppHandler) StartHandler(bot *tgbotapi.BotAPI, msg *tgbotapi.Message, userRepo *repositories.UserRepo) {

	u, err := userRepo.GetUserByTelegramID(msg.From.ID)

	if err != nil {
		str := err.Error()
		appHandler.Reply(bot, msg, str)
		return
	}
	if *u.GetId() != 0 {
		appHandler.Reply(bot, msg, "Вы уже зарегистрированы! Для изменения данных используйте /edit")
		return
	}

	regStates[msg.Chat.ID] = RegHeight
	regData[msg.Chat.ID] = map[string]string{}

	appHandler.Reply(bot, msg, "Добро пожаловать! 🎉\nНачнём регистрацию.\nВведите ваш рост в сантиметрах:")
}

func (appHandler *AppHandler) HandleRegistration(bot *tgbotapi.BotAPI, msg *tgbotapi.Message) bool {
	state, ok := regStates[msg.Chat.ID]
	if !ok || state == RegNone {
		return false
	}

	text := strings.TrimSpace(msg.Text)

	switch state {

	case RegHeight:
		val, err := strconv.Atoi(text)
		if err != nil || val < 120 || val > 250 {
			appHandler.Reply(bot, msg, "Введите корректный рост в сантиметрах, например 170:")
			return true
		}
		regData[msg.Chat.ID]["height"] = text
		regStates[msg.Chat.ID] = RegWeight
		appHandler.Reply(bot, msg, "Введите ваш вес (в кг), например 65.5:")
		return true

	case RegWeight:
		val, err := strconv.ParseFloat(text, 64)
		if err != nil || val < 30 || val > 300 {
			appHandler.Reply(bot, msg, "Введите корректный вес в кг, например 65.5:")
			return true
		}
		regData[msg.Chat.ID]["weight"] = text
		regStates[msg.Chat.ID] = RegAge
		appHandler.Reply(bot, msg, "Введите ваш возраст:")
		return true

	case RegAge:
		keyboardHandler := adapter.NewKeyboardHandler()
		val, err := strconv.Atoi(text)
		if err != nil || val < 10 || val > 100 {
			appHandler.Reply(bot, msg, "Введите корректный возраст, например 25:")
			return true
		}
		regData[msg.Chat.ID]["age"] = text
		regStates[msg.Chat.ID] = RegGoal

		msgOut := tgbotapi.NewMessage(msg.Chat.ID, "Выберите вашу цель:")
		msgOut.ReplyMarkup = keyboardHandler.GoalButtons()
		bot.Send(msgOut)
		return true

	case RegGoal:
		return true

	case RegActivity:
		return true
	}

	return false
}

func (appHandler *AppHandler) FinalizeRegistration(bot *tgbotapi.BotAPI, tgID int64, userRepo *repositories.UserRepo,
	actHandler *ActHandler, callbackHandler *CallbackHandler) {
	d := regData[tgID]

	height, _ := strconv.Atoi(d["height"])
	weight, _ := strconv.ParseFloat(d["weight"], 64)
	age, _ := strconv.Atoi(d["age"])

	var goal domain.Goal
	switch d["goal"] {
	case "lose":
		goal = domain.GoalLose
	case "gain":
		goal = domain.GoalGain
	default:
		goal = domain.GoalMaintain
	}

	var act domain.ActivityLevel
	switch d["activity"] {
	case "low":
		act = domain.ActivityLow
	case "medium":
		act = domain.ActivityMedium
	case "good":
		act = domain.ActivityGood
	case "high":
		act = domain.ActivityHigh
	}

	u, err := userRepo.CreateUser(tgID, height, weight, age, goal, act)
	if err != nil {
		callbackHandler.Send(bot, tgID, "Произошла ошибка при регистрации. Попробуйте снова.")
		return
	}

	cal := actHandler.CalcDailyCalories(u)
	userRepo.UpdateGoalCalories(cal, *u.GetId())

	callbackHandler.Send(bot, tgID, fmt.Sprintf(
		"Регистрация завершена! 🎉\n\n"+
			"Ваши параметры:\n"+
			"• Рост: %d см\n"+
			"• Вес: %.1f кг\n"+
			"• Возраст: %d\n"+
			"• Цель: %s\n"+
			"• Уровень активности: %s\n\n"+
			"Ваша дневная норма: *%d ккал*",
		height, weight, age, d["goal"], d["activity"], cal,
	))

	appHandler.ShowMainMenuAfterRegistration(bot, tgID)

	delete(regStates, tgID)
	delete(regData, tgID)
}

func (appHandler *AppHandler) ShowMainMenuAfterRegistration(bot *tgbotapi.BotAPI, chatID int64) {
	text := `🏠 *Главное меню*

Выберите действие:`

	keyboardHandler := adapter.NewKeyboardHandler()
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyMarkup = keyboardHandler.MainMenuKeyboard()
	msg.ParseMode = "Markdown"
	bot.Send(msg)
}
