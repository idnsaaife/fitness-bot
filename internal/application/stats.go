package application

import (
	"fitness-bot/internal/adapter/repositories"
	"fitness-bot/internal/domain"
	"fmt"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (appHandler *AppHandler) StatsHandler(bot *tgbotapi.BotAPI, msg *tgbotapi.Message, u domain.User,
	weightRepo *repositories.WeightChangesRepo, actRepo *repositories.ActivityRepo) {
	var firstWeight float64
	var lastWeight float64

	row := weightRepo.SelectWeightAsc(*u.GetId())
	row.Scan(&firstWeight)

	row = weightRepo.SelectWeightDesc(*u.GetId())
	row.Scan(&lastWeight)

	monthAgo := time.Now().AddDate(0, -1, 0).Format("2006-01-02 15:04:05")
	row2 := actRepo.CalculateCountActivitiesFromMonth(*u.GetId(), monthAgo)
	var count int
	if row2 != nil {
		row2.Next()
		row2.Scan(&count)
		row2.Close()
	}

	text := "📊 Статистика:\n"
	if firstWeight == 0 {
		text += fmt.Sprintf("Вес: сейчас %.1f кг\n", *u.GetWeightKg())
	} else {
		text += fmt.Sprintf("Вес: %.1f кг (первый) → %.1f кг (последний)\n", firstWeight, lastWeight)
	}
	text += fmt.Sprintf("Тренировок за последний месяц: %d\n", count)

	appHandler.Reply(bot, msg, text)
}
