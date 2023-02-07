package web

import (
	"embed"
	"html/template"
	"io/fs"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	//go:embed templates
	htmlTemplates embed.FS
)

func AttachWebRoutes(e *gin.Engine) {
	if t, err := template.ParseFS(htmlTemplates, "templates/*.gohtml"); err != nil {
		log.Panicf("failed to parse embedded templates: %v", err)
	} else {
		e.SetHTMLTemplate(t)
	}

	if assets, err := fs.Sub(htmlTemplates, "templates/assets"); err != nil {
		log.Panic(err)
	} else {
		e.StaticFS("/assets", http.FS(assets))
	}
	log.Println("using embedded files - rebuild to get changes")

	e.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.gohtml", gin.H{})
	})
}
