package checkout

import (
	"github.com/alexroden/checkout-kata-go/pkg/dynamodb"
	"github.com/alexroden/checkout-kata-go/pkg/errors"
)

func (c *Checkout) ScanItem(sku string) {
	if c.session == "" {
		c.errors = append(c.errors, errors.BadRequest("session not set"))

		return
	}

	if _, err := c.db.GetBasket(c.session); err != nil {
		c.errors = append(c.errors, err)

		return
	}

	if _, err := c.db.GetProduct(sku); err != nil {
		c.errors = append(c.errors, err)

		return
	}

	if err := c.db.PutBasketItem(&dynamodb.BasketItemRow{
		BasketId: c.session,
		Sku:      sku,
		Quantity: 1,
	}); err != nil {
		c.errors = append(c.errors, err)

		return
	}

	return
}
