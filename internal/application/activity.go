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

type ActivityQHandler struct{}

func ActivityHandler(bot *tgbotapi.BotAPI, msg *tgbotapi.Message, u domain.User) {
	text := `🏃 *Добавление активности*

Выберите тип активности:`

	msgOut := tgbotapi.NewMessage(msg.Chat.ID, text)
	msgOut.ReplyMarkup = adapter.ActivityInlineKeyboard()
	msgOut.ParseMode = "Markdown"
	bot.Send(msgOut)
}

// func no usages
func AddActivityHandler(bot *tgbotapi.BotAPI, msg *tgbotapi.Message, u domain.User) {

	args := strings.SplitN(msg.Text, " ", 3)
	if len(args) < 3 {
		Reply(bot, msg, "Использование:\n/addactivity <мин> <тип>\nТипы: бег, эллипс, велик, силовая, ходьба\nПример: /addactivity 30 бег")
		return
	}
	mins, err := strconv.Atoi(args[1])
	if err != nil {
		Reply(bot, msg, "Неверный формат минут")
		return
	}
	atype := strings.ToLower(strings.TrimSpace(args[2]))
	cal := CaloriesForActivity(atype, mins, u.WeightKg)

	//_, err = adapter.DB.Exec("INSERT INTO activities (user_id, atype, duration_min, calories) VALUES (?, ?, ?, ?)", u.ID, atype, mins, cal)
	if err != nil {
		Reply(bot, msg, "Ошибка сохранения активности")
		return
	}

	//_, _ = adapter.DB.Exec("UPDATE users SET calories_today = calories_today - ? WHERE id = ?", cal, u.ID)
	Reply(bot, msg, fmt.Sprintf("Занятие: %s, %d минут — ~%d ккал сожжено", atype, mins, cal))
}

func IsAddingActivity(chatID int64) bool {
	_, exists := activityStates[chatID]
	return exists
}

func HandleActivityDuration(bot *tgbotapi.BotAPI, msg *tgbotapi.Message, u domain.User, actRepo *repositories.ActivityRepo, userRepo *repositories.UserRepo) {
	chatID := msg.Chat.ID
	activityType := activityStates[chatID]

	duration, err := strconv.Atoi(msg.Text)
	if err != nil || duration <= 0 {
		Reply(bot, msg, "Введите корректную длительность в минутах:")
		return
	}

	calories := CaloriesForActivity(activityType, duration, u.WeightKg)

	err = actRepo.InsertActivityInBase(u.ID, activityType, duration, calories)
	//Create func for db actions
	//_, err = adapter.DB.Exec("INSERT INTO activities (user_id, atype, duration_min, calories) VALUES (?, ?, ?, ?)",
	//	u.ID, activityType, duration, calories)
	if err != nil {
		Reply(bot, msg, "Ошибка сохранения активности")
		delete(activityStates, chatID)
		return
	}

	userRepo.UpdateTodayCalories(calories, u.ID)

	//_, _ = adapter.DB.Exec("UPDATE users SET calories_today = calories_today - ? WHERE id = ?", calories, u.ID)

	Reply(bot, msg, fmt.Sprintf("✅ Активность: %s, %d минут — ~%d ккал сожжено", activityType, duration, calories))
	delete(activityStates, chatID)
	ShowMainMenu(bot, msg, u)
}
