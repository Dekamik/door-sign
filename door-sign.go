package main

import (
	"door-sign/configuration"
	"door-sign/handlers"
	"embed"
	"github.com/gin-gonic/gin"
	"html/template"
)

//go:embed templates
var templateFS embed.FS

func main() {
	var conf configuration.Config
	conf = *conf.ReadYamlConfig()
	siteID := handlers.GetSLSiteId(conf)

	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		t := template.Must(template.ParseFS(templateFS,
			"templates/index.html",
			"templates/htmx_sl.html",
			"templates/htmx_yr.html"))
		t.Execute(c.Writer, gin.H{
			"yr": handlers.UpdateYR(conf, 4),
			"sl": handlers.UpdateSL(conf, siteID, 4),
		})
	})

	r.GET("/htmx/sl.html", func(c *gin.Context) {
		t := template.Must(template.ParseFS(templateFS,
			"templates/htmx_sl.html"))
		t.Execute(c.Writer, gin.H{"sl": handlers.UpdateSL(conf, siteID, 4)})
	})

	r.GET("/htmx/yr.html", func(c *gin.Context) {
		t := template.Must(template.ParseFS(templateFS,
			"templates/htmx_yr.html"))
		t.Execute(c.Writer, gin.H{"yr": handlers.UpdateYR(conf, 4)})
	})

	r.Static("/assets", "./assets")

	r.Run()
}
