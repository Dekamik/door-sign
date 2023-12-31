package config

import (
	"github.com/BurntSushi/toml"
	"os"
)

type Config struct {
	App struct {
		Port int
	}
	SL struct {
		BusStop              string
		SLDeparturesV4Key    string
		SLServiceAlertsV2Key string
		SLStopLookupV1Key    string
		SLStopsAndLinesV2Key string
		Deviations           *struct {
			SiteID         *int
			LineNumbers    []string
			TransportModes []string
		}
	}
	Weather struct {
		Lat    float32
		Lon    float32
		Colors struct {
			TempMin float64
			TempMid float64
			TempMax float64

			TempColorCoolCoolest string
			TempColorCoolHottest string
			TempColorMid         string
			TempColorHotCoolest  string
			TempColorHotHottest  string

			ClassPrecip   string
			ClassNoPrecip string
		}
	}
}

func ReadConfig() *Config {
	tomlData, err := os.ReadFile("config.toml")
	if err != nil {
		panic(err)
	}
	str := string(tomlData)
	var conf Config
	_, err = toml.Decode(str, &conf)
	if err != nil {
		panic(err)
	}
	return &conf
}
