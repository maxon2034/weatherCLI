package main

import (
	"encoding/json"
	"fmt"
	"os"
	"weatherCLI/internal/domain"
	"weatherCLI/internal/provider/openmeteo"
)

func main() {
	todayForecast := domain.Today{}
	name, lat, lon, err := openmeteo.Geocode("Minsk")

	if err != nil {
		fmt.Println("error in locating city:", err)
	}

	todayForecast, err = openmeteo.GetCurrentWeather(lat, lon)
	todayForecast.City = name
	if err != nil {
		fmt.Println("error in generating todays forecast: ", err)
	}
	fmt.Println(todayForecast)
	forecastJSON, err := json.Marshal(todayForecast)
	os.WriteFile("forecast.json", forecastJSON, 0700)
}
