package spentcalories

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep                    = 0.65 // средняя длина шага.
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

func parseTraining(data string) (int, string, time.Duration, error) {
	parts := strings.Split(data, ",")
	if len(parts) != 3 {
		return 0, "", 0, fmt.Errorf("неверный формат данных: ожидается 3 значения, получено %d", len(parts))
	}

	steps, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, "", 0, fmt.Errorf("некорректное количество шагов: %w", err)
	}

	if steps <= 0 {
		return 0, "", 0, errors.New("количество шагов должно быть больше 0")
	}

	duration, err := time.ParseDuration(parts[2])
	if err != nil {
		return 0, "", 0, fmt.Errorf("некорректная продолжительность: %w", err)
	}

	if duration <= 0 {
		return 0, "", 0, errors.New("продолжительность должна быть больше 0")
	}

	return steps, parts[1], duration, nil
}

func distance(steps int, height float64) float64 {
	stepLength := height * stepLengthCoefficient

	return float64(steps) * stepLength / mInKm
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}

	return distance(steps, height) / duration.Hours()
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, activity, duration, err := parseTraining(data)
	if err != nil {
		log.Println(err)
		return "", err
	}

	var calories float64

	switch activity {
	case "Бег":
		calories, err = RunningSpentCalories(steps, weight, height, duration)
	case "Ходьба":
		calories, err = WalkingSpentCalories(steps, weight, height, duration)
	default:
		return "", errors.New("неизвестный тип тренировки")
	}

	if err != nil {
		log.Println(err)
		return "", err
	}

	return fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n",
		activity,
		duration.Hours(),
		distance(steps, height),
		meanSpeed(steps, height, duration),
		calories,
	), nil
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if err := checkParams(steps, weight, height, duration); err != nil {
		return 0, err
	}

	speed := meanSpeed(steps, height, duration)

	return weight * speed * duration.Minutes() / minInH, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if err := checkParams(steps, weight, height, duration); err != nil {
		return 0, err
	}

	speed := meanSpeed(steps, height, duration)
	calories := weight * speed * duration.Minutes() / minInH

	return calories * walkingCaloriesCoefficient, nil
}

func checkParams(steps int, weight, height float64, duration time.Duration) error {
	if steps <= 0 {
		return errors.New("количество шагов должно быть больше 0")
	}

	if weight <= 0 {
		return errors.New("вес должен быть больше 0")
	}

	if height <= 0 {
		return errors.New("рост должен быть больше 0")
	}

	if duration <= 0 {
		return errors.New("продолжительность должна быть больше 0")
	}

	return nil
}
