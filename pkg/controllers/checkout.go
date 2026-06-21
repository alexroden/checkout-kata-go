package controllers

import (
	"net/http"

	"github.com/alexroden/checkout-kata-go/pkg/dynamodb"
	"github.com/alexroden/checkout-kata-go/pkg/models"
	"github.com/alexroden/checkout-kata-go/pkg/response"
	"github.com/alexroden/checkout-kata-go/pkg/uuid"
	"github.com/gin-gonic/gin"
)

var isTesting bool = false

type CheckoutController struct {
	db dynamodb.Repository
}

func NewCheckoutController(db dynamodb.Repository) *CheckoutController {
	return &CheckoutController{
		db: db,
	}
}

func (c *CheckoutController) StartSession(ctx *gin.Context) {
	id := uuid.New(isTesting)

	if err := c.db.PutBasket(&dynamodb.BasketRow{
		Id: id,
	}); err != nil {
		response.HandleErrorResponse(ctx, err)
		return
	}

	response.HandleResponse(ctx, http.StatusCreated, &models.Basket{
		Id: id,
	}, nil)
}
