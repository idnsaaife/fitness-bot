package repositories

import (
	"database/sql"
	"fitness-bot/internal/domain"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (uRepo *UserRepo) GetQueryWaterReminders(bot *tgbotapi.BotAPI) (*sql.Rows, error) {
	rows, err := uRepo.db.Query("SELECT tg_id, water_interval_minutes FROM users WHERE water_interval_minutes > 0")
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (uRepo *UserRepo) UpdateTodayCalories(calories, uID int) {
	_, _ = uRepo.db.Exec("UPDATE users SET calories_today = calories_today - ? WHERE id = ?", calories, uID)

}

func (uRepo *UserRepo) CreateUser(tgID int64, heightCm int, weightKg float64, age int, goal domain.Goal, activity domain.ActivityLevel) (domain.User, error) {
	res, err := uRepo.db.Exec(
		`INSERT OR IGNORE INTO users (tg_id, height_cm, weight_kg, age, goal, activity_level, calories_goal, registered_at)
         VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		tgID, heightCm, weightKg, age, string(goal), string(activity), 0, time.Now(),
	)
	if err != nil {
		return domain.User{}, err
	}
	_ = res

	return uRepo.GetUserByTelegramID(tgID)
}

func (uRepo *UserRepo) GetUserByTelegramID(tgID int64) (domain.User, error) {
	var u domain.User
	uId := u.GetId()
	uTgId := u.GetTgID()
	uCaloriesGoal := u.GetCaloriesGoal()
	uWaterGoal := u.GetWaterGoal()
	uWaterToday := u.GetWaterToday()
	uCaloriesToday := u.GetCaloriesToday()
	uHeightCm := u.GetHeightCm()
	uWeightKg := u.GetWeightKg()
	uAge := u.GetAge()
	uGoal := u.GetGoal()
	uActivityLevel := u.GetActivityLevel()
	uWaterIntervalMinutes := u.GetWaterIntervalMinutes()
	//uRegisteredAt := u.GetRegisteredAt()

	var regAt sql.NullString

	err := uRepo.db.QueryRow(`
		SELECT id, tg_id, calories_goal, water_goal, water_today, calories_today, 
		       height_cm, weight_kg, age, goal, activity_level, water_interval_minutes, registered_at 
		FROM users WHERE tg_id = ?`, tgID).
		Scan(uId, uTgId, uCaloriesGoal, uWaterGoal, uWaterToday, uCaloriesToday,
			uHeightCm, uWeightKg, uAge, uGoal, uActivityLevel, uWaterIntervalMinutes, &regAt)

	if err != nil {
		if err == sql.ErrNoRows {

			return domain.User{}, nil
		}
		return domain.User{}, err
	}

	if regAt.Valid {
		_, _ = time.Parse("2006-01-02 15:04:05", regAt.String)
	}

	return u, nil
}

func (uRepo *UserRepo) UpdateWaterIntervalMinutes(id, minutes int) {
	_, _ = uRepo.db.Exec("UPDATE users SET water_interval_minutes = ? WHERE id = ?", minutes, id)
}

func (uRepo *UserRepo) UpdateWaterToday(ml, id int) {
	_, _ = uRepo.db.Exec("UPDATE users SET water_today = water_today + ? WHERE id = ?", ml, id)
}

func (uRepo *UserRepo) UpdateWaterIntervalMinutesOff(id int) {
	_, _ = uRepo.db.Exec("UPDATE users SET water_interval_minutes = 0 WHERE id = ?", id)
}

func (uRepo *UserRepo) UpdateUserParams(height int, weight float64, age int, goal, act string, id int) error {
	_, err := uRepo.db.Exec("UPDATE users SET height_cm = ?, weight_kg = ?, age = ?, goal = ?, activity_level = ? WHERE id = ?",
		height, weight, age, goal, act, id)
	return err
}

func (uRepo *UserRepo) UpdateGoalCalories(newCal, id int) {
	_, _ = uRepo.db.Exec("UPDATE users SET calories_goal = ? WHERE id = ?", newCal, id)
}

func (uRepo *UserRepo) UpdateCalories(kcal, id int) {
	_, _ = uRepo.db.Exec("UPDATE users SET calories_today = calories_today + ? WHERE id = ?", kcal, id)

}
