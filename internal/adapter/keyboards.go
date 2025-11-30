package adapter

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func MainMenuKeyboard() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("📊 Статистика"),
			tgbotapi.NewKeyboardButton("🍎 Добавить еду"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("💧 Вода"),
			tgbotapi.NewKeyboardButton("🏃 Активность"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("✏️ Редактировать данные"),
			tgbotapi.NewKeyboardButton("📋 Проверить питание"),
		),
	)
}

// Минимальная клавиатура (только старт)
func StartKeyboard() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("/start"),
		),
	)
}

func WaterInlineKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("💧 250 мл", "water:250"),
			tgbotapi.NewInlineKeyboardButtonData("💧 500 мл", "water:500"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⏰ Напоминания: Выкл", "water:off"),
			tgbotapi.NewInlineKeyboardButtonData("⏰ Напоминания: 1ч", "water:60"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⏰ Напоминания: 2ч", "water:120"),
			tgbotapi.NewInlineKeyboardButtonData("⏰ Напоминания: 4ч", "water:240"),
		),
	)
}

func ActivityInlineKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🏃 Бег", "activity:run"),
			tgbotapi.NewInlineKeyboardButtonData("🚶 Ходьба", "activity:walk"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🚴 Велосипед", "activity:bike"),
			tgbotapi.NewInlineKeyboardButtonData("💪 Силовая", "activity:strength"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🏃‍♂️ Эллипс", "activity:elliptical"),
			tgbotapi.NewInlineKeyboardButtonData("⚡ Другое", "activity:other"),
		),
	)
}

func QuickFoodInlineKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🍎 Яблоко (52 ккал)", "food:apple"),
			tgbotapi.NewInlineKeyboardButtonData("🍌 Банан (89 ккал)", "food:banana"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🥪 Сэндвич (300 ккал)", "food:sandwich"),
			tgbotapi.NewInlineKeyboardButtonData("🍚 Рис (130 ккал)", "food:rice"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("➕ Своя еда", "food:custom"),
		),
	)
}

func GoalButtons() tgbotapi.InlineKeyboardMarkup {
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

func ActivityButtons() tgbotapi.InlineKeyboardMarkup {
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
