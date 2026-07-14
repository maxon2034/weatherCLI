package ui

import (
	"fmt"
	"strconv"
	"time"
	"weatherCLI/internal/cache"
	"weatherCLI/internal/domain"
)

func colorTemp(celsius float64, text string) string {
	switch {
	case celsius < 0:
		return blue + text + reset
	case celsius >= 0 && celsius <= 25:
		return yellow + text + reset
	case celsius > 25:
		return red + text + reset
	}
	return text
}

func Header(city string, cached bool, fetchedAt time.Time) string {
	fetchStr := fetchedAt.Format("4")
	if !cached {
		return city + " • обновлено " + fetchStr + " мин назад\n"
	}
	return city + " • обновлено " + fetchStr + " мин назад • " + "из кэша\n"
}

func iconForCondition(cond string) string {
	switch cond {
	case "Ясно", "Преимущественно ясно":
		return "☀"

	case "Переменная облачность", "Пасмурно":
		return "☁"

	case "Туман", "Осаждающийся туман (изморозь)":
		return "🌫" // Или обычное облако "☁", если нужно строго по списку

	case "Слабая морось", "Умеренная морось", "Интенсивная морось",
		"Слабая переохлажденная морось", "Интенсивная переохлажденная морось",
		"Слабый дождь", "Умеренный дождь", "Сильный дождь",
		"Слабый переохлажденный дождь", "Сильный переохлажденный дождь",
		"Слабый ливневый дождь", "Умеренный ливневый дождь", "Сильный ливневый дождь":
		return "🌧"

	case "Слабый снегопад", "Умеренный снегопад", "Сильный снегопад", "Снежные зерна":
		return "❄"

	case "Слабый ливневый снег", "Сильный ливневый снег":
		return "🌨" // Ливневый снег. Можно заменить на "❄"

	case "Гроза", "Гроза со слабым градом", "Гроза с сильным градом":
		return "⛈"

	default:
		return "❓" // На случай, если придет неизвестная строка
	}
}

func RenderToday(t domain.Today) string {
	c := cache.New()
	_, fetchedAt, ok := c.Get("today:" + t.City)
	header := Header(t.City, ok, fetchedAt)
	c.Set("today:"+t.City, t, 5*time.Minute)

	icon := iconForCondition(t.Condition)
	announcement := "Сегодня в " + t.City + " [ " + icon + " ]\n"

	// Здесь жесткое выравнивание не критично, так как это не таблица,
	// но передаем просто чистые строки температуры в цвет
	tempStr := fmt.Sprintf("%.1f", t.TemperatureC)
	feelsStr := fmt.Sprintf("%.1f", t.FeelsLikeC)

	temp := "Температура воздуха: " + colorTemp(t.TemperatureC, tempStr) + "°С (ощущается как " + colorTemp(t.FeelsLikeC, feelsStr) + "°С)\n"
	cond := "Условия: " + t.Condition + "\n"
	windSp := fmt.Sprintln("Скорость ветра:", t.WindSpeedMS, "; Направление ветра: ", t.WindDirectionDeg, "°")
	humid := fmt.Sprintln("Влажность воздуха: ", t.HumidityPercent, "%; Атмосферное давление: ", t.PressureHPa, " ГПа")
	visPrec := fmt.Sprintln("Видимость составляет", t.VisibilityKm, "км; Количество осадков: ", t.PrecipitationMm, "мм")
	return fmt.Sprint(header, announcement, temp, cond, windSp, humid, visPrec)
}

func RenderHourly(list []domain.HourlyEntry) string {
	first := "Почасовой прогноз (" + strconv.Itoa(len(list)) + " час.)"
	header := fmt.Sprintf("%-8s | %6s | %7s | %10s", "Время", "t°C", "Осадки", "Ветер м/с")
	split := "---------+--------+---------+------------" + "\n"
	final := first + "\n" + split + header + "\n" + split
	for _, v := range list {
		// 1. Форматируем температуру в строку с фиксированной шириной (например, 5 символов: " 21.3")
		rawTempStr := fmt.Sprintf("%5.1f", v.TemperatureC)
		// 2. Раскрашиваем уже выровненную строку
		coloredTemp := colorTemp(v.TemperatureC, rawTempStr)

		// Выводим %s без дополнительного указания ширины для температуры, так как она уже отформатирована
		final += fmt.Sprintf("%-8s | %s  | %6d%% | %10.1f", v.Time.Format("15:04"), coloredTemp, v.POPPercent, v.WindSpeedMS) + "\n"
	}
	return final + split
}

func RenderDaily(list []domain.DailyEntry) string {
	first := "Прогноз на неделю"
	header := fmt.Sprintf("%-12s | %-20s | %6s | %6s | %7s", "Дата", "Условия", "Мин°C", "Макс°C", "Осадки")
	split := "-------------+----------------------+--------+--------+----------" + "\n"
	final := first + "\n" + split + header + "\n" + split
	for _, v := range list {
		// Форматируем отдельно минимальную и максимальную температуру под ширину 5 символов
		rawMinStr := fmt.Sprintf("%5.1f", v.TempMinC)
		rawMaxStr := fmt.Sprintf("%5.1f", v.TempMaxC)

		coloredMin := colorTemp(v.TempMinC, rawMinStr)
		coloredMax := colorTemp(v.TempMaxC, rawMaxStr)

		// Объединяем иконку и дату, задаем фиксированную ширину (например, 12 символов под дату с иконкой)
		dateWithIcon := v.Date.Format("02 Jan") + " " + iconForCondition(v.Condition)

		final += fmt.Sprintf("%-12s | %-20s | %s | %s | %7d%%", dateWithIcon, v.Condition, coloredMin, coloredMax, v.POPPercent) + "\n"
	}
	return final + split
}

func RenderMenu() string {
	split := "────────────────────────────────────────────────────────────\n"
	menu := "[1] Почасовой (12 ч)  [2] На 7 дней  [C] Сменить город  [R] Обновить  [Q] Выход"
	return split + menu
}
