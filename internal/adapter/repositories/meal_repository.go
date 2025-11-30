package repositories

import (
	"database/sql"
)

type MealRepo struct {
	Db *sql.DB
}

func NewMealRepo(db *sql.DB) *MealRepo {
	return &MealRepo{Db: db}
}

func (mealRepo *MealRepo) SaveFoodWithCalories(id int, desc string, kcal int) error {
	_, err := mealRepo.Db.Exec("INSERT INTO meals (user_id, description, calories) VALUES (?, ?, ?)", id, desc, kcal)
	return err
}

func (mealRepo *MealRepo) GetAllFoodByDay(id int, start string) (*sql.Rows, error) {
	rows, err := mealRepo.Db.Query("SELECT description, calories, created_at FROM meals WHERE user_id = ? AND created_at >= ?", id, start)
	if err != nil {
		return nil, err
	}
	return rows, nil
}
