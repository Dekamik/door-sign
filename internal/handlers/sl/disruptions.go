package sl

import (
	"door-sign/internal/config"
	"door-sign/internal/integrations/v1"
	"log/slog"
	"time"
)

type Deviations struct {
	ErrorMessage string
	Deviations   []Deviation
}

type Deviation struct {
	Header    string
	MainNews  bool
	Details   string
	Scope     string
	From      time.Time
	Until     time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

func GetDeviations(conf config.Config) Deviations {
	modes := make([]v1.TransportMode, 0)
	for _, item := range conf.SL.Deviations.TransportModes {
		mode := v1.TransportModeMap[item]
		modes = append(modes, mode)
	}

	args := v1.SLGetDeviationsArgs{
		SiteID:        conf.SL.Deviations.SiteID,
		LineNumber:    conf.SL.Deviations.LineNumbers,
		TransportMode: modes,
	}

	res, err := v1.SLGetDeviations(conf.SL.SLServiceAlertsV2Key, args)
	if err != nil {
		slog.Error("an error occurred when calling SL API", "err", err)
		return Deviations{
			ErrorMessage: err.Error(),
		}
	}

	deviations := make([]Deviation, 0)
	for _, item := range res.ResponseData {
		deviation := Deviation{
			Header:    item.Header,
			MainNews:  item.MainNews,
			Details:   item.Details,
			Scope:     item.Scope,
			From:      item.FromDateTime,
			Until:     item.ToDateTime,
			CreatedAt: item.Created,
			UpdatedAt: item.Updated,
		}
		deviations = append(deviations, deviation)
	}

	message := res.Message
	if message == "" {
		message = "No data"
	}

	return Deviations{
		ErrorMessage: message,
		Deviations:   deviations,
	}
}
