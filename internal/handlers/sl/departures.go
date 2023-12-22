package sl

import (
	"door-sign/internal/config"
	"door-sign/internal/helpers"
	"door-sign/internal/integrations"
	"fmt"
	"log"
	"log/slog"
	"net/url"
	"regexp"
	"time"
)

type Departures struct {
	ErrorMessage string
	Departures   []Departure
}

type Departure struct {
	TransportMode string
	LineNumber    string
	Destination   string
	DisplayTime   string
}

func formatDisplayTime(displayTime string, expectedAt string) string {
	result := displayTime
	re := regexp.MustCompile(`\d\d:\d\d`)

	if re.MatchString(displayTime) {
		now := time.Now()

		time, err := time.ParseInLocation("2006-01-02T15:04:05", expectedAt, time.Local)
		if err != nil {
			return "err 1"
		}

		mins := time.Local().Sub(now).Minutes()
		if mins <= 1 {
			result = "1 min"
		} else if mins <= 30 {
			result = fmt.Sprintf("%.0f min", mins)
		}
	}

	return result
}

func extractDepartures(items []integrations.SLDeparturesResponseItem) []Departure {
	departures := make([]Departure, 0)

	for _, item := range items {
		departure := Departure{
			TransportMode: helpers.SLTransportModeIcons[item.TransportMode],
			LineNumber:    item.LineNumber,
			Destination:   item.Destination,
			DisplayTime:   formatDisplayTime(item.DisplayTime, item.ExpectedDateTime),
		}
		departures = append(departures, departure)
	}

	return departures
}

func GetDepartures(conf config.Config, siteId string, maxLength int) Departures {
	res, err := integrations.SLGetDepartures(conf.SL.SLDeparturesV4Key, siteId, 60)
	if err != nil {
		slog.Error("an error occurred when calling SL API", "err", err)
		return Departures{
			ErrorMessage: err.Error(),
		}
	}

	departures := make([]Departure, 0)
	departures = append(departures, extractDepartures(res.ResponseData.Metros)...)
	departures = append(departures, extractDepartures(res.ResponseData.Buses)...)
	departures = append(departures, extractDepartures(res.ResponseData.Trains)...)
	departures = append(departures, extractDepartures(res.ResponseData.Trams)...)

	message := res.Message
	if message == "" {
		message = "No data"
	}

	length := len(departures)
	if maxLength < length {
		length = maxLength
	}

	return Departures{
		ErrorMessage: message,
		Departures:   departures[0:length],
	}
}

func GetSLSiteID(conf config.Config) string {
	escapedBusStop := url.QueryEscape(conf.SL.BusStop)
	res, err := integrations.SLStopLookup(conf.SL.SLStopLookupV1Key, escapedBusStop, 1)
	if err != nil {
		log.Fatalln(err)
	}
	return res.ResponseData[0].SiteId
}
