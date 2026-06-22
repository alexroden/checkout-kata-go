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

type GetBasketItemSuite struct {
	suite.Suite
	basketId string
	sku      string
	row      *BasketItemRow
	item     map[string]types.AttributeValue
}

func (s *GetBasketItemSuite) SetupTest() {
	s.basketId = uuid.NilUUID
	s.sku = "A"

	s.row = &BasketItemRow{
		BasketId: s.basketId,
		Sku:      s.sku,
	}

	row, err := s.row.Marshal()
	s.NoError(err)

	s.item = row
}

func (s *GetBasketItemSuite) DB(ctx context.Context, err error) repositories.DynamoDBAPI {
	result := &mocks.MockDynamoDBAPI{}

	result.On(
		"GetItem",
		ctx,
		&dynamodb.GetItemInput{
			TableName: aws.String(BASKET_ITEMS_TABLE),
			Key: map[string]types.AttributeValue{
				"basketId": &types.AttributeValueMemberS{Value: s.basketId},
				"sku":      &types.AttributeValueMemberS{Value: s.sku},
			},
		},
	).Return(
		&dynamodb.GetItemOutput{
			Item: s.item,
		},
		err,
	).Once()

	return result
}
func (s *GetBasketItemSuite) DynamoDB(err error) Repository {
	result, e := New(&Config{
		Tables: map[string]string{
			BASKET_ITEMS_TABLE: BASKET_ITEMS_TABLE,
		},
	})
	s.NoError(e)

	result.SetDb(s.DB(result.Context(), err))

	return result
}

func (s *GetBasketItemSuite) TestSuccess() {
	res, err := s.DynamoDB(nil).GetBasketItem(s.basketId, s.sku)
	s.NoError(err)

	s.IsType(&BasketItemRow{}, res)
}

func (s *GetBasketItemSuite) TestError() {
	_, err := s.DynamoDB(errors.New("Invalid Credentials")).GetBasketItem(s.basketId, s.sku)
	s.Error(err)
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestGetBasketItemSuite(t *testing.T) {
	suite.Run(t, new(GetBasketItemSuite))
}
