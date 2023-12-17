package integrations

import (
	"fmt"
)

type SLStopLookupResponse struct {
	StatusCode    int         `json:"StatusCode"`
	Message       interface{} `json:"Message"`
	ExecutionTime int         `json:"ExecutionTime"`
	ResponseData  []struct {
		Name     string      `json:"Name"`
		SiteId   string      `json:"SiteId"`
		Type     string      `json:"Type"`
		X        string      `json:"X"`
		Y        string      `json:"Y"`
		Products interface{} `json:"Products"`
	} `json:"ResponseData"`
}

type SLDeparturesResponse struct {
	StatusCode    int    `json:"StatusCode"`
	Message       string `json:"Message"`
	ExecutionTime int    `json:"ExecutionTime"`
	ResponseData  struct {
		LatestUpdate string `json:"LatestUpdate"`
		DataAge      int    `json:"DataAge"`
		Metros       []struct {
			GroupOfLine          string `json:"GroupOfLine"`
			DisplayTime          string `json:"DisplayTime"`
			TransportMode        string `json:"TransportMode"`
			LineNumber           string `json:"LineNumber"`
			Destination          string `json:"Destination"`
			JourneyDirection     int    `json:"JourneyDirection"`
			StopAreaName         string `json:"StopAreaName"`
			StopAreaNumber       int    `json:"StopAreaNumber"`
			StopPointNumber      int    `json:"StopPointNumber"`
			StopPointDesignation string `json:"StopPointDesignation"`
			TimeTabledDateTime   string `json:"TimeTabledDateTime"`
			ExpectedDateTime     string `json:"ExpectedDateTime"`
			JourneyNumber        int    `json:"JourneyNumber"`
			Deviation            struct {
				Text            string `json:"Text"`
				Consequence     string `json:"Consequence"`
				ImportanceLevel int    `json:"ImportanceLevel"`
			} `json:"Deviation"`
		} `json:"Metros"`
		Buses []struct {
			GroupOfLine          interface{} `json:"GroupOfLine"`
			TransportMode        string      `json:"TransportMode"`
			LineNumber           string      `json:"LineNumber"`
			Destination          string      `json:"Destination"`
			JourneyDirection     int         `json:"JourneyDirection"`
			StopAreaName         string      `json:"StopAreaName"`
			StopAreaNumber       int         `json:"StopAreaNumber"`
			StopPointNumber      int         `json:"StopPointNumber"`
			StopPointDesignation string      `json:"StopPointDesignation"`
			TimeTabledDateTime   string      `json:"TimeTabledDateTime"`
			ExpectedDateTime     *string     `json:"ExpectedDateTime"`
			DisplayTime          string      `json:"DisplayTime"`
			JourneyNumber        int         `json:"JourneyNumber"`
			Deviation            struct {
				Text            string `json:"Text"`
				Consequence     string `json:"Consequence"`
				ImportanceLevel int    `json:"ImportanceLevel"`
			} `json:"Deviation"`
		} `json:"Buses"`
		Trains []struct {
			SecondaryDestinationName interface{} `json:"SecondaryDestinationName"`
			GroupOfLine              string      `json:"GroupOfLine"`
			TransportMode            string      `json:"TransportMode"`
			LineNumber               string      `json:"LineNumber"`
			Destination              string      `json:"Destination"`
			JourneyDirection         int         `json:"JourneyDirection"`
			StopAreaName             string      `json:"StopAreaName"`
			StopAreaNumber           int         `json:"StopAreaNumber"`
			StopPointNumber          int         `json:"StopPointNumber"`
			StopPointDesignation     string      `json:"StopPointDesignation"`
			TimeTabledDateTime       string      `json:"TimeTabledDateTime"`
			ExpectedDateTime         string      `json:"ExpectedDateTime"`
			DisplayTime              string      `json:"DisplayTime"`
			JourneyNumber            int         `json:"JourneyNumber"`
			Deviations               []struct {
				Text            string `json:"Text"`
				Consequence     string `json:"Consequence"`
				ImportanceLevel int    `json:"ImportanceLevel"`
			} `json:"Deviations"`
		} `json:"Trains"`
		Trams []struct {
			TransportMode        string      `json:"TransportMode"`
			LineNumber           string      `json:"LineNumber"`
			Destination          string      `json:"Destination"`
			JourneyDirection     int         `json:"JourneyDirection"`
			GroupOfLine          string      `json:"GroupOfLine"`
			StopAreaName         string      `json:"StopAreaName"`
			StopAreaNumber       int         `json:"StopAreaNumber"`
			StopPointNumber      int         `json:"StopPointNumber"`
			StopPointDesignation interface{} `json:"StopPointDesignation"`
			TimeTabledDateTime   string      `json:"TimeTabledDateTime"`
			ExpectedDateTime     string      `json:"ExpectedDateTime"`
			DisplayTime          string      `json:"DisplayTime"`
			JourneyNumber        int         `json:"JourneyNumber"`
			Deviation            struct {
				Text            string `json:"Text"`
				Consequence     string `json:"Consequence"`
				ImportanceLevel int    `json:"ImportanceLevel"`
			} `json:"Deviation"`
		} `json:"Trams"`
		Ships               []interface{} `json:"Ships"`
		StopPointDeviations []struct {
			StopInfo struct {
				StopAreaNumber int    `json:"StopAreaNumber"`
				StopAreaName   string `json:"StopAreaName"`
				TransportMode  string `json:"TransportMode"`
				GroupOfLine    string `json:"GroupOfLine"`
			} `json:"StopInfo"`
			Deviation struct {
				Text            string `json:"Text"`
				Consequence     string `json:"Consequence"`
				ImportanceLevel int    `json:"ImportanceLevel"`
			} `json:"Deviation"`
		} `json:"StopPointDeviations"`
	} `json:"ResponseData"`
}

func SLStopLookup(apiKey string, searchString string, maxResults int) (*SLStopLookupResponse, error) {
	url := fmt.Sprintf("https://api.sl.se/api2/typeahead.json?key=%s&maxresults=%d&searchstring=%s", apiKey, maxResults, searchString)
	return get[SLStopLookupResponse](url)
}

func SLGetDepartures(apiKey string, siteId string, timeWindow int) (*SLDeparturesResponse, error) {
	url := fmt.Sprintf("http://api.sl.se/api2/realtimedeparturesV4.JSON?key=%s&siteid=%s&timewindow=%d", apiKey, siteId, timeWindow)
	return get[SLDeparturesResponse](url)
}
