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

func (c *Client) getCurrentWeather(ctx context.Context, lat, lon float64) (domain.Today, error) {
	forecast := Forecast{}
	today := domain.Today{}
	url := fmt.Sprintf("https://api.open-meteo.com/v1/forecast?latitude=%f&longitude=%f&current=temperature_2m&current=apparent_temperature&current=weather_code&wind_speed_unit=ms&current=wind_speed_10m&current=wind_direction_10m&current=relative_humidity_2m&current=pressure_msl&current=visibility&current=precipitation", lat, lon)
	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
	resp, err := c.HTTPClient.Do(req)
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

func (c *Client) GetToday(ctx context.Context, city string) (domain.Today, error) {
	name, lat, lon, err := c.geocode(ctx, city)
	if err != nil {
		return domain.Today{}, fmt.Errorf("error in getting coordinates: %w", err)
	}
	today, err := c.getCurrentWeather(ctx, lat, lon)
	if err != nil {
		return domain.Today{}, fmt.Errorf("error in generating forecast: %w", err)
	}
	today.City = name
	return today, nil
}
