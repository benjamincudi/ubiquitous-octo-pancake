package web

import (
	"embed"
	"encoding/base64"
	"errors"
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
		errMsg, err := getWithdrawError(c)
		if err != nil && !errors.Is(err, http.ErrNoCookie) {
			c.SetCookie("withdrawError", "", 0, "", "", true, true)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		c.SetCookie("withdrawError", "", 0, "", "", true, true)
		c.HTML(http.StatusOK, "index.gohtml", gin.H{
			"error": errMsg,
		})
	})
}

func getWithdrawError(c *gin.Context) (string, error) {
	eCookie, err := c.Request.Cookie("withdrawError")
	if err != nil {
		return "", err
	}
	b, err := base64.StdEncoding.DecodeString(eCookie.Value)
	if err != nil {
		log.Printf("b64 decode error: %v\n", err)
	}
	return string(b), err
}
