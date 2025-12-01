package adapter

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	sql *sql.DB
}

func (db *DB) GetDB() *sql.DB {
	return db.sql
}

//var DB *sql.DB

func InitDB() (*DB, error) {
	var err error
	db, err := sql.Open("sqlite3", "fitness.db")
	if err != nil {
		log.Fatal(err)
	}

	// enable foreign keys
	db.Exec("PRAGMA foreign_keys = ON;")

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		tg_id INTEGER UNIQUE,
		calories_goal INTEGER DEFAULT 2000,
		water_goal INTEGER DEFAULT 2000,
		water_today INTEGER DEFAULT 0,
		calories_today INTEGER DEFAULT 0,
		height_cm INTEGER DEFAULT 170,
		weight_kg REAL DEFAULT 70,
		age INTEGER DEFAULT 30,
		goal TEXT DEFAULT 'maintain',
		activity_level TEXT DEFAULT 'medium',
		water_interval_minutes INTEGER DEFAULT 0,
		registered_at DATETIME DEFAULT (datetime('now'))
	);

	CREATE TABLE IF NOT EXISTS meals (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER,
		description TEXT,
		calories INTEGER,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS activities (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER,
		atype TEXT,
		duration_min INTEGER,
		calories INTEGER,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS weight_logs (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER,
		weight REAL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
	);
	`)

	if err != nil {
		log.Fatal(err)
	}

	wrapper := &DB{sql: db}
	go func() {
		wrapper.StartWaterRemindersOnBoot()
	}()
	return wrapper, nil
}

//func (db *DB) CreateUser(tgID int64, heightCm int, weightKg float64, age int, goal domain.Goal, activity domain.ActivityLevel) (domain.User, error) {
//	res, err := db.SQL.Exec(
//		`INSERT OR IGNORE INTO users (tg_id, height_cm, weight_kg, age, goal, activity_level, calories_goal, registered_at)
//         VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
//		tgID, heightCm, weightKg, age, string(goal), string(activity), 0, time.Now(),
//	)
//	if err != nil {
//		return domain.User{}, err
//	}
//	_ = res
//
//	return db.GetUserByTelegramID(tgID)
//}
//
//func (db *DB) GetUserByTelegramID(tgID int64) (domain.User, error) {
//	var u domain.User
//	var regAt sql.NullString
//
//	err := db.SQL.QueryRow(`
//		SELECT id, tg_id, calories_goal, water_goal, water_today, calories_today,
//		       height_cm, weight_kg, age, goal, activity_level, water_interval_minutes, registered_at
//		FROM users WHERE tg_id = ?`, tgID).
//		Scan(&u.ID, &u.TgID, &u.CaloriesGoal, &u.WaterGoal, &u.WaterToday, &u.CaloriesToday,
//			&u.HeightCm, &u.WeightKg, &u.Age, &u.Goal, &u.ActivityLevel, &u.WaterIntervalMinutes, &regAt)
//
//	if err != nil {
//		if err == sql.ErrNoRows {
//
//			return domain.User{}, nil
//		}
//		return domain.User{}, err
//	}
//
//	if regAt.Valid {
//		u.RegisteredAt, _ = time.Parse("2006-01-02 15:04:05", regAt.String)
//	}
//
//	return u, nil
//}

func (db *DB) StartWaterRemindersOnBoot() {
	rows, err := db.sql.Query("SELECT tg_id, water_interval_minutes FROM users WHERE water_interval_minutes > 0")
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var tgID int64
		var mins int
		rows.Scan(&tgID, &mins)
		_ = tgID
		_ = mins
	}
}
