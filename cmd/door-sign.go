package main

import (
	"door-sign/internal/config"
	"door-sign/internal/handlers"
	"door-sign/web"
	"html/template"

	"github.com/gin-gonic/gin"
)

func main() {
	const rowCount int = 4

	conf := *config.ReadConfig()
	siteID := handlers.GetSLSiteID(conf)
	YR := handlers.NewYR()

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
			"time":  handlers.GetTime(),
			"yr":    YR.GetForecasts(conf, rowCount),
			"yrNow": YR.GetCurrent(conf),
			"sl":    handlers.UpdateSL(conf, siteID, rowCount),
		})
	})

	router.GET("/htmx/time.html", func(c *gin.Context) {
		t := template.Must(template.ParseFS(web.TemplateFS,
			"templates/htmx_time.html"))
		t.Execute(c.Writer, gin.H{"time": handlers.GetTime()})
	})

	router.GET("/htmx/yr_now.html", func(c *gin.Context) {
		t := template.Must(template.ParseFS(web.TemplateFS,
			"templates/htmx_yr_now.html"))
		t.Execute(c.Writer, gin.H{"yrNow": YR.GetCurrent(conf)})
	})

	router.GET("/htmx/sl.html", func(c *gin.Context) {
		t := template.Must(template.ParseFS(web.TemplateFS,
			"templates/htmx_sl.html"))
		t.Execute(c.Writer, gin.H{"sl": handlers.UpdateSL(conf, siteID, rowCount)})
	})

	router.GET("/htmx/yr_forecast.html", func(c *gin.Context) {
		t := template.Must(template.ParseFS(web.TemplateFS,
			"templates/htmx_yr_forecast.html"))
		t.Execute(c.Writer, gin.H{"yr": YR.GetForecasts(conf, rowCount)})
	})

	router.GET("/weather", func(c *gin.Context) {
		t := template.Must(template.ParseFS(web.TemplateFS,
			"templates/weather.html",
			"templates/imports.html",
			"templates/htmx_navbar.html",
			"templates/htmx_yr_full_forecast.html"))
		t.Execute(c.Writer, gin.H{
			"nav": "weather",
			"yr": YR.GetFullForecasts(conf),
		})
	})

	router.GET("/disruptions", func(c *gin.Context) {
		t := template.Must(template.ParseFS(web.TemplateFS,
			"templates/disruptions.html",
			"templates/imports.html",
			"templates/htmx_navbar.html"))
		t.Execute(c.Writer, gin.H{
			"nav": "disruptions",
		})
	})

	router.Static("/css", "./web/css")
	router.Static("/images", "./web/images")

	router.Run()
}
