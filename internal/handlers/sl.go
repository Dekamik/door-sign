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

func formatDisplayTime(original string) string {
	result := original
	re := regexp.MustCompile(`\d\d:\d\d`)

	if re.MatchString(original) {
		now, err := time.Parse("15:04", time.Now().Format("15:04"))
		if err != nil {
			return "err 1"
		}

		time, err := time.Parse("15:04", original)
		if err != nil {
			return "err 2"
		}

		mins := time.Sub(now).Minutes()

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
			DisplayTime:   formatDisplayTime(item.DisplayTime),
		}
		departures = append(departures, departure)
	}

	for _, item := range res.ResponseData.Buses {
		departure := SLDeparture{
			TransportMode: helpers.SLTransportModeIcons[item.TransportMode],
			LineNumber:    item.LineNumber,
			Destination:   item.Destination,
			DisplayTime:   formatDisplayTime(item.DisplayTime),
		}
		departures = append(departures, departure)
	}

	for _, item := range res.ResponseData.Trains {
		departure := SLDeparture{
			TransportMode: helpers.SLTransportModeIcons[item.TransportMode],
			LineNumber:    item.LineNumber,
			Destination:   item.Destination,
			DisplayTime:   formatDisplayTime(item.DisplayTime),
		}
		departures = append(departures, departure)
	}

	for _, item := range res.ResponseData.Trams {
		departure := SLDeparture{
			TransportMode: helpers.SLTransportModeIcons[item.TransportMode],
			LineNumber:    item.LineNumber,
			Destination:   item.Destination,
			DisplayTime:   formatDisplayTime(item.DisplayTime),
		}
		departures = append(departures, departure)
	}

	return SLDepartures{
		ErrorMessage: res.Message,
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
