package ui

import (
	"context"
	"fmt"
	"strconv"
	"time"
	"weatherCLI/internal/cache"
	"weatherCLI/internal/config"
	"weatherCLI/internal/domain"
	"weatherCLI/internal/provider/openmeteo"
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

	temp := "Температура воздуха: " + colorTemp(t.TemperatureC, fmt.Sprint(t.TemperatureC)) + "°С (ощущается как " + colorTemp(t.FeelsLikeC, fmt.Sprint(t.FeelsLikeC)) + "°С)\n"
	cond := "Условия: " + t.Condition + "\n"
	windSp := fmt.Sprintln("Скорость ветра:", t.WindSpeedMS, "; Направление ветра: ", t.WindDirectionDeg, "°")
	humid := fmt.Sprintln("Влажность воздуха: ", t.HumidityPercent, "%; Атмосферное давление: ", t.PressureHPa, " ГПа")
	visPrec := fmt.Sprintln("Видимость составляет", t.VisibilityKm, "км; Количество осадков: ", t.PrecipitationMm, "мм")
	return fmt.Sprint(header, announcement, temp, cond, windSp, humid, visPrec)
}

func RenderHourly(list []domain.HourlyEntry) string {
	first := "Почасовой прогноз (" + strconv.Itoa(len(list)) + "час.)"
	header := fmt.Sprintf("%-8s | %6s | %7s | %7s", "Время", "t°C", "Осадки", "Ветер м/с")
	split := "---------+--------+---------+-----------" + "\n"
	final := first + "\n" + split + header + "\n" + split
	for _, v := range list {
		final += fmt.Sprintf("%-8s | %6.1s | %6d%% |  %8.1f", v.Time.Format("15:04"), colorTemp(v.TemperatureC, fmt.Sprint(v.TemperatureC)), v.POPPercent, v.WindSpeedMS) + "\n"
	}
	return final + split
}

func RenderDaily(list []domain.DailyEntry) string {
	first := "Прогноз на неделю"
	header := fmt.Sprintf("%-9s | %-20s | %7s | %7s | %7s", "Дата", "", "Мин°C", "Макс°C", "Осадки")
	split := "----------+----------------------+---------+--------+----------" + "\n"
	final := first + "\n" + split + header + "\n" + split
	for _, v := range list {
		final += fmt.Sprintf("%-9s | %-20s | %10s |  %10s | %7d%%", v.Date.Format("02 Jan")+" "+iconForCondition(v.Condition), v.Condition, colorTemp(v.TempMinC, fmt.Sprint(v.TempMinC)), colorTemp(v.TempMaxC, fmt.Sprint(v.TempMaxC)), v.POPPercent) + "\n"
	}
	return final + split
}

func RenderMenu() string {
	c := openmeteo.NewClient()
	ctx := context.Background()

	conf, err := config.Load()
	if err != nil {
		return fmt.Sprint("Error in loading configuration:", err)
	}
	defCity := conf.DefaultCity
	today, err := c.GetToday(ctx, defCity)
	if err != nil {
		return fmt.Sprint("Error in getting today's forecast: ", err)
	}
	split := "────────────────────────────────────────────────────────────\n"
	menu := "[1] Почасовой (12 ч)  [2] На 7 дней  [C] Сменить город  [R] Обновить  [Q] Выход"
	return RenderToday(today) + split + menu
}
func RenderMenuMine() string {
	c := openmeteo.NewClient()
	ctx := context.Background()

	conf, err := config.Load()
	if err != nil {
		return fmt.Sprint("Error in loading configuration:", err)
	}
	defCity := conf.DefaultCity
	today, err := c.GetToday(ctx, defCity)
	if err != nil {
		return fmt.Sprint("Error in getting today's forecast: ", err)
	}
	fmt.Println(RenderToday(today))
OuterLoop:
	for {
		var input string
		fmt.Scan(&input)
		switch input {
		case "1":
			fmt.Println("in progress...")
			time.Sleep(5 * time.Second)
			fmt.Println(RenderToday(today))
			continue
		case "2":
			fmt.Println("in regress..")
			time.Sleep(5 * time.Second)
			fmt.Println(RenderToday(today))
			continue
		case "Q", "q":
			fmt.Println("Спасибо за то, что остаетесь с нами")
			break OuterLoop
		}
	}
	return ""
}
