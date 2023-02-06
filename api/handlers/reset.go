package handlers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

type resetDS interface {
	ResetDay(context.Context) error
}
type ResetController struct {
	DS resetDS
}

func (ctrl ResetController) HandlePost(c *gin.Context) {
	if err := ctrl.DS.ResetDay(c); err != nil {
		_ = c.Error(err)
	}
	c.Redirect(http.StatusSeeOther, "/")
}
