package main

import (
	"door-sign/integrations"
	"door-sign/templates"
	"embed"
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
	"html/template"
	"log"
	"os"
)

type config struct {
	Weather struct {
		Lat float32 `yaml:"lat"`
		Lon float32 `yaml:"lon"`
	} `yaml:"weather"`
}

func (conf *config) readConf() *config {
	yamlFile, err := os.ReadFile("config.yml")
	if err != nil {
		log.Panicln(err)
	}
	err = yaml.Unmarshal(yamlFile, conf)
	if err != nil {
		log.Panicln(err)
	}
	return conf
}

//go:embed templates
var templateFS embed.FS

func main() {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		t := template.Must(template.ParseFS(templateFS,
			"templates/index.html",
			"templates/htmx_sl.html",
			"templates/htmx_yr.html"))
		t.Execute(c.Writer, gin.H{"yr": updateYR(4)})
	})

	r.GET("/htmx/sl.html", func(c *gin.Context) {
		t := template.Must(template.ParseFS(templateFS,
			"templates/htmx_sl.html"))
		t.Execute(c.Writer, gin.H{"sl": nil})
	})

	r.GET("/htmx/yr.html", func(c *gin.Context) {
		t := template.Must(template.ParseFS(templateFS,
			"templates/htmx_yr.html"))
		t.Execute(c.Writer, gin.H{"yr": updateYR(4)})
	})

	r.Run()
}

func updateYR(maxLength int) []templates.YRForecast {
	var conf config
	conf.readConf()

	res := integrations.YRGet(conf.Weather.Lat, conf.Weather.Lon)

	forecasts := make([]templates.YRForecast, 0)
	for i, item := range res.Properties.Timeseries {
		time := item.Time.Local().Format("15:04")
		if i != 0 &&
			time != "00:00" &&
			time != "08:00" &&
			time != "12:00" &&
			time != "18:00" {
			continue
		}
		forecast := templates.YRForecast{
			Time:          time,
			Temperature:   item.Data.Instant.Details.AirTemperature,
			Symbol:        item.Data.Next6Hours.Summary.SymbolCode,
			Precipitation: item.Data.Next6Hours.Details.PrecipitationAmount,
		}
		forecasts = append(forecasts, forecast)
	}

	return forecasts[0:maxLength]
}
