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
	//lenStep                    = 0.65 // средняя длина шага.
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

func parseTraining(data string) (int, string, time.Duration, error) {
	parts := strings.Split(data, ",")
	if len(parts) != 3 {
		return 0, "", 0, fmt.Errorf("Некорректный формат данных")
	}

	steps, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, "", 0, fmt.Errorf("ошибка преобразования количества шагов в int: %v", err)
	}
	if steps <= 0 {
		return 0, "", 0, fmt.Errorf("Количество шагов не может быть меньше или равно 0")
	}

	activity := parts[1]

	duration, err := time.ParseDuration(parts[2])
	if err != nil {
		return 0, "", 0, fmt.Errorf("ошибка преобразования продолжительности прогулки в time.Duration: %v", err)
	}
	if duration.Minutes() <= 0 {
		return 0, "", 0, fmt.Errorf("время активности не может быть меньше или равно 0")
	}

	return steps, activity, duration, nil
}

func distance(steps int, height float64) float64 {
	if height <= 0 {
		return 0
	}
	return (float64(steps) * height * stepLengthCoefficient) / mInKm
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration.Minutes() <= 0 {
		return 0
	}
	dist := distance(steps, height)
	hours := duration.Hours()
	return dist / hours
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, activity, duration, err := parseTraining(data)
	if err != nil {
		log.Println(err)
		return "", err
	}

	dist := distance(steps, height)
	avgSpeed := meanSpeed(steps, height, duration)
	var calories float64
	switch activity {
	case "Бег":
		calories, err = RunningSpentCalories(steps, weight, height, duration)
	case "Ходьба":
		calories, err = WalkingSpentCalories(steps, weight, height, duration)
	default:
		return "", fmt.Errorf("неизвестный тип тренировки")
	}
	if err != nil {
		log.Println(err)
		return "", err
	}

	return fmt.Sprintf(`Тип тренировки: %s
Длительность: %.2f ч.
Дистанция: %.2f км.
Скорость: %.2f км/ч
Сожгли калорий: %.2f
`, activity, duration.Hours(), dist, avgSpeed, calories), nil

}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	var err error
	if steps <= 0 {
		err = fmt.Errorf("количество шагов не может быть меньше 0. ")
	}
	if weight <= 0 {
		err = errors.Join(err, fmt.Errorf("вес не может быть меньше или равен 0. "))
	}
	if height <= 0 {
		err = errors.Join(err, fmt.Errorf("рост не может быть меньше или равен 0. "))
	}
	if duration.Minutes() <= 0 {
		err = errors.Join(err, fmt.Errorf("время активности не может быть меньше или равно 0. "))
	}
	if err != nil {
		return 0, err
	}

	avgSpeed := meanSpeed(steps, height, duration)

	return (weight * avgSpeed * duration.Minutes()) / minInH, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	var err error
	if steps <= 0 {
		err = fmt.Errorf("количество шагов не может быть меньше 0. ")
	}
	if weight <= 0 {
		err = errors.Join(err, fmt.Errorf("вес не может быть меньше 0. "))
	}
	if height <= 0 {
		err = errors.Join(err, fmt.Errorf("рост не может быть меньше 0. "))
	}
	if duration.Minutes() <= 0 {
		err = errors.Join(err, fmt.Errorf("время активности не может быть меньше или равно 0. "))
	}
	if err != nil {
		return 0, err
	}

	avgSpeed := meanSpeed(steps, height, duration)

	return ((weight * avgSpeed * duration.Minutes()) / minInH) * walkingCaloriesCoefficient, nil
}
