package dynamodb

import (
	"context"

	"github.com/alexroden/checkout-kata-go/pkg/repositories"
)

type Repository interface {
	Context() context.Context
	GetBasket(id string) (*BasketRow, error)
	GetBasketItem(basketId, sku string) (*BasketItemRow, error)
	GetBasketItems(basketId string) ([]*BasketItemRow, error)
	GetProduct(sku string) (*ProductRow, error)
	GetProducts(skus ...string) ([]*ProductRow, error)
	GetSpecials(skus ...string) ([]*SpecialRow, error)
	PutBasket(row *BasketRow) error
	PutBasketItem(row *BasketItemRow) error
	SetDb(db repositories.DynamoDBAPI)
}
