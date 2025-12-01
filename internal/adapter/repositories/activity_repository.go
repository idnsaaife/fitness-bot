package repositories

import (
	"database/sql"
)

type ActivityRepo struct {
	db *sql.DB
}

func NewActivityRepo(db *sql.DB) *ActivityRepo {
	return &ActivityRepo{db: db}
}

func (actRepo *ActivityRepo) InsertActivityInBase(id int, activityType string, duration, calories int) error {
	_, err := actRepo.db.Exec("INSERT INTO activities (user_id, atype, duration_min, calories) VALUES (?, ?, ?, ?)",
		id, activityType, duration, calories)
	if err != nil {
		return err
	}
	return nil
}

func (actRepo *ActivityRepo) CalculateCountActivitiesFromMonth(id int, month string) *sql.Rows {
	row2, _ := actRepo.db.Query("SELECT COUNT(*) FROM activities WHERE user_id = ? AND created_at >= ?", id, month)
	return row2
}
