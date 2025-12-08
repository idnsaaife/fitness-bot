package application

import (
	"fitness-bot/internal/adapter"
	"fitness-bot/internal/adapter/repositories"
	"fitness-bot/internal/domain"
	"fmt"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type ActHandler struct {
	bot *tgbotapi.BotAPI
}

func NewActHandler(Bot *tgbotapi.BotAPI) *ActHandler {
	return &ActHandler{bot: Bot}
}

func (actHandler *ActHandler) ActivityHandler(msg *tgbotapi.Message) {
	text := `🏃 *Добавление активности*

Выберите тип активности:`

	keyboardHandler := adapter.NewKeyboardHandler()
	msgOut := tgbotapi.NewMessage(msg.Chat.ID, text)
	msgOut.ReplyMarkup = keyboardHandler.ActivityInlineKeyboard()
	msgOut.ParseMode = "Markdown"
	actHandler.bot.Send(msgOut)
}

func (actHandler *ActHandler) IsAddingActivity(chatID int64) bool {
	_, exists := activityStates[chatID]
	return exists
}

func (actHandler *ActHandler) HandleActivityDuration(msg *tgbotapi.Message, u domain.User, actRepo *repositories.ActivityRepo,
	userRepo *repositories.UserRepo, appHandler *AppHandler) {
	chatID := msg.Chat.ID
	activityType := activityStates[chatID]

	duration, err := strconv.Atoi(msg.Text)
	if err != nil || duration <= 0 {
		appHandler.Reply(actHandler.bot, msg, "Введите корректную длительность в минутах:")
		return
	}

	calories := actHandler.CaloriesForActivity(activityType, duration, *u.GetWeightKg())

	err = actRepo.InsertActivityInBase(*u.GetId(), activityType, duration, calories)
	if err != nil {
		appHandler.Reply(actHandler.bot, msg, "Ошибка сохранения активности")
		delete(activityStates, chatID)
		return
	}

	userRepo.UpdateTodayCalories(calories, *u.GetId())

	appHandler.Reply(actHandler.bot, msg, fmt.Sprintf("✅ Активность: %s, %d минут — ~%d ккал сожжено", activityType, duration, calories))
	delete(activityStates, chatID)
	appHandler.ShowMainMenu(actHandler.bot, msg)
}
