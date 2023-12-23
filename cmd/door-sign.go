package main

import (
	"door-sign/internal/config"
	"door-sign/internal/handlers/sl"
	"door-sign/internal/handlers/timeanddate"
	"door-sign/internal/handlers/yr"
	"door-sign/web"
	"html/template"

	"github.com/gin-gonic/gin"
)

func main() {
	const rowCount int = 4

	conf := *config.ReadConfig()
	siteID := sl.GetSLSiteID(conf)
	yrHandler := yr.New()

	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		t := template.Must(template.ParseFS(web.TemplateFS,
			"templates/index.html",
			"templates/imports.html",
			"templates/htmx_navbar.html",
			"templates/htmx_time.html",
			"templates/htmx_yr_now.html",
			"templates/htmx_sl.html",
			"templates/htmx_yr_forecast.html"))
		t.Execute(c.Writer, gin.H{
			"nav":   "index",
			"time":  timeanddate.GetTime(),
			"yr":    yrHandler.GetForecasts(conf, rowCount),
			"yrNow": yrHandler.GetCurrent(conf),
			"sl":    sl.GetDepartures(conf, siteID, rowCount),
		})
	})

	router.GET("/htmx/time.html", func(c *gin.Context) {
		t := template.Must(template.ParseFS(web.TemplateFS,
			"templates/htmx_time.html"))
		t.Execute(c.Writer, gin.H{"time": timeanddate.GetTime()})
	})

	router.GET("/htmx/yr_now.html", func(c *gin.Context) {
		t := template.Must(template.ParseFS(web.TemplateFS,
			"templates/htmx_yr_now.html"))
		t.Execute(c.Writer, gin.H{"yrNow": yrHandler.GetCurrent(conf)})
	})

	router.GET("/htmx/sl.html", func(c *gin.Context) {
		t := template.Must(template.ParseFS(web.TemplateFS,
			"templates/htmx_sl.html"))
		t.Execute(c.Writer, gin.H{"sl": sl.GetDepartures(conf, siteID, rowCount)})
	})

	router.GET("/htmx/yr_forecast.html", func(c *gin.Context) {
		t := template.Must(template.ParseFS(web.TemplateFS,
			"templates/htmx_yr_forecast.html"))
		t.Execute(c.Writer, gin.H{"yr": yrHandler.GetForecasts(conf, rowCount)})
	})

	router.GET("/weather", func(c *gin.Context) {
		t := template.Must(template.ParseFS(web.TemplateFS,
			"templates/weather.html",
			"templates/imports.html",
			"templates/htmx_navbar.html",
			"templates/htmx_yr_full_forecast.html"))
		t.Execute(c.Writer, gin.H{
			"nav": "weather",
			"yr":  yrHandler.GetFullForecasts(conf),
		})
	})

	router.GET("/disruptions", func(c *gin.Context) {
		t := template.Must(template.ParseFS(web.TemplateFS,
			"templates/disruptions.html",
			"templates/imports.html",
			"templates/htmx_navbar.html",
			"templates/htmx_sl_deviations.html"))
		t.Execute(c.Writer, gin.H{
			"nav": "disruptions",
			"sl":  sl.GetDeviations(conf),
		})
	})

	router.Static("/css", "./web/css")
	router.Static("/images", "./web/images")

	router.Run()
}
