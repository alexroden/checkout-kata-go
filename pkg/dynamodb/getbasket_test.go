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

type GetBasketSuite struct {
	suite.Suite
	id   string
	row  *BasketRow
	item map[string]types.AttributeValue
}

func (s *GetBasketSuite) SetupTest() {
	s.id = uuid.NilUUID

	s.row = &BasketRow{
		Id: uuid.NilUUID,
	}

	row, err := s.row.Marshal()
	s.NoError(err)

	s.item = row
}

func (s *GetBasketSuite) DB(ctx context.Context, err error) repositories.DynamoDBAPI {
	result := &mocks.MockDynamoDBAPI{}

	result.On(
		"GetItem",
		ctx,
		&dynamodb.GetItemInput{
			TableName: aws.String(BASKETS_TABLE),
			Key: map[string]types.AttributeValue{
				"id": &types.AttributeValueMemberS{Value: s.id},
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
func (s *GetBasketSuite) DynamoDB(err error) Repository {
	result, e := New(&Config{
		Tables: map[string]string{
			BASKETS_TABLE: BASKETS_TABLE,
		},
	})
	s.NoError(e)

	result.SetDb(s.DB(result.Context(), err))

	return result
}

func (s *GetBasketSuite) TestSuccess() {
	res, err := s.DynamoDB(nil).GetBasket(s.id)
	s.NoError(err)

	s.IsType(&BasketRow{}, res)
}

func (s *GetBasketSuite) TestError() {
	_, err := s.DynamoDB(errors.New("Invalid Credentials")).GetBasket(s.id)
	s.Error(err)
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestGetBasketSuite(t *testing.T) {
	suite.Run(t, new(GetBasketSuite))
}
