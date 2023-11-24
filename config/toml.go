package config

import (
	"log"
	"os"
	"github.com/BurntSushi/toml"
)

type Config struct {
	App struct {
		Port int
	}
	Departures struct {
		BusStop              string
		SLDeparturesV4Key    string
		SLServiceAlertsV2Key string
		SLStopLookupV1Key    string
		SLStopsAndLinesV2Key string
	}
	Weather struct {
		Lat float32
		Lon float32
		Colors struct {
			TempQ1 float64
			TempQ2 float64
			TempQ3 float64
			TempQ4 float64

			ClassQ1 string
			ClassQ2 string
			ClassQ3 string
			ClassQ4 string

			ClassPrecip string
			ClassNoPrecip string
		}
	}
}

func ReadConfig() *Config {
	tomlData, err := os.ReadFile("config.toml")
	if err != nil {
		log.Panicln(err)
	}
	str := string(tomlData)
	var conf Config
	_, err = toml.Decode(str, &conf)
	if err != nil {
		log.Panicln(err)
	}
	return &conf
}
