package main

import (
	"door-sign/config"
	"door-sign/handlers"
	"embed"
	"github.com/gin-gonic/gin"
	"html/template"
)

//go:embed templates
var templateFS embed.FS

func main() {
	conf := *config.ReadConfig()
	siteID := handlers.GetSLSiteID(conf)
	YR := handlers.NewYR()

	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		t := template.Must(template.ParseFS(templateFS,
			"templates/index.html",
			"templates/imports.html",
			"templates/htmx_navbar.html",
			"templates/htmx_time.html",
			"templates/htmx_yr_now.html",
			"templates/htmx_sl.html",
			"templates/htmx_yr_forecast.html"))
		t.Execute(c.Writer, gin.H{
			"nav":   "index",
			"time":  handlers.GetTime(),
			"yr":    YR.GetForecasts(conf, 4),
			"yrNow": YR.GetCurrent(conf),
			"sl":    handlers.UpdateSL(conf, siteID, 4),
		})
	})

	router.GET("/htmx/time.html", func(c *gin.Context) {
		t := template.Must(template.ParseFS(templateFS,
			"templates/htmx_time.html"))
		t.Execute(c.Writer, gin.H{"time": handlers.GetTime()})
	})

	router.GET("/htmx/yr_now.html", func(c *gin.Context) {
		t := template.Must(template.ParseFS(templateFS,
			"templates/htmx_yr_now.html"))
		t.Execute(c.Writer, gin.H{"yrNow": YR.GetCurrent(conf)})
	})

	router.GET("/htmx/sl.html", func(c *gin.Context) {
		t := template.Must(template.ParseFS(templateFS,
			"templates/htmx_sl.html"))
		t.Execute(c.Writer, gin.H{"sl": handlers.UpdateSL(conf, siteID, 4)})
	})

	router.GET("/htmx/yr_forecast.html", func(c *gin.Context) {
		t := template.Must(template.ParseFS(templateFS,
			"templates/htmx_yr_forecast.html"))
		t.Execute(c.Writer, gin.H{"yr": YR.GetForecasts(conf, 4)})
	})

	router.GET("/weather", func(c *gin.Context) {
		t := template.Must(template.ParseFS(templateFS,
			"templates/weather.html",
			"templates/imports.html",
			"templates/htmx_navbar.html"))
		t.Execute(c.Writer, gin.H{
			"nav": "weather",
		})
	})

	router.Static("/assets", "./assets")

	router.Run()
}
