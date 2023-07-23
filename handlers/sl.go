package handlers

import (
	"door-sign/configuration"
	"door-sign/helpers"
	"door-sign/integrations"
	"log"
	"net/url"
)

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

func UpdateSL(conf configuration.Config, siteId string, maxLength int) []SLDeparture {
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
			DisplayTime:   item.DisplayTime,
		}
		departures = append(departures, departure)
	}
	for _, item := range res.ResponseData.Buses {
		departure := SLDeparture{
			TransportMode: helpers.SLTransportModeIcons[item.TransportMode],
			LineNumber:    item.LineNumber,
			Destination:   item.Destination,
			DisplayTime:   item.DisplayTime,
		}
		departures = append(departures, departure)
	}
	for _, item := range res.ResponseData.Trains {
		departure := SLDeparture{
			TransportMode: helpers.SLTransportModeIcons[item.TransportMode],
			LineNumber:    item.LineNumber,
			Destination:   item.Destination,
			DisplayTime:   item.DisplayTime,
		}
		departures = append(departures, departure)
	}
	for _, item := range res.ResponseData.Trams {
		departure := SLDeparture{
			TransportMode: helpers.SLTransportModeIcons[item.TransportMode],
			LineNumber:    item.LineNumber,
			Destination:   item.Destination,
			DisplayTime:   item.DisplayTime,
		}
		departures = append(departures, departure)
	}
	return departures[0:min(maxLength, len(departures))]
}

func GetSLSiteId(conf configuration.Config) string {
	escapedBusStop := url.QueryEscape(conf.Departures.BusStop)
	res, err := integrations.SLStopLookup(conf.Departures.SLStopLookupV1Key, escapedBusStop, 1)
	if err != nil {
		log.Fatalln(err)
	}
	return res.ResponseData[0].SiteId
}
