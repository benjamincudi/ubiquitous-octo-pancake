package handlers

import (
	"atillm/datastore"
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type withdrawDS interface {
	Withdraw(ctx context.Context, info datastore.WithdrawalInfo) error
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
		switch err.(type) {
		case datastore.ErrTooLarge, datastore.ErrDailyAmount, datastore.ErrDailyUsage, datastore.ErrATMSupply:
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		default:
			_ = c.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Oops! Something went wrong, please try again."})
		}
	} else {
		c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("We've dispatched a rat to bring your $%d.", form.Amount)})
	}
}
