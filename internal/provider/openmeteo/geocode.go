package openmeteo

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type apiResponse struct {
	Results []struct {
		Name    string  `json:"name"`
		Country string  `json:"country"`
		Lat     float64 `json:"latitude"`
		Lon     float64 `json:"longitude"`
	} `json:"results"`
}

func geocode(city string) (name string, lat, lon float64, err error) {
	apiResponse := apiResponse{}
	url := "https://geocoding-api.open-meteo.com/v1/search?name=" + city
	resp, err := http.Get(url)
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
