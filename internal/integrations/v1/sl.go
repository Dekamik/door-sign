package v1

import (
	"door-sign/internal/integrations"
	"fmt"
	"strings"
	"time"
)

type TransportMode int64

const (
	Bus TransportMode = iota
	Metro
	Train
	Ship
	Tram
)

var TransportModeString = map[TransportMode]string{
	Bus:   "bus",
	Metro: "metro",
	Train: "train",
	Ship:  "ship",
	Tram:  "tram",
}

var TransportModeMap = map[string]TransportMode{
	"bus":   Bus,
	"metro": Metro,
	"train": Train,
	"ship":  Ship,
	"tram":  Tram,
}

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

type SLDeparturesResponseItem struct {
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
}

type SLDeparturesResponse struct {
	StatusCode    int    `json:"StatusCode"`
	Message       string `json:"Message"`
	ExecutionTime int    `json:"ExecutionTime"`
	ResponseData  struct {
		LatestUpdate        string                     `json:"LatestUpdate"`
		DataAge             int                        `json:"DataAge"`
		Metros              []SLDeparturesResponseItem `json:"Metros"`
		Buses               []SLDeparturesResponseItem `json:"Buses"`
		Trains              []SLDeparturesResponseItem `json:"Trains"`
		Trams               []SLDeparturesResponseItem `json:"Trams"`
		Ships               []interface{}              `json:"Ships"`
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

type SLDeviationsResponse struct {
	StatusCode   int    `json:"StatusCode"`
	Message      string `json:"Message"`
	ResponseData []struct {
		Created      time.Time `json:"Created"`
		MainNews     bool      `json:"MainNews"`
		SortOrder    int       `json:"SortOrder"`
		Header       string    `json:"Header"`
		Details      string    `json:"Details"`
		Scope        string    `json:"Scope"`
		FromDateTime time.Time `json:"FromDateTime"`
		ToDateTime   time.Time `json:"ToDateTime"`
		Updated      time.Time `json:"Updated"`
	} `json:"ResponseData"`
}

func SLStopLookup(apiKey string, searchString string, maxResults int) (*SLStopLookupResponse, error) {
	url := fmt.Sprintf("https://api.sl.se/api2/typeahead.json?key=%s&maxresults=%d&searchstring=%s", apiKey, maxResults, searchString)
	return integrations.Get[SLStopLookupResponse](url)
}

func SLGetDepartures(apiKey string, siteId string, timeWindow int) (*SLDeparturesResponse, error) {
	url := fmt.Sprintf("http://api.sl.se/api2/realtimedeparturesV4.JSON?key=%s&siteid=%s&timewindow=%d", apiKey, siteId, timeWindow)
	return integrations.Get[SLDeparturesResponse](url)
}

type SLGetDeviationsArgs struct {
	FromDate      *time.Time
	ToDate        *time.Time
	SiteID        *int
	LineNumber    []string
	TransportMode []TransportMode
}

func SLGetDeviations(apiKey string, args SLGetDeviationsArgs) (*SLDeviationsResponse, error) {
	if args.FromDate == nil && args.ToDate == nil && args.SiteID == nil && len(args.LineNumber) == 0 && len(args.TransportMode) == 0 {
		return &SLDeviationsResponse{
			Message: "No args setup for deviations",
		}, nil
	}

	if len(args.LineNumber) > 10 {
		return nil, fmt.Errorf("received %d line numbers - only 10 line numbers permitted", len(args.LineNumber))
	}

	if !((args.FromDate != nil && args.ToDate != nil) || (args.FromDate == nil && args.ToDate == nil)) {
		return nil, fmt.Errorf("both FromDate and ToDate must be set or unset simultaneously")
	}

	var transportModes string
	for i, mode := range args.TransportMode {
		if i == 0 {
			transportModes += TransportModeString[mode]
			continue
		}
		transportModes += "," + TransportModeString[mode]
	}

	lineNumbers := strings.Join(args.LineNumber, ",")

	var siteID string
	if args.SiteID != nil {
		siteID = string(*args.SiteID)
	}

	var fromDate string
	if args.FromDate != nil {
		fromDate = args.FromDate.Local().Format("2006-01-02")
	}

	var toDate string
	if args.ToDate != nil {
		fromDate = args.ToDate.Local().Format("2006-01-02")
	}

	url := fmt.Sprintf("http://api.sl.se/api2/deviations.JSON?key=%s&transportMode=%s&lineNumber=%s&siteId=%s&fromDate=%s&toDate=%s",
		apiKey, transportModes, lineNumbers, siteID, fromDate, toDate)

	return integrations.Get[SLDeviationsResponse](url)
}
