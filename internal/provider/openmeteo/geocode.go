package openmeteo

type apiResponse struct {
	Results []struct {
		Name    string  `json:"name"`
		Country string  `json:"country"`
		Lat     float64 `json:"latitude"`
		Lon     float64 `json:"longitude"`
	} `json:"results"`
}
