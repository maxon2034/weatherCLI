package openmeteo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"weatherCLI/internal/domain"
)

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

func GetCurrentWeather(lat, lon float64) (domain.Today, error) {
	forecast := Forecast{}
	today := domain.Today{}
	url := fmt.Sprintf("https://api.open-meteo.com/v1/forecast?latitude=%f&longitude=%f&current=temperature_2m&current=apparent_temperature&current=weather_code&wind_speed_unit=ms&current=wind_speed_10m&current=wind_direction_10m&current=relative_humidity_2m&current=pressure_msl&current=visibility&current=precipitation", lat, lon)
	resp, err := http.Get(url)
	if err != nil {
		return domain.Today{}, fmt.Errorf("Error in getting response: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == 404 {
			return domain.Today{}, fmt.Errorf("Forecast not found")
		}
		return domain.Today{}, fmt.Errorf("Http error: %d", resp.StatusCode)
	}
	defer resp.Body.Close()

	json.NewDecoder(resp.Body).Decode(&forecast)
	today.TemperatureC = forecast.Current.TemperatureC
	today.FeelsLikeC = forecast.Current.FeelsLikeC
	today.Condition = weatherCodeToText(forecast.Current.WeatherCode)
	today.WindSpeedMS = forecast.Current.WindSpeedMS
	today.WindDirectionDeg = forecast.Current.WindDirectionDeg
	today.HumidityPercent = forecast.Current.HumidityPercent
	today.PressureHPa = forecast.Current.PressureHPa
	today.VisibilityKm = forecast.Current.VisibilityKm
	today.PrecipitationMm = forecast.Current.PrecipitationMm
	today.UpdatedAt, err = time.Parse("2006-01-02T15:04", forecast.Current.UpdatedAt)
	if err != nil {
		return domain.Today{}, fmt.Errorf("Error in converting time: %w", err)
	}
	return today, nil
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
