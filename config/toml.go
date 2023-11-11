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
