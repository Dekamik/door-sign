package handlers

import (
	"door-sign/internal/config"
	"door-sign/internal/helpers"
	"door-sign/internal/integrations"
	"log"
	"time"
)

type YR interface {
	GetCurrent(conf config.Config) YRForecast
	GetForecasts(conf config.Config, maxLength int) []YRForecast
}

type YRImpl struct {
	CacheLifetime  time.Duration
	CachedResponse *Cache[*integrations.YRResponse]
}

var _ YR = &YRImpl{}

type YRForecast struct {
	Time               string
	Temperature        float64
	TemperatureColor   string
	SymbolCode         string
	SymbolID           string
	Precipitation      float64
	PrecipitationColor string
}

type Cache[T any] struct {
	ExpiresAt time.Time
	Data      T
}

func (y *YRImpl) getForecasts(conf config.Config) *integrations.YRResponse {
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

func (y *YRImpl) getTemperatureColorClass(conf config.Config, temperature float64) string {
	switch {

	case temperature < conf.Weather.Colors.TempQ1:
		return conf.Weather.Colors.ClassQ1

	case temperature < conf.Weather.Colors.TempQ2:
		return conf.Weather.Colors.ClassQ2

	case temperature < conf.Weather.Colors.TempQ3:
		return conf.Weather.Colors.ClassQ3

	case temperature < conf.Weather.Colors.TempQ4:
		return conf.Weather.Colors.ClassQ4

	default:
		return ""
	}
}

func (y *YRImpl) getPrecipitationColorClass(conf config.Config, precipitation float64) string {
	if precipitation > 0 {
		return conf.Weather.Colors.ClassPrecip
	}
	return conf.Weather.Colors.ClassNoPrecip
}

func (y *YRImpl) GetCurrent(conf config.Config) YRForecast {
	res := y.getForecasts(conf)
	latest := res.Properties.Timeseries[0]
	return YRForecast{
		Time:               latest.Time.Local().Format("15:04"),
		Temperature:        latest.Data.Instant.Details.AirTemperature,
		TemperatureColor:   y.getTemperatureColorClass(conf, latest.Data.Instant.Details.AirTemperature),
		SymbolCode:         latest.Data.Next6Hours.Summary.SymbolCode,
		SymbolID:           helpers.YRSymbolsID[latest.Data.Next6Hours.Summary.SymbolCode],
		Precipitation:      latest.Data.Next6Hours.Details.PrecipitationAmount,
		PrecipitationColor: y.getPrecipitationColorClass(conf, latest.Data.Next6Hours.Details.PrecipitationAmount),
	}
}

func (y *YRImpl) GetForecasts(conf config.Config, maxLength int) []YRForecast {
	res := y.getForecasts(conf)

	forecasts := make([]YRForecast, 0)
	for _, item := range res.Properties.Timeseries {
		timeStr := item.Time.Local().Format("15")
		if timeStr != "00" &&
			timeStr != "06" &&
			timeStr != "12" &&
			timeStr != "18" {
			continue
		}
		forecast := YRForecast{
			Time:               timeStr + "-" + item.Time.Local().Add(time.Hour*6).Format("15"),
			Temperature:        item.Data.Instant.Details.AirTemperature,
			TemperatureColor:   y.getTemperatureColorClass(conf, item.Data.Instant.Details.AirTemperature),
			SymbolCode:         item.Data.Next6Hours.Summary.SymbolCode,
			SymbolID:           helpers.YRSymbolsID[item.Data.Next6Hours.Summary.SymbolCode],
			Precipitation:      item.Data.Next6Hours.Details.PrecipitationAmount,
			PrecipitationColor: y.getPrecipitationColorClass(conf, item.Data.Next6Hours.Details.PrecipitationAmount),
		}
		forecasts = append(forecasts, forecast)
	}

	return forecasts[0:maxLength]
}

func NewYR() YR {
	return &YRImpl{
		CacheLifetime:  time.Minute * 5,
		CachedResponse: nil,
	}
}
