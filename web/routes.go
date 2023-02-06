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
		res, err := getWithdrawRes(c)
		if err != nil && !errors.Is(err, http.ErrNoCookie) {
			c.SetCookie("withdrawError", "", 0, "", "", true, true)
			c.SetCookie("withdrawSuccess", "", 0, "", "", true, true)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		c.HTML(http.StatusOK, "index.gohtml", gin.H{
			"error":   res.errMsg,
			"success": res.success,
		})
	})
}

type withdrawResult struct {
	errMsg, success string
}

func getWithdrawRes(c *gin.Context) (withdrawResult, error) {
	var errMsg, success string
	eCookie, err := c.Request.Cookie("withdrawError")
	switch err {
	case http.ErrNoCookie:
	case nil:
		msg, decodeErr := base64.StdEncoding.DecodeString(eCookie.Value)
		if decodeErr != nil {
			log.Printf("b64 decode error: %v\n", decodeErr)
		}
		errMsg = string(msg)
	default:
		return withdrawResult{}, err
	}

	eCookie, err = c.Request.Cookie("withdrawSuccess")
	switch err {
	case http.ErrNoCookie:
	case nil:
		msg, decodeErr := base64.StdEncoding.DecodeString(eCookie.Value)
		if decodeErr != nil {
			log.Printf("b64 decode error: %v\n", decodeErr)
		}
		success = string(msg)
	default:
		return withdrawResult{}, err
	}
	return withdrawResult{
		errMsg,
		success,
	}, err
}
