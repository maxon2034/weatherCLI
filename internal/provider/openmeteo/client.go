package openmeteo

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"weatherCLI/internal/domain"
	"weatherCLI/internal/retry"
)

type Client struct {
	HTTPClient *http.Client
}

func NewClient() *Client {
	client := Client{}
	client.HTTPClient = &http.Client{}
	client.HTTPClient.Timeout = time.Second * 5
	return &client
}

func (c *Client) geocode(ctx context.Context, city string) (name string, lat, lon float64, err error) {
	apiResponse := apiResponse{}
	url := "https://geocoding-api.open-meteo.com/v1/search?name=" + city
	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
	err = retry.Do(ctx, 5, time.Millisecond*250, func() error {
		resp, err := c.HTTPClient.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			if resp.StatusCode == 404 {
				return fmt.Errorf("City is not found")
			}
			return fmt.Errorf("unexpected status: %d", resp.StatusCode)
		}

		err = json.NewDecoder(resp.Body).Decode(&apiResponse)
		if err != nil {
			return fmt.Errorf("Error in decoding json: %w", err)
		}
		return err
	})
	if err != nil {
		return "", 0, 0, fmt.Errorf("Error in network availabilty: %w", err)
	}
	return apiResponse.Results[0].Name, apiResponse.Results[0].Lat, apiResponse.Results[0].Lon, nil
}

func (c *Client) forecast(ctx context.Context, lat, lon float64, days int) (*forecastResp, error) {
	forecast := forecastResp{}
	url := fmt.Sprintf("https://api.open-meteo.com/v1/forecast?latitude=%f&longitude=%f&forecast_days=%d&current=temperature_2m,apparent_temperature,weather_code,wind_speed_10m,wind_direction_10m,relative_humidity_2m,pressure_msl,visibility,precipitation&hourly=temperature_2m,precipitation_probability,wind_speed_10m&daily=temperature_2m_min,temperature_2m_max,precipitation_probability_max,wind_speed_10m_max,weather_code", lat, lon, days)
	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
	err := retry.Do(ctx, 5, time.Millisecond*250, func() error {
		resp, err := c.HTTPClient.Do(req)
		if err != nil {
			return fmt.Errorf("Error in getting response: %w", err)
		}
		if resp.StatusCode != http.StatusOK {
			if resp.StatusCode == 404 {
				return fmt.Errorf("Forecast not found")
			}
			return fmt.Errorf("Http error: %d", resp.StatusCode)
		}
		defer resp.Body.Close()

		err = json.NewDecoder(resp.Body).Decode(&forecast)
		if err != nil {
			return fmt.Errorf("Error in converting time: %w", err)
		}
		return err
	})
	if err != nil {
		return &forecastResp{}, fmt.Errorf("Error in network availabilty: %w", err)
	}
	return &forecast, nil
}

func (c *Client) GetToday(ctx context.Context, city string) (domain.Today, error) {
	today := domain.Today{}
	name, lat, lon, err := c.geocode(ctx, city)
	if err != nil {
		return domain.Today{}, fmt.Errorf("error in getting coordinates: %w", err)
	}
	forecast, err := c.forecast(ctx, lat, lon, 0)
	if err != nil {
		return domain.Today{}, fmt.Errorf("error in generating forecast: %w", err)
	}

	forecast.Current.City = name
	today.City = forecast.Current.City
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
		return domain.Today{}, fmt.Errorf("error in parsing time: %w", err)
	}

	return today, nil
}

func (c *Client) GetHourly(ctx context.Context, city string, hours int) ([]domain.HourlyEntry, error) {
	hourly := make([]domain.HourlyEntry, hours)
	_, lat, lon, err := c.geocode(ctx, city)
	if err != nil {
		return []domain.HourlyEntry{}, fmt.Errorf("error in getting coordinates: %w", err)
	}
	forecast, err := c.forecast(ctx, lat, lon, hours)
	if err != nil {
		return []domain.HourlyEntry{}, fmt.Errorf("error in generating forecast: %w", err)
	}
	for i := 0; i < hours; i++ {
		hourly[i].Time, err = time.Parse("2006-01-02T15:04", forecast.Hourly.Time[i])
		if err != nil {
			return []domain.HourlyEntry{}, fmt.Errorf("error in parsing time on iteration %d : %w", i, err)
		}
		hourly[i].TemperatureC = forecast.Hourly.TemperatureC[i]
		hourly[i].POPPercent = forecast.Hourly.POPPercent[i]
		hourly[i].WindSpeedMS = forecast.Hourly.WindSpeedMS[i]
	}
	return hourly, nil
}

func (c *Client) GetDaily(ctx context.Context, city string, days int) ([]domain.DailyEntry, error) {
	daily := make([]domain.DailyEntry, days)
	_, lat, lon, err := c.geocode(ctx, city)
	if err != nil {
		return []domain.DailyEntry{}, fmt.Errorf("error in getting coordinates: %w", err)
	}
	forecast, err := c.forecast(ctx, lat, lon, days)
	if err != nil {
		return []domain.DailyEntry{}, fmt.Errorf("error in generating forecast: %w", err)
	}

	for i := 0; i < days; i++ {
		daily[i].Date, err = time.Parse("2006-01-02", forecast.Daily.Date[i])
		if err != nil {
			return []domain.DailyEntry{}, fmt.Errorf("error in parsing time on iteration %d : %w", i, err)
		}
		daily[i].TempMinC = forecast.Daily.TempMinC[i]
		daily[i].TempMaxC = forecast.Daily.TempMaxC[i]
		daily[i].POPPercent = forecast.Daily.POPPercent[i]
		daily[i].WindSpeedMS = forecast.Daily.WindSpeedMS[i]
		daily[i].Condition = weatherCodeToText(forecast.Daily.Condition[i])
	}
	return daily, err
}

func (c *Client) GetTodayFormatted(ctx context.Context, city string) (string, error) {
	today, err := c.GetToday(ctx, city)
	if err != nil {
		return "", fmt.Errorf("Unable to format the forecast: %w", err)
	}
	forecast := "Прогноз погоды на " + today.UpdatedAt.Format("01-02-2006 15:04") + "\n"
	cityF := fmt.Sprintln("Город: ", today.City)
	temp := fmt.Sprintln("Температура воздуха: ", today.TemperatureC, "°С (ощущается как ", today.FeelsLikeC, "°С)")
	cond := today.Condition + "\n"
	windSp := fmt.Sprintln("Скорость ветра:", today.WindSpeedMS, "; Направление ветра: ", today.WindDirectionDeg, "°")
	humid := fmt.Sprintln("Влажность воздуха: ", today.HumidityPercent, "%; Атмосферное давление: ", today.PressureHPa, " ГПа")
	visPrec := fmt.Sprintln("Видимость составляет", today.VisibilityKm, "км; Количество осадков: ", today.PrecipitationMm, "мм")
	return fmt.Sprint(forecast, cityF, temp, cond, windSp, humid, visPrec), nil
}
