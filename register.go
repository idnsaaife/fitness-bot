package main

import (
	"fmt"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type RegState string

const (
	RegNone       RegState = ""
	RegHeight     RegState = "height"
	RegWeight     RegState = "weight"
	RegAge        RegState = "age"
	RegGoal       RegState = "goal"
	RegActivity   RegState = "activity"
	RegCompleted  RegState = "done"
)

var regStates = map[int64]RegState{}
var regData = map[int64]map[string]string{}

func StartHandler(bot *tgbotapi.BotAPI, msg *tgbotapi.Message) {

	u, error := GetUserByTelegramID(msg.From.ID)

	if (error != nil) {
		str := error.Error()
		reply(bot, msg, str)
		return
	}
	if u.ID != 0 {
		reply(bot, msg, "Вы уже зарегистрированы! Для изменения данных используйте /edit")
		return
	}

	regStates[msg.Chat.ID] = RegHeight
	regData[msg.Chat.ID] = map[string]string{}

	reply(bot, msg, "Добро пожаловать! 🎉\nНачнём регистрацию.\nВведите ваш рост в сантиметрах:")
}

func HandleRegistration(bot *tgbotapi.BotAPI, msg *tgbotapi.Message) bool {
	state, ok := regStates[msg.Chat.ID]
	if !ok || state == RegNone {
		return false
	}

	text := strings.TrimSpace(msg.Text)

	switch state {

	case RegHeight:
		val, err := strconv.Atoi(text)
		if err != nil || val < 120 || val > 250 {
			reply(bot, msg, "Введите корректный рост в сантиметрах, например 170:")
			return true
		}
		regData[msg.Chat.ID]["height"] = text
		regStates[msg.Chat.ID] = RegWeight
		reply(bot, msg, "Введите ваш вес (в кг), например 65.5:")
		return true

	case RegWeight:
		val, err := strconv.ParseFloat(text, 64)
		if err != nil || val < 30 || val > 300 {
			reply(bot, msg, "Введите корректный вес в кг, например 65.5:")
			return true
		}
		regData[msg.Chat.ID]["weight"] = text
		regStates[msg.Chat.ID] = RegAge
		reply(bot, msg, "Введите ваш возраст:")
		return true

	case RegAge:
		val, err := strconv.Atoi(text)
		if err != nil || val < 10 || val > 100 {
			reply(bot, msg, "Введите корректный возраст, например 25:")
			return true
		}
		regData[msg.Chat.ID]["age"] = text
		regStates[msg.Chat.ID] = RegGoal


		msgOut := tgbotapi.NewMessage(msg.Chat.ID, "Выберите вашу цель:")
		msgOut.ReplyMarkup = goalButtons()
		bot.Send(msgOut)
		return true

	case RegGoal:
		// пользователь должен выбрать из inline кнопки — текст сюда не дойдёт
		return true

	case RegActivity:
		// тоже inline кнопки
		return true
	}

	return false
}

func goalButtons() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Похудеть", "goal:lose"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Набрать массу", "goal:gain"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Поддерживать вес", "goal:maintain"),
		),
	)
}

func activityButtons() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Низкая", "activity:low"),
			tgbotapi.NewInlineKeyboardButtonData("Средняя", "activity:medium"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Приемлемая", "activity:good"),
			tgbotapi.NewInlineKeyboardButtonData("Высокая", "activity:high"),
		),
	)
}

func HandleCallback(bot *tgbotapi.BotAPI, cb *tgbotapi.CallbackQuery) {
	data := cb.Data
	tgID := cb.Message.Chat.ID
	state := regStates[tgID]

	if strings.HasPrefix(data, "goal:") && state == RegGoal {
		goal := strings.TrimPrefix(data, "goal:")
		regData[tgID]["goal"] = goal
		regStates[tgID] = RegActivity

		msg := tgbotapi.NewMessage(tgID, "Выберите вашу активность:")
		msg.ReplyMarkup = activityButtons()
		bot.Send(msg)

		ack(bot, cb)
		return
	}

	if strings.HasPrefix(data, "activity:") && state == RegActivity {
		activity := strings.TrimPrefix(data, "activity:")
		regData[tgID]["activity"] = activity
		regStates[tgID] = RegCompleted

		ack(bot, cb)

		// теперь регистрируем пользователя в БД
		FinalizeRegistration(bot, tgID)
		return
	}

	ack(bot, cb)
}

func FinalizeRegistration(bot *tgbotapi.BotAPI, tgID int64) {
	d := regData[tgID]

	height, _ := strconv.Atoi(d["height"])
	weight, _ := strconv.ParseFloat(d["weight"], 64)
	age, _ := strconv.Atoi(d["age"])

	var goal Goal
	switch d["goal"] {
	case "lose":
		goal = GoalLose
	case "gain":
		goal = GoalGain
	default:
		goal = GoalMaintain
	}

	var act ActivityLevel
	switch d["activity"] {
	case "low":
		act = ActivityLow
	case "medium":
		act = ActivityMedium
	case "good":
		act = ActivityGood
	case "high":
		act = ActivityHigh
	}

	u, err := CreateUser(tgID, height, weight, age, goal, act)
	if err != nil {
		send(bot, tgID, "Произошла ошибка при регистрации. Попробуйте снова.")
		return
	}

	// считаем норму
	cal := CalcDailyCalories(u)
	DB.Exec("UPDATE users SET calories_goal = ? WHERE id = ?", cal, u.ID)

	send(bot, tgID, fmt.Sprintf(
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

	// зачистка
	delete(regStates, tgID)
	delete(regData, tgID)
}

func send(bot *tgbotapi.BotAPI, chatID int64, text string) {
	m := tgbotapi.NewMessage(chatID, text)
	m.ParseMode = "Markdown"
	bot.Send(m)
}

func ack(bot *tgbotapi.BotAPI, cb *tgbotapi.CallbackQuery) {
	bot.Request(tgbotapi.NewCallback(cb.ID, ""))
}
