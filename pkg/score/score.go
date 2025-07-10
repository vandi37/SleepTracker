package score

import (
	"time"
)

const (
	MaxTimeIn36HourFormat = 2160 // 36 hours in mins
)

func GetAge(birth time.Time) float64 {
	return time.Since(birth).Hours() / 365.2425 * 24
}

// CalculateSleepScore calculates a score (1–100) in 36 hour format
func CalculateSleepScore(age float64, sleepTime int16, wakeTime int16) int16 {
	minSleep, maxSleep := getSleepDurationByAge(age)
	sleepStartMin, sleepStartMax, wakeUpMin, wakeUpMax := getIdealSleepWindow(age)
	duration := wakeTime - sleepTime
	durationScore := evaluateDurationScore(duration, minSleep, maxSleep)
	sleepStartScore := evaluateSleepStartScore(sleepTime, sleepStartMin, sleepStartMax)
	wakeTimeScore := evaluateWakeTimeScore(wakeTime, wakeUpMin, wakeUpMax)
	totalScore := durationScore + sleepStartScore + wakeTimeScore
	return clamp(totalScore, 1, 100)
}

func getSleepDurationByAge(age float64) (int16, int16) {
	switch {
	case age < 1: // 0-11 months
		return 780, 1020
	case age < 3: // 1–2 years
		return 660, 840 // 11–14
	case age < 6: // 3–5 years
		return 600, 780 // 10–13
	case age < 14: // 6–13 years
		return 540, 660 // 9–11
	case age < 18: // 14–17 years
		return 480, 600 // 8–10
	case age < 65: // 18–64 years
		return 420, 540 // 7–9
	default: // 65+ years
		return 420, 480 // 7–8
	}
}

func getIdealSleepWindow(age float64) (sleepStartMin, sleepStartMax, wakeUpMin, wakeUpMax int16) {
	switch {
	case age < 13: // Kids
		return 480, 570, 1080, 1200 // 20:00–21:30 → 6:00–8:00
	case age < 18: // Teenagers
		return 540, 630, 1140, 1260 // 21:00–22:30 → 7:00–9:00
	case age < 65: // Adults
		return 600, 720, 1110, 1260 // 22:00–00:00 → 6:30–9:00
	default: // Old
		return 540, 660, 1020, 1140 // 21:00–23:00 → 5:00–7:00
	}
}

func evaluateDurationScore(duration, minSleep, maxSleep int16) int16 {
	if duration < minSleep {
		penalty := (minSleep - duration) / 30
		return max(0, 50-6*penalty)
	} else if duration > maxSleep {
		penalty := (duration - maxSleep) / 30
		return max(0, 50-2*penalty)
	}
	return 50
}

func evaluateSleepStartScore(sleepTime, sleepStartMin, sleepStartMax int16) int16 {
	if sleepTime <= sleepStartMin {
		penalty := (sleepStartMin - sleepTime + 29) / 30
		return max(0, 25-penalty)
	} else if sleepTime <= sleepStartMax {
		return 25
	} else {
		penalty := (sleepTime - sleepStartMax + 29) / 30
		return max(0, 25-3*penalty)
	}
}

func evaluateWakeTimeScore(wakeUpMin, wakeUpMax, wakeTime int16) int16 {
	if wakeTime <= wakeUpMin {
		penalty := (wakeUpMin - wakeTime + 29) / 30
		return max(0, 25-3*penalty)
	} else if wakeTime <= wakeUpMax {
		return 25
	} else {
		penalty := (wakeTime - wakeUpMax + 29) / 30
		return max(0, 25-penalty)
	}
}

func clamp(value, minimum, maximum int16) int16 {
	return max(minimum, min(maximum, value))
}
