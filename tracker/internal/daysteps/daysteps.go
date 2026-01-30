package daysteps

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/spentcalories"
)

const (
	stepLength = 0.65 // длина шага в метрах
	mInKm      = 1000 // метров в километре
)

func parsePackage(data string) (int, time.Duration, error) {
	parts := strings.Split(data, ",")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("неверный формат данных")
	}

	// Парсим количество шагов
	stepsStr := strings.TrimSpace(parts[0])
	if stepsStr == "" {
		return 0, 0, fmt.Errorf("количество шагов не указано")
	}
	
	steps, err := strconv.Atoi(stepsStr)
	if err != nil {
		return 0, 0, fmt.Errorf("ошибка парсинга количества шагов")
	}

	// Проверяем, что шаги больше 0
	if steps <= 0 {
		return 0, 0, fmt.Errorf("количество шагов должно быть положительным числом")
	}

	// Парсим продолжительность
	durationStr := strings.TrimSpace(parts[1])
	if durationStr == "" {
		return 0, 0, fmt.Errorf("продолжительность не указана")
	}
	
	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		return 0, 0, fmt.Errorf("ошибка парсинга продолжительности")
	}

	// Проверяем, что продолжительность больше 0
	if duration <= 0 {
		return 0, 0, fmt.Errorf("продолжительность должна быть положительной")
	}

	return steps, duration, nil
}

func DayActionInfo(data string, weight, height float64) string {
	steps, duration, err := parsePackage(data)
	if err != nil {
		// Выводим ошибку в лог
		fmt.Printf("Ошибка: %v\n", err)
		return ""
	}

	// Проверяем количество шагов
	if steps <= 0 {
		return ""
	}

	// Вычисляем дистанцию в метрах
	distanceMeters := float64(steps) * stepLength
	
	// Переводим дистанцию в километры
	distanceKm := distanceMeters / mInKm

	// Вычисляем калории
	calories, err := spentcalories.WalkingSpentCalories(steps, weight, height, duration)
	if err != nil {
		fmt.Printf("Ошибка расчета калорий: %v\n", err)
		return ""
	}

	// Формируем строку результата
	result := fmt.Sprintf(
		"Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.",
		steps,
		distanceKm,
		calories,
	)

	return result
}
