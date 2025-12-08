package application

import (
	"fitness-bot/internal/adapter"
	"fitness-bot/internal/adapter/repositories"
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

func (waterHandler *WaterHandler) HandlerWater(msg *tgbotapi.Message) {
	text := `💧 *Управление водой*

• Нажмите на кнопку чтобы добавить воду
• Установите напоминания`

	keyboardHandler := adapter.NewKeyboardHandler()
	msgOut := tgbotapi.NewMessage(msg.Chat.ID, text)
	msgOut.ReplyMarkup = keyboardHandler.WaterInlineKeyboard()
	msgOut.ParseMode = "Markdown"
	waterHandler.bot.Send(msgOut)
}

var waterReminders = map[int64]chan bool{}

func (waterHandler *WaterHandler) StartWaterReminders(uRepo *repositories.UserRepo) {
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
