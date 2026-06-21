package checkout

import (
	"github.com/alexroden/checkout-kata-go/pkg/dynamodb"
	"github.com/alexroden/checkout-kata-go/pkg/uuid"
)

func (c *Checkout) StartSession() (string, error) {
	result := uuid.New(isTesting)
	if err := c.db.PutBasket(&dynamodb.BasketRow{
		Id: result,
	}); err != nil {
		return "", err
	}

	return result, nil
}
