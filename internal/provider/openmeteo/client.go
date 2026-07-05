package openmeteo

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"weatherCLI/internal/domain"
)

type Client struct {
	HTTPClient *http.Client
}

func NewClient() *Client {
	client := Client{}
	client.HTTPClient.Timeout = time.Second * 5
	return &client
}

func (c *Client) geocode(ctx context.Context, city string) (name string, lat, lon float64, err error) {
	apiResponse := apiResponse{}
	url := "https://geocoding-api.open-meteo.com/v1/search?name=" + city
	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", 0, 0, fmt.Errorf("Error in http get: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == 404 {
			return "", 0, 0, fmt.Errorf("City is not found")
		}
		return "", 0, 0, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	err = json.NewDecoder(resp.Body).Decode(&apiResponse)
	if err != nil {
		return "", 0, 0, fmt.Errorf("Error in decoding json: %w", err)
	}
	return apiResponse.Results[0].Name, apiResponse.Results[0].Lat, apiResponse.Results[0].Lon, nil
}

func (c *Client) forecast(ctx context.Context, lat, lon float64, days int) (forecastResp, error) {
	forecast := forecastResp{}
	url := fmt.Sprintf("https://api.open-meteo.com/v1/forecast?latitude=%f&longitude=%f&forecast_days=%d&current=temperature_2m,apparent_temperature,weather_code,wind_speed_10m,wind_direction_10m,relative_humidity_2m,pressure_msl,visibility,precipitation&hourly=temperature_2m,precipitation_probability,wind_speed_10m&daily=temperature_2m_min,temperature_2m_max,precipitation_probability_max,wind_speed_10m_max,weather_code", lat, lon, days)
	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return forecastResp{}, fmt.Errorf("Error in getting response: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == 404 {
			return forecastResp{}, fmt.Errorf("Forecast not found")
		}
		return forecastResp{}, fmt.Errorf("Http error: %d", resp.StatusCode)
	}
	defer resp.Body.Close()

	json.NewDecoder(resp.Body).Decode(&forecast)
	if err != nil {
		return forecastResp{}, fmt.Errorf("Error in converting time: %w", err)
	}
	return forecast, nil
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
	today.UpdatedAt, err = time.Parse("2006-01-02", forecast.Current.UpdatedAt)
	if err != nil {
		return domain.Today{}, fmt.Errorf("error in parsing time: %w", err)
	}
	return today, nil
}

func (c *Client) GetHourly(ctx context.Context, city string, hours int) (domain.HourlyEntry, error) {
	hourly := domain.HourlyEntry{}
	_, lat, lon, err := c.geocode(ctx, city)
	if err != nil {
		return domain.HourlyEntry{}, fmt.Errorf("error in getting coordinates: %w", err)
	}
	forecast, err := c.forecast(ctx, lat, lon, 0)
	if err != nil {
		return domain.HourlyEntry{}, fmt.Errorf("error in generating forecast: %w", err)
	}
	for i := 0; i < hours; i++ {
		hourly.Time[i], err = time.Parse("2006-01-02T15:04", forecast.Hourly.Time[i])
		if err != nil {
			return domain.HourlyEntry{}, fmt.Errorf("error in parsing time on iteration %d : %w", i, err)
		}
		hourly.TemperatureC[i] = forecast.Hourly.TemperatureC[i]
		hourly.POPPercent[i] = forecast.Hourly.POPPercent[i]
		hourly.WindSpeedMS[i] = forecast.Hourly.WindSpeedMS[i]
	}
	return hourly, nil
}

func (c *Client) GetDaily(ctx context.Context, city string, days int) (domain.DailyEntry, error) {
	daily := domain.DailyEntry{}
	_, lat, lon, err := c.geocode(ctx, city)
	if err != nil {
		return domain.DailyEntry{}, fmt.Errorf("error in getting coordinates: %w", err)
	}
	forecast, err := c.forecast(ctx, lat, lon, days)
	if err != nil {
		return domain.DailyEntry{}, fmt.Errorf("error in generating forecast: %w", err)
	}

	for i := 0; i < days; i++ {
		daily.Date[i], err = time.Parse("2006-01-02", forecast.Daily.Date[i])
		if err != nil {
			return domain.DailyEntry{}, fmt.Errorf("error in parsing time on iteration %d : %w", i, err)
		}
		daily.TempMinC[i] = forecast.Daily.TempMinC[i]
		daily.TempMaxC[i] = forecast.Daily.TempMaxC[i]
		daily.POPPercent[i] = forecast.Daily.POPPercent[i]
		daily.WindSpeedMS[i] = forecast.Daily.WindSpeedMS[i]
		daily.Condition[i] = weatherCodeToText(forecast.Daily.Condition[i])
	}
	return daily, err
}
