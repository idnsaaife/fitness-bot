package repositories

import (
	"database/sql"
)

type ActivityRepo struct {
	Db *sql.DB
}

func NewActivityRepo(db *sql.DB) *ActivityRepo {
	return &ActivityRepo{Db: db}
}

func (actRepo *ActivityRepo) InsertActivityInBase(id int, activityType string, duration, calories int) error {
	_, err := actRepo.Db.Exec("INSERT INTO activities (user_id, atype, duration_min, calories) VALUES (?, ?, ?, ?)",
		id, activityType, duration, calories)
	if err != nil {
		return err
	}
	return nil
}

func (actRepo *ActivityRepo) CalculateCountActivitiesFromMonth(id int, month string) *sql.Rows {
	row2, _ := actRepo.Db.Query("SELECT COUNT(*) FROM activities WHERE user_id = ? AND created_at >= ?", id, month)
	return row2
}
