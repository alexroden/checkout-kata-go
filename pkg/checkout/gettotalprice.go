package checkout

import (
	"github.com/alexroden/checkout-kata-go/pkg/dynamodb"
	"github.com/alexroden/checkout-kata-go/pkg/errors"
)

type Pricing struct {
	*dynamodb.ProductRow
	*dynamodb.SpecialRow
}

func (c *Checkout) GetTotalPrice() int {
	if c.session == "" {
		c.errors = append(c.errors, errors.BadRequest("session not set"))

		return 0
	}

	items, err := c.db.GetBasketItems(c.session)
	if err != nil {
		c.errors = append(c.errors, err)

		return 0
	}

	skus := []string{}
	for _, item := range items {
		skus = append(skus, item.Sku)
	}

	products, err := c.db.GetProducts(skus...)
	if err != nil {
		c.errors = append(c.errors, err)

		return 0
	}

	specials, err := c.db.GetSpecials(skus...)
	if err != nil {
		c.errors = append(c.errors, err)

		return 0
	}

	pricing := map[string]*Pricing{}
	for _, product := range products {
		var special *dynamodb.SpecialRow
		for _, s := range specials {
			if product.Sku == s.Sku {
				special = s
				break
			}
		}

		pricing[product.Sku] = &Pricing{
			ProductRow: product,
			SpecialRow: special,
		}
	}

	var result int
	for _, item := range items {
		p := pricing[item.Sku]
		if p.SpecialRow == nil {
			result += item.Quantity * p.UnitPrice
			continue
		}

		offerCount := item.Quantity / p.SpecialRow.Quantity
		remainder := item.Quantity % p.SpecialRow.Quantity

		result += offerCount*p.SpecialRow.Price + remainder*p.UnitPrice
	}

	return result
}
