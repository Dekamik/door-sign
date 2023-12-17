package handlers

import (
	"door-sign/internal/config"
	"door-sign/internal/helpers"
	"door-sign/internal/integrations"
	"fmt"
	"log"
	"net/url"
	"regexp"
	"time"
)

type SLDepartures struct {
	ErrorMessage string
	Departures   []SLDeparture
}

type SLDeparture struct {
	TransportMode string
	LineNumber    string
	Destination   string
	DisplayTime   string
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
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

		localTime := time.Local()
		mins := localTime.Sub(now).Minutes()

		if mins <= 1 {
			result = "1 min"
		} else if mins <= 30 {
			result = fmt.Sprintf("%.0f min", mins)
		}
	}

	return result
}

func UpdateSL(conf config.Config, siteId string, maxLength int) SLDepartures {
	res, err := integrations.SLGetDepartures(conf.Departures.SLDeparturesV4Key, siteId, 60)
	if err != nil {
		log.Fatalln(err)
	}

	departures := make([]SLDeparture, 0)

	for _, item := range res.ResponseData.Metros {
		departure := SLDeparture{
			TransportMode: helpers.SLTransportModeIcons[item.TransportMode],
			LineNumber:    item.LineNumber,
			Destination:   item.Destination,
			DisplayTime:   formatDisplayTime(item.DisplayTime, item.ExpectedDateTime),
		}
		departures = append(departures, departure)
	}

	for _, item := range res.ResponseData.Buses {
		departure := SLDeparture{
			TransportMode: helpers.SLTransportModeIcons[item.TransportMode],
			LineNumber:    item.LineNumber,
			Destination:   item.Destination,
			DisplayTime:   formatDisplayTime(item.DisplayTime, *item.ExpectedDateTime),
		}
		departures = append(departures, departure)
	}

	for _, item := range res.ResponseData.Trains {
		departure := SLDeparture{
			TransportMode: helpers.SLTransportModeIcons[item.TransportMode],
			LineNumber:    item.LineNumber,
			Destination:   item.Destination,
			DisplayTime:   formatDisplayTime(item.DisplayTime, item.ExpectedDateTime),
		}
		departures = append(departures, departure)
	}

	for _, item := range res.ResponseData.Trams {
		departure := SLDeparture{
			TransportMode: helpers.SLTransportModeIcons[item.TransportMode],
			LineNumber:    item.LineNumber,
			Destination:   item.Destination,
			DisplayTime:   formatDisplayTime(item.DisplayTime, item.ExpectedDateTime),
		}
		departures = append(departures, departure)
	}

	message := res.Message
	if message == "" {
		message = "No data"
	}

	return SLDepartures{
		ErrorMessage: message,
		Departures: departures[0:min(maxLength, len(departures))],
	}
}

func GetSLSiteID(conf config.Config) string {
	escapedBusStop := url.QueryEscape(conf.Departures.BusStop)
	res, err := integrations.SLStopLookup(conf.Departures.SLStopLookupV1Key, escapedBusStop, 1)
	if err != nil {
		log.Fatalln(err)
	}
	return res.ResponseData[0].SiteId
}
