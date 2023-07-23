package integrations

import (
	"fmt"
	"time"
)

type YRResponse struct {
	Type     string `json:"type"`
	Geometry struct {
		Type        string    `json:"type"`
		Coordinates []float64 `json:"coordinates"`
	} `json:"geometry"`
	Properties struct {
		Meta struct {
			UpdatedAt time.Time `json:"updated_at"`
			Units     struct {
				AirPressureAtSeaLevel string `json:"air_pressure_at_sea_level"`
				AirTemperature        string `json:"air_temperature"`
				CloudAreaFraction     string `json:"cloud_area_fraction"`
				PrecipitationAmount   string `json:"precipitation_amount"`
				RelativeHumidity      string `json:"relative_humidity"`
				WindFromDirection     string `json:"wind_from_direction"`
				WindSpeed             string `json:"wind_speed"`
			} `json:"units"`
		} `json:"meta"`
		Timeseries []struct {
			Time time.Time `json:"time"`
			Data struct {
				Instant struct {
					Details struct {
						AirPressureAtSeaLevel float64 `json:"air_pressure_at_sea_level"`
						AirTemperature        float64 `json:"air_temperature"`
						CloudAreaFraction     float64 `json:"cloud_area_fraction"`
						RelativeHumidity      float64 `json:"relative_humidity"`
						WindFromDirection     float64 `json:"wind_from_direction"`
						WindSpeed             float64 `json:"wind_speed"`
					} `json:"details"`
				} `json:"instant"`
				Next12Hours struct {
					Summary struct {
						SymbolCode string `json:"symbol_code"`
					} `json:"summary"`
				} `json:"next_12_hours,omitempty"`
				Next1Hours struct {
					Summary struct {
						SymbolCode string `json:"symbol_code"`
					} `json:"summary"`
					Details struct {
						PrecipitationAmount float64 `json:"precipitation_amount"`
					} `json:"details"`
				} `json:"next_1_hours,omitempty"`
				Next6Hours struct {
					Summary struct {
						SymbolCode string `json:"symbol_code"`
					} `json:"summary"`
					Details struct {
						PrecipitationAmount float64 `json:"precipitation_amount"`
					} `json:"details"`
				} `json:"next_6_hours,omitempty"`
			} `json:"data"`
		} `json:"timeseries"`
	} `json:"properties"`
}

func YRGetLocationForecast(lat float32, lon float32) (*YRResponse, error) {
	url := fmt.Sprintf("https://api.met.no/weatherapi/locationforecast/2.0/compact?lat=%f&lon=%f", lat, lon)
	return get[YRResponse](url)
}
