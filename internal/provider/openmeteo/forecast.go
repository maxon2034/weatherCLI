package openmeteo

type forecastResp struct {
	Current struct {
		City             string
		TemperatureC     float64 `json:"temperature_2m"`
		FeelsLikeC       float64 `json:"apparent_temperature"`
		WeatherCode      int     `json:"weather_code"`
		WindSpeedMS      float64 `json:"wind_speed_10m"`
		WindDirectionDeg int     `json:"wind_direction_10m"`
		HumidityPercent  int     `json:"relative_humidity_2m"`
		PressureHPa      float64 `json:"pressure_msl"`
		VisibilityKm     float64 `json:"visibility"`
		PrecipitationMm  float64 `json:"precipitation"`
		UpdatedAt        string  `json:"time"`
	} `json:"current"`

	Hourly struct {
		Time         []string  `json:"time"`
		TemperatureC []float64 `json:"temperature_2m"`
		POPPercent   []int     `json:"precipitation_probability"`
		WindSpeedMS  []float64 `json:"wind_speed_10m"`
	} `json:"hourly"`

	Daily struct {
		Date        []string  `json:"time"`
		TempMinC    []float64 `json:"temperature_2m_min"`
		TempMaxC    []float64 `json:"temperature_2m_max"`
		POPPercent  []int     `json:"precipitation_probability_max"`
		WindSpeedMS []float64 `json:"wind_speed_10m_max"`
		Condition   []int     `json:"weather_code"` // API возвращает int, а не string!
	} `json:"daily"`
}

func weatherCodeToText(code int) string {
	switch code {
	case 0:
		return "Ясно"
	case 1:
		return "Преимущественно ясно"
	case 2:
		return "Переменная облачность"
	case 3:
		return "Пасмурно"
	case 45:
		return "Туман"
	case 48:
		return "Осаждающийся туман (изморозь)"
	case 51:
		return "Слабая морось"
	case 53:
		return "Умеренная морось"
	case 55:
		return "Интенсивная морось"
	case 56:
		return "Слабая переохлажденная морось"
	case 57:
		return "Интенсивная переохлажденная морось"
	case 61:
		return "Слабый дождь"
	case 63:
		return "Умеренный дождь"
	case 65:
		return "Сильный дождь"
	case 66:
		return "Слабый переохлажденный дождь"
	case 67:
		return "Сильный переохлажденный дождь"
	case 71:
		return "Слабый снегопад"
	case 73:
		return "Умеренный снегопад"
	case 75:
		return "Сильный снегопад"
	case 77:
		return "Снежные зерна"
	case 80:
		return "Слабый ливневый дождь"
	case 81:
		return "Умеренный ливневый дождь"
	case 82:
		return "Сильный ливневый дождь"
	case 85:
		return "Слабый ливневый снег"
	case 86:
		return "Сильный ливневый снег"
	case 95:
		return "Гроза"
	case 96:
		return "Гроза со слабым градом"
	case 99:
		return "Гроза с сильным градом"
	default:
		return "Неизвестные погодные условия"
	}
}
