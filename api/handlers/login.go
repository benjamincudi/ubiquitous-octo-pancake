package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type LoginController struct{}

func (ctrl LoginController) HandlePost(c *gin.Context) {
	var cr UserInfo
	if err := c.Bind(&cr); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if tokenString, err := getSignedJwt(c, cr); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	} else {
		c.SetCookie(authToken, tokenString, 3600, "", "", true, true)
		c.Redirect(http.StatusSeeOther, "/account")
	}
}

func (ctrl LoginController) HandleDelete(c *gin.Context) {
	c.SetCookie(authToken, "expired", -1, "", "", true, true)
	c.Redirect(http.StatusSeeOther, "/")
}
