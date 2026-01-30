package spentcalories

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	mInKm                     = 1000  // метров в километре
	minInH                    = 60    // минут в часе
	stepLengthCoefficient     = 0.414 // коэффициент длины шага от роста
	walkingCaloriesCoefficient = 0.029 // коэффициент для ходьбы
)

func parseTraining(data string) (int, string, time.Duration, error) {
	parts := strings.Split(data, ",")
	if len(parts) != 3 {
		return 0, "", 0, fmt.Errorf("неверный формат данных: ожидается 3 части, получено %d", len(parts))
	}

	// Парсим количество шагов
	stepsStr := strings.TrimSpace(parts[0])
	steps, err := strconv.Atoi(stepsStr)
	if err != nil {
		return 0, "", 0, fmt.Errorf("ошибка парсинга количества шагов: %v", err)
	}

	// Получаем вид активности
	activity := strings.TrimSpace(parts[1])

	// Парсим продолжительность
	durationStr := strings.TrimSpace(parts[2])
	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		return 0, "", 0, fmt.Errorf("ошибка парсинга продолжительности: %v", err)
	}

	return steps, activity, duration, nil
}

func distance(steps int, height float64) float64 {
	// Вычисляем длину шага
	stepLength := height * stepLengthCoefficient
	
	// Вычисляем дистанцию в метрах
	distanceMeters := float64(steps) * stepLength
	
	// Переводим в километры
	return distanceMeters / mInKm
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	// Проверяем продолжительность
	if duration <= 0 {
		return 0
	}

	// Вычисляем дистанцию
	dist := distance(steps, height)
	
	// Переводим продолжительность в часы
	durationHours := duration.Hours()
	
	if durationHours == 0 {
		return 0
	}
	
	// Вычисляем среднюю скорость
	return dist / durationHours
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	// Проверка корректности входных параметров
	if steps <= 0 {
		return 0, fmt.Errorf("количество шагов должно быть положительным числом")
	}
	if weight <= 0 {
		return 0, fmt.Errorf("вес должен быть положительным числом")
	}
	if height <= 0 {
		return 0, fmt.Errorf("рост должен быть положительным числом")
	}
	if duration <= 0 {
		return 0, fmt.Errorf("продолжительность должна быть положительной")
	}

	// Рассчитываем среднюю скорость
	speed := meanSpeed(steps, height, duration)
	if speed <= 0 {
		return 0, fmt.Errorf("невозможно рассчитать скорость")
	}

	// Переводим продолжительность в минуты
	durationMinutes := duration.Minutes()

	// Рассчитываем калории по формуле
	calories := (weight * speed * durationMinutes) / minInH
	
	return calories, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	// Проверка корректности входных параметров
	if steps <= 0 {
		return 0, fmt.Errorf("количество шагов должно быть положительным числом")
	}
	if weight <= 0 {
		return 0, fmt.Errorf("вес должен быть положительным числом")
	}
	if height <= 0 {
		return 0, fmt.Errorf("рост должен быть положительным числом")
	}
	if duration <= 0 {
		return 0, fmt.Errorf("продолжительность должна быть положительной")
	}

	// Рассчитываем среднюю скорость
	speed := meanSpeed(steps, height, duration)
	if speed <= 0 {
		return 0, fmt.Errorf("невозможно рассчитать скорость")
	}

	// Переводим продолжительность в минуты
	durationMinutes := duration.Minutes()

	// Рассчитываем калории по формуле для бега
	calories := (weight * speed * durationMinutes) / minInH
	
	// Применяем корректирующий коэффициент для ходьбы
	calories *= walkingCaloriesCoefficient
	
	return calories, nil
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, activity, duration, err := parseTraining(data)
	if err != nil {
		return "", fmt.Errorf("ошибка парсинга тренировки: %v", err)
	}

	var calories float64
	
	// Вычисляем дистанцию
	dist := distance(steps, height)
	
	// Вычисляем среднюю скорость
	speed := meanSpeed(steps, height, duration)
	
	// Определяем тип тренировки и вычисляем калории
	switch activity {
	case "Ходьба":
		calories, err = WalkingSpentCalories(steps, weight, height, duration)
		if err != nil {
			return "", fmt.Errorf("ошибка расчета калорий для ходьбы: %v", err)
		}
	case "Бег":
		calories, err = RunningSpentCalories(steps, weight, height, duration)
		if err != nil {
			return "", fmt.Errorf("ошибка расчета калорий для бега: %v", err)
		}
	default:
		return "", fmt.Errorf("неизвестный тип тренировки")
	}

	// Формируем строку с информацией о тренировке
	result := fmt.Sprintf("Тип тренировки: %s\n", activity)
	result += fmt.Sprintf("Длительность: %.2f ч.\n", duration.Hours())
	result += fmt.Sprintf("Дистанция: %.2f км.\n", dist)
	result += fmt.Sprintf("Скорость: %.2f км/ч\n", speed)
	result += fmt.Sprintf("Сожгли калорий: %.2f", calories)
	
	return result, nil
}
