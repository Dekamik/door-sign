package yr

import (
	"door-sign/internal/config"
	"door-sign/internal/handlers/timeanddate"
	"door-sign/internal/helpers"
	"door-sign/internal/integrations/v1"
	"log/slog"
	"math"
	"os"
	"time"
)

type YR interface {
	GetCurrent(conf config.Config) YRForecast
	GetForecasts(conf config.Config, maxLength int) []YRForecast
	GetFullForecasts(conf config.Config) []YRFullForecast
}

type YRImpl struct {
	CacheLifetime          time.Duration
	CachedForecastResponse *Cache[*v1.YRResponse]
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

type YRFullForecast struct {
	Time                     string
	AirPressureAtSeaLevel    float64
	Temperature              float64
	TemperatureColor         string
	TemperatureMax           float64
	TemperatureMaxColor      string
	TemperatureMin           float64
	TemperatureMinColor      string
	DewPointTemperature      float64
	DewPointTemperatureColor string
	CloudAreaFraction        float64
	CloudAreaFractionHigh    float64
	CloudAreaFractionLow     float64
	CloudAreaFractionMedium  float64
	FogAreaFraction          float64
	RelativeHumidity         float64
	UltravioletIndexClearSky float64
	WindFromDirection        float64
	WindSpeed                float64
	WindSpeedOfGust          float64
	WindSpeedPercentile10    float64
	WindSpeedPercentile90    float64
	SymbolName               string
	SymbolCode               string
	SymbolID                 string
	Precipitation            float64
	PrecipitationColor       string
	PrecipitationMax         float64
	PrecipitationMaxColor    string
	PrecipitationMin         float64
	PrecipitationMinColor    string
	PrecipitationChance      float64
	ThunderChance            float64
}

type Cache[T any] struct {
	ExpiresAt time.Time
	Data      T
}

func (y *YRImpl) getForecasts(conf config.Config) *v1.YRResponse {
	if y.CachedForecastResponse != nil && time.Now().Before(y.CachedForecastResponse.ExpiresAt) {
		slog.Info("YR: Getting cached response")
		return y.CachedForecastResponse.Data
	}

	slog.Info("YR: Getting new repsonse from met.no")
	res, err := v1.YRGetLocationForecast(conf.Weather.Lat, conf.Weather.Lon)
	if err != nil {
		// HACK: needs proper error handling
		slog.Error("error occurred when calling YR.no", "err", err)
	}
	y.CachedForecastResponse = &Cache[*v1.YRResponse]{
		ExpiresAt: time.Now().Add(y.CacheLifetime),
		Data:      res,
	}

	return res
}

func (y *YRImpl) getTemperatureColor(conf config.Config, temperature float64) string {
	if temperature == conf.Weather.Colors.TempMid {
		return conf.Weather.Colors.TempColorMid
	} else if temperature <= conf.Weather.Colors.TempMin {
		return conf.Weather.Colors.TempColorCoolCoolest
	} else if temperature >= conf.Weather.Colors.TempMax {
		return conf.Weather.Colors.TempColorHotHottest
	}

	if temperature > conf.Weather.Colors.TempMid {
		value := math.Abs(temperature / conf.Weather.Colors.TempMax)

		c, err := helpers.LerpHexString(conf.Weather.Colors.TempColorHotCoolest,
			conf.Weather.Colors.TempColorHotHottest, value)
		if err != nil {
			slog.Error("error occurred when lerping colors", "err", err)
			return conf.Weather.Colors.TempColorMid
		}

		return c

	}

	value := math.Abs(temperature / conf.Weather.Colors.TempMin)

	c, err := helpers.LerpHexString(conf.Weather.Colors.TempColorCoolHottest,
		conf.Weather.Colors.TempColorCoolCoolest, value)
	if err != nil {
		slog.Error("error occurred when lerping colors", "err", err)
		return conf.Weather.Colors.TempColorMid
	}

	return c
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
		TemperatureColor:   y.getTemperatureColor(conf, latest.Data.Instant.Details.AirTemperature),
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
			TemperatureColor:   y.getTemperatureColor(conf, item.Data.Instant.Details.AirTemperature),
			SymbolCode:         item.Data.Next6Hours.Summary.SymbolCode,
			SymbolID:           helpers.YRSymbolsID[item.Data.Next6Hours.Summary.SymbolCode],
			Precipitation:      item.Data.Next6Hours.Details.PrecipitationAmount,
			PrecipitationColor: y.getPrecipitationColorClass(conf, item.Data.Next6Hours.Details.PrecipitationAmount),
		}
		forecasts = append(forecasts, forecast)
	}

	return forecasts[0:maxLength]
}

func (y *YRImpl) GetFullForecasts(conf config.Config) []YRFullForecast {
	res := y.getForecasts(conf)

	forecasts := []YRFullForecast{}
	for _, item := range res.Properties.Timeseries {
		forecast := YRFullForecast{
			Time:                     item.Time.Local().Format("15:04") + " - " + timeanddate.GetDateStr(item.Time.Local()),
			AirPressureAtSeaLevel:    item.Data.Instant.Details.AirPressureAtSeaLevel,
			Temperature:              item.Data.Instant.Details.AirTemperature,
			TemperatureColor:         y.getTemperatureColor(conf, item.Data.Instant.Details.AirTemperature),
			TemperatureMax:           item.Data.Instant.Details.AirTemperaturePercentile90,
			TemperatureMaxColor:      y.getTemperatureColor(conf, item.Data.Instant.Details.AirTemperaturePercentile90),
			TemperatureMin:           item.Data.Instant.Details.AirTemperaturePercentile10,
			TemperatureMinColor:      y.getTemperatureColor(conf, item.Data.Instant.Details.AirTemperaturePercentile10),
			DewPointTemperature:      item.Data.Instant.Details.DewPointTemperature,
			DewPointTemperatureColor: y.getTemperatureColor(conf, item.Data.Instant.Details.DewPointTemperature),
			CloudAreaFraction:        item.Data.Instant.Details.CloudAreaFraction,
			CloudAreaFractionHigh:    item.Data.Instant.Details.CloudAreaFractionHigh,
			CloudAreaFractionLow:     item.Data.Instant.Details.CloudAreaFractionLow,
			CloudAreaFractionMedium:  item.Data.Instant.Details.CloudAreaFractionMedium,
			FogAreaFraction:          item.Data.Instant.Details.FogAreaFraction,
			RelativeHumidity:         item.Data.Instant.Details.RelativeHumidity,
			UltravioletIndexClearSky: item.Data.Instant.Details.UltravioletIndexClearSky,
			WindFromDirection:        item.Data.Instant.Details.WindFromDirection,
			WindSpeed:                item.Data.Instant.Details.WindSpeed,
			WindSpeedOfGust:          item.Data.Instant.Details.WindSpeedOfGust,
			WindSpeedPercentile10:    item.Data.Instant.Details.WindSpeedPercentile10,
			WindSpeedPercentile90:    item.Data.Instant.Details.WindSpeedPercentile90,
			SymbolName:               helpers.Capitalize(item.Data.Next1Hours.Summary.SymbolCode),
			SymbolCode:               item.Data.Next1Hours.Summary.SymbolCode,
			SymbolID:                 helpers.YRSymbolsID[item.Data.Next1Hours.Summary.SymbolCode],
			Precipitation:            item.Data.Next1Hours.Details.PrecipitationAmount,
			PrecipitationColor:       y.getPrecipitationColorClass(conf, item.Data.Next1Hours.Details.PrecipitationAmount),
			PrecipitationMax:         item.Data.Next1Hours.Details.PrecipitationAmountMax,
			PrecipitationMaxColor:    y.getPrecipitationColorClass(conf, item.Data.Next1Hours.Details.PrecipitationAmountMax),
			PrecipitationMin:         item.Data.Next1Hours.Details.PrecipitationAmountMin,
			PrecipitationMinColor:    y.getPrecipitationColorClass(conf, item.Data.Next1Hours.Details.PrecipitationAmountMin),
			PrecipitationChance:      item.Data.Next1Hours.Details.ProbabilityOfPrecipitation,
			ThunderChance:            item.Data.Next1Hours.Details.ProbabilityOfThunder,
		}
		forecasts = append(forecasts, forecast)
	}

	return forecasts
}

func New() YR {
	return &YRImpl{
		CacheLifetime:          time.Minute * 5,
		CachedForecastResponse: nil,
	}
}
