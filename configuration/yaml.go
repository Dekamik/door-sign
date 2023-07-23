package configuration

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

type Config struct {
	App struct {
		Port int `yaml:"port"`
	}
	Departures struct {
		BusStop              string `yaml:"bus_stop"`
		SLDeparturesV4Key    string `yaml:"sl_departures_v4_key"`
		SLServiceAlertsV2Key string `yaml:"sl_service_alerts_v2_key"`
		SLStopLookupV1Key    string `yaml:"sl_stop_lookup_v1_key"`
		SLStopsAndLinesV2Key string `yaml:"sl_stops_and_lines_v2_key"`
	}
	Weather struct {
		Lat float32 `yaml:"lat"`
		Lon float32 `yaml:"lon"`
	} `yaml:"weather"`
}

func (c *Config) ReadYamlConfig() *Config {
	yamlFile, err := os.ReadFile("config.yml")
	if err != nil {
		log.Panicln(err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Panicln(err)
	}
	return c
}
