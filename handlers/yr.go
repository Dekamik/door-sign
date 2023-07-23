package handlers

import (
	"door-sign/configuration"
	"door-sign/helpers"
	"door-sign/integrations"
	"log"
)

type YRForecast struct {
	Time          string
	Temperature   float64
	SymbolCode    string
	SymbolID      string
	Precipitation float64
}

func GetNowcast() *YRForecast {
	return nil
}

func UpdateYR(conf configuration.Config, maxLength int) []YRForecast {
	res, err := integrations.YRGetLocationForecast(conf.Weather.Lat, conf.Weather.Lon)
	if err != nil {
		log.Fatalln(err)
	}

	forecasts := make([]YRForecast, 0)
	for i, item := range res.Properties.Timeseries {
		time := item.Time.Local().Format("15:04")
		if i != 0 &&
			time != "00:00" &&
			time != "08:00" &&
			time != "12:00" &&
			time != "18:00" {
			continue
		}
		forecast := YRForecast{
			Time:          time,
			Temperature:   item.Data.Instant.Details.AirTemperature,
			SymbolCode:    item.Data.Next6Hours.Summary.SymbolCode,
			SymbolID:      helpers.YRSymbolsID[item.Data.Next6Hours.Summary.SymbolCode],
			Precipitation: item.Data.Next6Hours.Details.PrecipitationAmount,
		}
		forecasts = append(forecasts, forecast)
	}

	return forecasts[0:maxLength]
}
