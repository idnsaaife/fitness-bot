package repositories

import "database/sql"

type WaterNotificationRepo struct {
	Db *sql.DB
}

func NewWaterNotificationRepo(db *sql.DB) *WaterNotificationRepo {
	return &WaterNotificationRepo{Db: db}
}
