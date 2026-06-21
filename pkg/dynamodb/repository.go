package dynamodb

import (
	"context"

	"github.com/alexroden/checkout-kata-go/pkg/repositories"
)

type Repository interface {
	Context() context.Context
	PutBasket(row *BasketRow) error
	SetDb(db repositories.DynamoDBAPI)
}
