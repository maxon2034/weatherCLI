package openmeteo

type Forecast struct {
	Current struct {
		TemperatureC     float64 `json:"temperature_2m"`
		FeelsLikeC       float64 `json:"apparent_temperature"`
		WeatherCode      int     `json:"weather_code"`
		WindSpeedMS      float64 `json:"wind_speed_10m"`
		WindDirectionDeg int     `json:"wind_direction_10m"`
		HumidityPercent  int     `json:"relative_humidity_2m"`
		PressureHPa      int     `json:"pressure_msl"`
		VisibilityKm     float64 `json:"visibility"`
		PrecipitationMm  float64 `json:"precipitation"`
		UpdatedAt        string  `json:"time"`
	} `json:"current"`
}

type forecastResp struct {
	Current struct {
		City             string
		TemperatureC     float64 `json:"temperature_2m"`
		FeelsLikeC       float64 `json:"apparent_temperature"`
		WeatherCode      int     `json:"weather_code"`
		WindSpeedMS      float64 `json:"wind_speed_10m"`
		WindDirectionDeg int     `json:"wind_direction_10m"`
		HumidityPercent  int     `json:"relative_humidity_2m"`
		PressureHPa      int     `json:"pressure_msl"`
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
		return "Ясное небо"
	case 1:
		return "Преимущественно ясно"
	case 2:
		return "Переменная облачность"
	case 3:
		return "Пасмурно"
	case 45:
		return "Туман"
	case 48:
		return "Иней с выпадением осадков"
	case 51:
		return "Легкий струйный полив"
	case 52:
		return "Умеренный струйный полив"
	case 53:
		return "Плотный струйный полив"
	case 56:
		return "Легкая ледяная морось"
	case 57:
		return "Насыщенная ледяная морось"
	case 61:
		return "Слабый дождь"
	case 63:
		return "Умеренный дождь"
	case 65:
		return "Сильный дождь"
	case 66:
		return "Слабый ледяной дождь"
	case 67:
		return "Легкий ледяной дождь"
	case 71:
		return "Слабый снегопад"
	case 73:
		return "Умеренный снегопад"
	case 75:
		return "Сильный снегопад"
	case 77:
		return "Град"
	case 80:
		return "Легкий ливень"
	case 81:
		return "Умеренный ливень"
	case 82:
		return "Жестокий ливень"
	case 85:
		return "Легкий дождь со снегом"
	case 86:
		return "Сильный дождь со снегом"
	case 95:
		return "Гром"
	case 96:
		return "Гром с легким градом"
	case 99:
		return "Гром с сильным градом"
	default:
		return ""
	}
}
