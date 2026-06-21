package dynamodb

import (
	"context"

	"github.com/alexroden/checkout-kata-go/pkg/repositories"
)

type Repository interface {
	Context() context.Context
	GetBasket(id string) (*BasketRow, error)
	PutBasket(row *BasketRow) error
	PutBasketItem(row *BasketItemRow) error
	SetDb(db repositories.DynamoDBAPI)
}
