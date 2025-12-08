package repositories

import (
	"database/sql"
)

type MealRepo struct {
	db *sql.DB
}

func NewMealRepo(db *sql.DB) *MealRepo {
	return &MealRepo{db: db}
}

func (mealRepo *MealRepo) SaveFoodWithCalories(id int, desc string, kcal int) error {
	_, err := mealRepo.db.Exec("INSERT INTO meals (user_id, description, calories) VALUES (?, ?, ?)", id, desc, kcal)
	return err
}

func (mealRepo *MealRepo) GetAllFoodByDay(id int, start string) (*sql.Rows, error) {
	rows, err := mealRepo.db.Query("SELECT description, calories, created_at FROM meals WHERE user_id = ? AND created_at >= ?", id, start)
	if err != nil {
		return nil, err
	}
	return rows, nil
}
