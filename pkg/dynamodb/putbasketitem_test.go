package dynamodb

import (
	"context"
	"errors"
	"testing"

	mocks "github.com/alexroden/checkout-kata-go/internal/mocks/pkg/repositories"
	"github.com/alexroden/checkout-kata-go/pkg/repositories"
	"github.com/alexroden/checkout-kata-go/pkg/uuid"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/stretchr/testify/suite"
)

type PutBasketItemSuite struct {
	suite.Suite
	row  *BasketItemRow
	item map[string]types.AttributeValue
}

func (s *PutBasketItemSuite) SetupTest() {
	s.row = &BasketItemRow{
		BasketId: uuid.NilUUID,
		Sku:      "A",
		Quantity: 1,
	}

	row, err := s.row.Marshal()
	s.NoError(err)

	s.item = row
}

func (s *PutBasketItemSuite) DB(ctx context.Context, err error) repositories.DynamoDBAPI {
	result := &mocks.MockDynamoDBAPI{}

	result.On(
		"PutItem",
		ctx,
		&dynamodb.PutItemInput{
			TableName: aws.String(BASKET_ITEMS_TABLE),
			Item:      s.item,
		},
	).Return(nil, err).Once()

	return result
}
func (s *PutBasketItemSuite) DynamoDB(err error) Repository {
	result, e := New(&Config{
		Tables: map[string]string{
			BASKET_ITEMS_TABLE: BASKET_ITEMS_TABLE,
		},
	})
	s.NoError(e)

	result.SetDb(s.DB(result.Context(), err))

	return result
}

func (s *PutBasketItemSuite) TestSuccess() {
	err := s.DynamoDB(nil).PutBasketItem(s.row)
	s.NoError(err)
}

func (s *PutBasketItemSuite) TestError() {
	err := s.DynamoDB(errors.New("Invalid Credentials")).PutBasketItem(s.row)
	s.Error(err)
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestPutBasketItemSuite(t *testing.T) {
	suite.Run(t, new(PutBasketItemSuite))
}
