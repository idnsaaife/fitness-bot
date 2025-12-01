package application

import (
	"fitness-bot/internal/adapter"
	"fitness-bot/internal/adapter/repositories"
	"fitness-bot/internal/domain"
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type WaterHandler struct {
	bot *tgbotapi.BotAPI
}

func NewWaterHandler(Bot *tgbotapi.BotAPI) *WaterHandler {
	return &WaterHandler{bot: Bot}
}

func (waterHandler *WaterHandler) HandlerWater(bot *tgbotapi.BotAPI, msg *tgbotapi.Message, u domain.User) {
	text := `💧 *Управление водой*

• Нажмите на кнопку чтобы добавить воду
• Установите напоминания`

	keyboardHandler := adapter.NewKeyboardHandler()
	msgOut := tgbotapi.NewMessage(msg.Chat.ID, text)
	msgOut.ReplyMarkup = keyboardHandler.WaterInlineKeyboard()
	msgOut.ParseMode = "Markdown"
	waterHandler.bot.Send(msgOut)
}

// func no usage
//func WaterCommandHandler(bot *tgbotapi.BotAPI, msg *tgbotapi.Message, u domain.User, args []string) {
//	if len(args) < 2 {
//		HandlerWater(bot, msg, u)
//		return
//	}
//	if args[1] == "off" {
//		//_, _ = adapter.DB.Exec("UPDATE users SET water_interval_minutes = 0 WHERE id = ?", u.ID)
//		Reply(bot, msg, "Напоминания о воде отключены.")
//		return
//	}
//	hours, err := strconv.Atoi(args[1])
//	if err != nil || !(hours == 1 || hours == 2 || hours == 4) {
//		Reply(bot, msg, "Неверный аргумент. Разрешены: 1,2,4 или off")
//		return
//	}
//	mins := hours * 60
//	//_, _ = adapter.DB.Exec("UPDATE users SET water_interval_minutes = ? WHERE id = ?", mins, u.ID)
//	Reply(bot, msg, fmt.Sprintf("Напоминания установлены каждые %d часов.", hours))
//
//	StartWaterReminderForUser(bot, u.TgID, mins)
//}

var waterReminders = map[int64]chan bool{}

func (waterHandler *WaterHandler) StartWaterReminders(bot *tgbotapi.BotAPI, uRepo *repositories.UserRepo) {
	rows, err := uRepo.GetQueryWaterReminders(waterHandler.bot)
	if err != nil {
		log.Println(err)
	}
	//если что убрать дефер
	defer rows.Close()

	for rows.Next() {
		var tgID int64
		var mins int
		rows.Scan(&tgID, &mins)
		waterHandler.StartWaterReminderForUser(waterHandler.bot, tgID, mins)
	}
}

func (waterHandler *WaterHandler) StartWaterReminderForUser(bot *tgbotapi.BotAPI, tgID int64, mins int) {

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
