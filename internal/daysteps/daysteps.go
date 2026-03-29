// Пакет daysteps отвечает за учёт активности в течение дня.
package daysteps

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
	// Формат времени
	//layout = 3h50m
)

// strconv.Atoi(s string) (int, error)
func parsePackage(data string) (int, time.Duration, error) {
	// Разделяем строку на слайс строк
	slice := strings.Split(data, ",")
	if len(slice) != 2 {
		return 0, 0, errors.New("Длина слайса меньше или больше 2")
	}
	// steps преобразует из слайса количество шагов
	steps, err := strconv.Atoi(slice[0])
	if err != nil {
		err = fmt.Errorf("ошибка преобразовния строки в целое число : %v, err")
		return 0, 0, err
	}
	if steps <= 0 {
		return 0, 0, errors.New("количество шагов меньше или равно 0")
	}

	// tDuration преобразует из слайса строку в тип time.Duration
	tDuration, err := time.ParseDuration(slice[1])
	if err != nil {
		err = fmt.Errorf("ошибка преобразовния строки в тип time.Duration : %v", err)
		return 0, 0, err
	}
	if tDuration <= 0 {
		return 0, 0, errors.New("Время меньше или равно 0")
	} // не факт что нужна

	return steps, tDuration, nil
}

func DayActionInfo(data string, weight, height float64) string {
	//Получаем данные о кол-ве шагов и продолжит. прогулки
	steps, tDuration, err := parsePackage(data)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	if steps <= 0 {
		var s string
		return s
	}
	distance := float64(steps) * stepLength
	distance /= float64(mInKm)
	caloriesSpent, err := WalkingSpentCalories(steps, weight, height, tDuration)
	if err != nil {
		return ""
	}
	return ""
}
