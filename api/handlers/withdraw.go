package handlers

import (
	"atillm/datastore"
	"context"
	"encoding/base64"
	"net/http"

	"github.com/gin-gonic/gin"
)

type withdrawDS interface {
	Withdraw(ctx context.Context, info datastore.WithdrawalInfo) error
}

type AccountWithdraw struct {
	DS withdrawDS
}

func (ctrl AccountWithdraw) HandlePost(c *gin.Context) {
	var form datastore.WithdrawalInfo
	if err := c.Bind(&form); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
	}
	if user, err := userFromContext(c); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	} else {
		form.PersonalAccountNumber = user.PersonalAccountNumber
	}
	// These should really be 400 errors, but we're engaging in shenanigans to avoid using any JS
	if err := ctrl.DS.Withdraw(c, form); err != nil {
		c.SetCookie("withdrawError", err.Error(), 60, "", "", true, true)
	}
	c.Redirect(http.StatusSeeOther, "/account")
}

type QuickWithdrawController struct {
	DS withdrawDS
}

func (ctrl QuickWithdrawController) HandlePost(c *gin.Context) {
	var form datastore.WithdrawalInfo
	if err := c.Bind(&form); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err := ctrl.DS.Withdraw(c, form); err != nil {
		// Normally we would want to use AbortWithError or be returning a meaningful
		// error to the client, but since we're doing this without any client-side JS,
		// we're engaging in some cookie shenanigans to show the error
		switch err.(type) {
		case datastore.ErrTooLarge, datastore.ErrDailyAmount, datastore.ErrDailyUsage:
			c.SetCookie("withdrawError", base64.StdEncoding.EncodeToString([]byte(err.Error())), 1, "", "", true, true)
		default:
			_ = c.Error(err)
			c.SetCookie("withdrawError", base64.StdEncoding.EncodeToString([]byte("Oops, something went wrong. Please try again.")), 1, "", "", true, true)
		}
	}
	// Normally this would be more like c.JSON(http.StatusOK, someResponseStruct)
	// but since we're just using built-in form behavior, we want to redirect the
	// browser back to a friendly view instead
	c.Redirect(http.StatusSeeOther, "/")
}
