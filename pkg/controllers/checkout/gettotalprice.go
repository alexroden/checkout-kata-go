package checkout

import (
	"net/http"

	"github.com/alexroden/checkout-kata-go/pkg/errors"
	"github.com/alexroden/checkout-kata-go/pkg/models"
	"github.com/alexroden/checkout-kata-go/pkg/response"
	"github.com/gin-gonic/gin"
)

func (c *CheckoutController) GetTotalPrice(ctx *gin.Context) {
	sessionId := ctx.Request.Header.Get("Session-Id")
	if sessionId == "" {
		response.HandleErrorResponse(ctx, errors.UnprocessableEntity("missing session id"))
		return
	}

	c.srv.SetSession(sessionId)

	total := c.srv.GetTotalPrice()
	if c.srv.HasError() {
		response.HandleErrorResponse(ctx, c.srv.Errors(0))
		return
	}

	response.HandleResponse(ctx, http.StatusOK, &models.BasketTotal{
		Total: total,
	}, nil)
}
