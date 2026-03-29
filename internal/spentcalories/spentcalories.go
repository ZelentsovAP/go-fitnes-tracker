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
	// Разделяем строку на слайс строк
	slice := strings.Split(data, ",")
	if len(slice) != 3 {
		return 0, "", 0, errors.New("Длина слайса меньше или больше 3")
	}
	// Пркобразует кол-во шагов в тип int
	steps, err := strconv.Atoi(slice[0])
	if err != nil {
		err = fmt.Errorf("ошибка преобразовния строки в целое число : %v", err)
		return 0, "", 0, err
	}
	if steps <= 0 {
		err = errors.New("Кол-во шагов меньше или равно 0")
		return 0, "", 0, err
	}
	// pDuration преобразует из слайса строку в тип time.Duration
	pDuration, err := time.ParseDuration(slice[2])
	if err != nil {
		err = fmt.Errorf("ошибка преобразовния строки в тип time.Duration : %v", err)
		return 0, "", 0, err
	}
	if pDuration <= 0 {
		err = errors.New("Время меньше или равно 0")
		return 0, "", 0, err
	}
	activityType := slice[1]
	return steps, activityType, pDuration, nil
}

func distance(steps int, height float64) float64 {
	stepLength := height * stepLengthCoefficient
	distanc := stepLength * float64(steps)
	distanc /= mInKm
	return distanc
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}
	distanc := distance(steps, height)
	hours := duration.Hours()
	midSpeed := distanc / hours
	return midSpeed
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	// Получаем значения из строки данных
	steps, activityType, pDuration, err := parseTraining(data)
	if err != nil {
		log.Println("Ошибка при парсинге данных:", err)
		return "", err
	}
	if steps <= 0 || weight <= 0 || height <= 0 || pDuration <= 0 {
		err = errors.New("Неверные входные данные")
		log.Println(err)
		return "", err
	}
	// Вычисляем необходимые параметры
	dist := distance(steps, height)
	midSpeed := meanSpeed(steps, height, pDuration)
	switch activityType {
	case "Walking", "Ходьба":
		calories, err := WalkingSpentCalories(steps, weight, height, pDuration)
		if err != nil {
			log.Println("Ошибка при вычислении калорий:", err)
			return "", err
		}
		resultStr := fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n", activityType, pDuration.Hours(), dist, midSpeed, calories)
		return resultStr, nil
	case "Running", "Бег":
		calories, err := RunningSpentCalories(steps, weight, height, pDuration)
		if err != nil {
			log.Println("Ошибка при вычислении калорий:", err)
			return "", err
		}
		resultStr := fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n", activityType, pDuration.Hours(), dist, midSpeed, calories)
		return resultStr, nil
	default:
		return "", errors.New("неизвестный тип тренировки")
	}
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if weight <= 0 || height <= 0 || steps <= 0 || duration <= 0 {
		return 0, errors.New("некоректные входные параметры")
	}
	midSpeed := meanSpeed(steps, height, duration)
	minutes := duration.Minutes()
	caloriesSpent := (weight * midSpeed * minutes) / minInH
	return caloriesSpent, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if weight <= 0 || height <= 0 || steps <= 0 || duration <= 0 {
		return 0, errors.New("некоректные входные параметры")
	}
	midSpeed := meanSpeed(steps, height, duration)
	minutes := duration.Minutes()
	calories := (minutes * weight * midSpeed) / 60
	calories *= walkingCaloriesCoefficient
	return calories, nil
}
