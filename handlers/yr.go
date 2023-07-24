package handlers

import (
	"door-sign/configuration"
	"door-sign/helpers"
	"door-sign/integrations"
	"log"
	"time"
)

type YR interface {
	GetCurrent(conf configuration.Config) YRForecast
	GetForecasts(conf configuration.Config, maxLength int) []YRForecast
}

type YRImpl struct {
	CacheLifetime  time.Duration
	CachedResponse *Cache[*integrations.YRResponse]
}

var _ YR = &YRImpl{}

type YRForecast struct {
	Time          string
	Temperature   float64
	SymbolCode    string
	SymbolID      string
	Precipitation float64
}

type Cache[T any] struct {
	ExpiresAt time.Time
	Data      T
}

func (y *YRImpl) getForecasts(conf configuration.Config) *integrations.YRResponse {
	if y.CachedResponse != nil && time.Now().Before(y.CachedResponse.ExpiresAt) {
		log.Println("YR: Getting cached response")
		return y.CachedResponse.Data
	}

	log.Println("YR: Getting new repsonse from met.no")
	res, err := integrations.YRGetLocationForecast(conf.Weather.Lat, conf.Weather.Lon)
	if err != nil {
		log.Fatalln(err)
	}
	y.CachedResponse = &Cache[*integrations.YRResponse]{
		ExpiresAt: time.Now().Add(y.CacheLifetime),
		Data:      res,
	}

	return res
}

func (y *YRImpl) GetCurrent(conf configuration.Config) YRForecast {
	res := y.getForecasts(conf)
	latest := res.Properties.Timeseries[0]
	return YRForecast{
		Time:          latest.Time.Local().Format("15:04"),
		Temperature:   latest.Data.Instant.Details.AirTemperature,
		SymbolCode:    latest.Data.Next6Hours.Summary.SymbolCode,
		SymbolID:      helpers.YRSymbolsID[latest.Data.Next6Hours.Summary.SymbolCode],
		Precipitation: latest.Data.Next6Hours.Details.PrecipitationAmount,
	}
}

func (y *YRImpl) GetForecasts(conf configuration.Config, maxLength int) []YRForecast {
	res := y.getForecasts(conf)

	forecasts := make([]YRForecast, 0)
	for _, item := range res.Properties.Timeseries {
		time := item.Time.Local().Format("15:04")
		if time != "00:00" &&
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

func NewYR() YR {
	return &YRImpl{
		CacheLifetime:  time.Second * 50,
		CachedResponse: nil,
	}
}
