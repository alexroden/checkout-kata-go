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

type GetBasketItemsSuite struct {
	suite.Suite
	basketId string
	rows     []*BasketItemRow
	items    []map[string]types.AttributeValue
}

func (s *GetBasketItemsSuite) SetupTest() {
	s.basketId = uuid.NilUUID
	s.rows = []*BasketItemRow{
		{
			BasketId: s.basketId,
			Sku:      "A",
		},
	}

	row, err := s.rows[0].Marshal()
	s.NoError(err)

	s.items = []map[string]types.AttributeValue{row}
}

func (s *GetBasketItemsSuite) DB(ctx context.Context, err error) repositories.DynamoDBAPI {
	result := &mocks.MockDynamoDBAPI{}

	result.On(
		"Query",
		ctx,
		&dynamodb.QueryInput{
			TableName:              aws.String(BASKET_ITEMS_TABLE),
			KeyConditionExpression: aws.String("basketId = :basketId"),
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":basketId": &types.AttributeValueMemberS{
					Value: s.basketId,
				},
			},
		},
	).Return(
		&dynamodb.QueryOutput{
			Items: s.items,
		},
		err,
	).Once()

	return result
}
func (s *GetBasketItemsSuite) DynamoDB(err error) Repository {
	result, e := New(&Config{
		Tables: map[string]string{
			BASKET_ITEMS_TABLE: BASKET_ITEMS_TABLE,
		},
	})
	s.NoError(e)

	result.SetDb(s.DB(result.Context(), err))

	return result
}

func (s *GetBasketItemsSuite) TestSuccess() {
	res, err := s.DynamoDB(nil).GetBasketItems(s.basketId)
	s.NoError(err)

	s.IsType([]*BasketItemRow{}, res)
}

func (s *GetBasketItemsSuite) TestError() {
	_, err := s.DynamoDB(errors.New("Invalid Credentials")).GetBasketItems(s.basketId)
	s.Error(err)
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestGetBasketItemsSuite(t *testing.T) {
	suite.Run(t, new(GetBasketItemsSuite))
}
