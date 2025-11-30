package repositories

import (
	"database/sql"
)

type WeightChangesRepo struct {
	Db *sql.DB
}

func NewWeightChangesRepo(db *sql.DB) *WeightChangesRepo {
	return &WeightChangesRepo{Db: db}
}

func (weightRepo *WeightChangesRepo) SelectWeightAsc(id int) *sql.Row {
	row := weightRepo.Db.QueryRow("SELECT weight FROM weight_logs WHERE user_id = ? ORDER BY created_at ASC LIMIT 1", id)
	return row
}

func (weightRepo *WeightChangesRepo) SelectWeightDesc(id int) *sql.Row {
	row := weightRepo.Db.QueryRow("SELECT weight FROM weight_logs WHERE user_id = ? ORDER BY created_at DESC LIMIT 1", id)
	return row
}
