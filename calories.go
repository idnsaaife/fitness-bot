package main

import (
	"math"
)

// Mifflin-St Jeor formula
func CalculateBMR(weightKg float64, heightCm int, age int, male bool) float64 {
	// BMR = 10W + 6.25H - 5A + 5 (male) / -161 (female)
	
	return 10*weightKg + 6.25*float64(heightCm) - 5*float64(age) - 161
}

func ActivityFactor(level ActivityLevel) float64 {
	switch level {
	case ActivityLow:
		return 1.2
	case ActivityMedium:
		return 1.375
	case ActivityGood:
		return 1.55
	case ActivityHigh:
		return 1.725
	default:
		return 1.375
	}
}

func GoalAdjustment(goal Goal) float64 {
	switch goal {
	case GoalLose:
		return -500.0 // дефицит в ккал/день
	case GoalGain:
		return 300.0 // профицит
	default:
		return 0.0
	}
}

// Рассчитать дневную норму калорий для пользователя
func CalcDailyCalories(u User) int {
	bmr := CalculateBMR(u.WeightKg, u.HeightCm, u.Age, true)
	af := ActivityFactor(u.ActivityLevel)
	goalAdj := GoalAdjustment(u.Goal)
	cal := bmr*af + goalAdj
	return int(math.Round(cal))
}

// Подсчёт сожжённых калорий по активности:
// calories = MET * weight_kg * hours
func CaloriesForActivity(activity string, durationMinutes int, weightKg float64) int {
	met := METForActivity(activity)
	hours := float64(durationMinutes) / 60.0
	cal := met * weightKg * hours
	return int(math.Round(cal))
}

func METForActivity(activity string) float64 {
	switch activity {
	case "бег", "run", "running":
		return 8.0
	case "эллипс", "elliptical":
		return 5.0
	case "велик", "bike", "cycling":
		return 7.5
	case "силовая", "strength", "weights":
		return 6.0
	case "ходьба", "walk", "walking":
		return 3.5
	default:
		return 4.0
	}
}
