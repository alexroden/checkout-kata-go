package checkout

import (
	"net/http"

	"github.com/alexroden/checkout-kata-go/pkg/errors"
	"github.com/alexroden/checkout-kata-go/pkg/response"
	"github.com/gin-gonic/gin"
)

func (c *CheckoutController) ScanItem(ctx *gin.Context) {
	sessionId := ctx.Request.Header.Get("Session-Id")
	if sessionId == "" {
		response.HandleErrorResponse(ctx, errors.UnprocessableEntity("missing session id"))
		return
	}

	c.srv.SetSession(sessionId)

	c.srv.ScanItem(ctx.Param("sku"))
	if c.srv.HasError() {
		response.HandleErrorResponse(ctx, c.srv.Errors(0))
		return
	}

	response.HandleResponse(ctx, http.StatusNoContent, nil, nil)
}
