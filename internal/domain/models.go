package domain

import "time"

type Goal string
type ActivityLevel string

const (
	GoalLose     Goal = "lose"
	GoalGain     Goal = "gain"
	GoalMaintain Goal = "maintain"

	ActivityLow    ActivityLevel = "low"
	ActivityMedium ActivityLevel = "medium"
	ActivityGood   ActivityLevel = "good"
	ActivityHigh   ActivityLevel = "high"
)

type User struct {
	iD            int
	tgID          int64
	caloriesGoal  int // дневная норма
	waterGoal     int // мл
	waterToday    int
	caloriesToday int

	heightCm      int
	weightKg      float64
	age           int
	goal          Goal
	activityLevel ActivityLevel

	waterIntervalMinutes int // 0 = выключено, 30, 60, 120
	registeredAt         time.Time
}

func (u *User) GetId() *int {
	return &u.iD
}

func (u *User) GetTgID() *int64 {
	return &u.tgID
}

func (u *User) GetCaloriesGoal() *int {
	return &u.caloriesGoal
}

func (u *User) GetWaterGoal() *int {
	return &u.waterGoal
}

func (u *User) GetWaterToday() *int {
	return &u.waterToday
}

func (u *User) GetCaloriesToday() *int {
	return &u.caloriesToday
}

func (u *User) GetHeightCm() *int {
	return &u.heightCm
}

func (u *User) GetWeightKg() *float64 {
	return &u.weightKg
}

func (u *User) GetAge() *int {
	return &u.age
}

func (u *User) GetGoal() *Goal {
	return &u.goal
}

func (u *User) GetActivityLevel() *ActivityLevel {
	return &u.activityLevel
}

func (u *User) GetWaterIntervalMinutes() *int {
	return &u.waterIntervalMinutes
}

func (u *User) SetCaloriesGoal(caloriesGoal int) {
	u.caloriesGoal = caloriesGoal
}

func (u *User) SetHeightCm(heightCm int) {
	u.heightCm = heightCm
}

func (u *User) SetWeightKg(weightKg float64) {
	u.weightKg = weightKg
}

func (u *User) SetAge(age int) {
	u.age = age
}

func (u *User) SetActivityLevel(act ActivityLevel) {
	u.activityLevel = act
}

func (u *User) SetGoal(goal Goal) {
	u.goal = goal
}
