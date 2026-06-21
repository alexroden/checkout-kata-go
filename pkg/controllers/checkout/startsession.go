package checkout

import (
	"net/http"

	"github.com/alexroden/checkout-kata-go/pkg/models"
	"github.com/alexroden/checkout-kata-go/pkg/response"
	"github.com/gin-gonic/gin"
)

func (c *CheckoutController) StartSession(ctx *gin.Context) {
	id, err := c.srv.StartSession()
	if err != nil {
		response.HandleErrorResponse(ctx, err)
		return
	}

	response.HandleResponse(ctx, http.StatusCreated, &models.Basket{
		Id: id,
	}, nil)
}
