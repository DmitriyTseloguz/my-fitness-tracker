package ftracker

import (
	"fmt"
	"math"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep   = 0.65  // средняя длина шага.
	mInKm     = 1000  // количество метров в километре.
	minInHour = 60    // количество минут в часе.
	kmhInMsec = 0.278 // коэффициент для преобразования км/ч в м/с.
	cmInMeter = 100   // количество сантиметров в метре.
)

// Константы для расчета калорий, расходуемых при беге.
const (
	runningCaloriesMeanSpeedMultiplier = 18.0 // множитель средней скорости.
	runningCaloriesMeanSpeedShift      = 1.79 // среднее количество сжигаемых калорий при беге.
)

// Константы для расчета калорий, расходуемых при ходьбе.
const (
	walkingCaloriesWeightMultiplier = 0.035 // множитель массы тела.
	walkingSpeedHeightMultiplier    = 0.029 // множитель роста.
)

// Константы для расчета калорий, расходуемых при плавании.
const (
	swimmingCaloriesMeanSpeedShift   = 1.1 // среднее количество сжигаемых колорий при плавании относительно скорости.
	swimmingCaloriesWeightMultiplier = 2   // множитель веса при плавании.
)

type TrainingType int

// Перечисление типов тренировок
const (
	Run TrainingType = iota
	Walk
	Swim
)

// AvailableTrainings - доступные виды тренировок
var AvailableTrainings = map[string]TrainingType{
	"Бег":      Run,
	"Ходьба":   Walk,
	"Плавание": Swim,
}

// distance возвращает дистанцию(в километрах), которую преодолел пользователь за время тренировки.
//
// Параметры:
//
// action int — количество совершенных действий (число шагов при ходьбе и беге, либо гребков при плавании).
func distance(action int) float64 {
	return float64(action) * lenStep / mInKm
}

// meanSpeed возвращает значение средней скорости движения во время тренировки.
//
// Параметры:
//
// action int — количество совершенных действий(число шагов при ходьбе и беге, либо гребков при плавании).
// duration float64 — длительность тренировки в часах.
func meanSpeed(action int, duration float64) float64 {
	if duration == 0 {
		return 0
	}

	distance := distance(action)

	return distance / duration
}

// ShowTrainingInfo возвращает строку с информацией о тренировке.
//
// Параметры:
//
// action int — количество совершенных действий(число шагов при ходьбе и беге, либо гребков при плавании).
// trainingType string — вид тренировки(Бег, Ходьба, Плавание).
// duration float64 — длительность тренировки в часах.
func ShowTrainingInfo(action int, trainingType string, duration, weight, height float64, lengthPool, countPool int) string {
	var training, isExist = AvailableTrainings[trainingType]

	if !isExist {
		return "неизвестный тип тренировки"
	}

	var distance = distance(action)
	var speed = meanSpeed(action, duration)
	var calories float64

	switch training {
	case Run:
		calories = RunningSpentCalories(action, weight, duration)
	case Walk:
		calories = WalkingSpentCalories(action, duration, weight, height)
	case Swim:
		speed = swimmingMeanSpeed(lengthPool, countPool, duration)
		calories = SwimmingSpentCalories(lengthPool, countPool, duration, weight)
	}

	var formatMessage = "Тип тренировки: %s\n" +
		"Длительность: %.2f ч.\n" +
		"Дистанция: %.2f км.\n" +
		"Скорость: %.2f км/ч\n" +
		"Сожгли калорий: %.2f\n"

	return fmt.Sprintf(
		formatMessage,
		trainingType, duration, distance, speed, calories,
	)
}

// RunningSpentCalories возвращает количество потраченных колорий при беге.
//
// Параметры:
//
// action int — количество совершенных действий(число шагов при ходьбе и беге, либо гребков при плавании).
// weight float64 — вес пользователя.
// duration float64 — длительность тренировки в часах.
func RunningSpentCalories(action int, weight, duration float64) float64 {
	var meanSpeedInKmH = meanSpeed(action, duration)
	var speedCaloriesRation = runningCaloriesMeanSpeedMultiplier * meanSpeedInKmH

	return (speedCaloriesRation * runningCaloriesMeanSpeedShift) * weight / mInKm * duration * minInHour
}

// WalkingSpentCalories возвращает количество потраченных калорий при ходьбе.
//
// Параметры:
//
// action int — количество совершенных действий(число шагов при ходьбе и беге, либо гребков при плавании).
// duration float64 — длительность тренировки в часах.
// weight float64 — вес пользователя.
// height float64 — рост пользователя.
func WalkingSpentCalories(action int, duration, weight, height float64) float64 {
	var meanSpeedMeterInSecond = meanSpeed(action, duration) * kmhInMsec
	var squaredSpeed = math.Pow(meanSpeedMeterInSecond, 2)
	var heightSpeedRation = squaredSpeed / (height / cmInMeter)
	var weightCaloriesRation = walkingCaloriesWeightMultiplier * weight
	var weightSpeedRation = walkingSpeedHeightMultiplier * weight

	return (weightCaloriesRation + heightSpeedRation*weightSpeedRation) * duration * minInHour
}

// swimmingMeanSpeed возвращает среднюю скорость при плавании.
//
// Параметры:
//
// lengthPool int — длина бассейна в метрах.
// countPool int — сколько раз пользователь переплыл бассейн.
// duration float64 — длительность тренировки в часах.
func swimmingMeanSpeed(lengthPool, countPool int, duration float64) float64 {
	if duration == 0 {
		return 0
	}

	return float64(lengthPool) * float64(countPool) / mInKm / duration
}

// SwimmingSpentCalories возвращает количество потраченных калорий при плавании.
//
// Параметры:
//
// lengthPool int — длина бассейна в метрах.
// countPool int — сколько раз пользователь переплыл бассейн.
// duration float64 — длительность тренировки в часах.
// weight float64 — вес пользователя.
func SwimmingSpentCalories(lengthPool, countPool int, duration, weight float64) float64 {
	var meanSpeedInKmH = swimmingMeanSpeed(lengthPool, countPool, duration)
	var weightCaloriesRation = swimmingCaloriesWeightMultiplier * weight

	return (meanSpeedInKmH + swimmingCaloriesMeanSpeedShift) * weightCaloriesRation * duration
}
