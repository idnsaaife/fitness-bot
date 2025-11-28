package main

import "time"

type Goal string
type ActivityLevel string

const (
	GoalLose   Goal = "lose"
	GoalGain   Goal = "gain"
	GoalMaintain Goal = "maintain"

	ActivityLow    ActivityLevel = "low"
	ActivityMedium ActivityLevel = "medium"
	ActivityGood   ActivityLevel = "good"
	ActivityHigh   ActivityLevel = "high"
)

type User struct {
	ID            int
	TgID          int64
	CaloriesGoal  int    // дневная норма
	WaterGoal     int    // мл
	WaterToday    int
	CaloriesToday int

	HeightCm      int
	WeightKg      float64
	Age           int
	Goal          Goal
	ActivityLevel ActivityLevel

	WaterIntervalMinutes int       // 0 = выключено, 60, 120, 240
	RegisteredAt         time.Time
}
