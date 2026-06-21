package checkout

import (
	"github.com/alexroden/checkout-kata-go/pkg/checkout"
)

type CheckoutController struct {
	srv checkout.Repository
}

func NewCheckoutController(srv checkout.Repository) *CheckoutController {
	return &CheckoutController{
		srv: srv,
	}
}
