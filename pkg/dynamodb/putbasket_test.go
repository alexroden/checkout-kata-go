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

type PutBasketSuite struct {
	suite.Suite
	row  *BasketRow
	item map[string]types.AttributeValue
}

func (s *PutBasketSuite) SetupTest() {
	s.row = &BasketRow{
		Id: uuid.New(true),
	}

	row, err := s.row.Marshal()
	s.NoError(err)

	s.item = row
}

func (s *PutBasketSuite) DB(ctx context.Context, err error) repositories.DynamoDBAPI {
	result := &mocks.MockDynamoDBAPI{}

	result.On(
		"PutItem",
		ctx,
		&dynamodb.PutItemInput{
			TableName: aws.String(BASKETS_TABLE),
			Item:      s.item,
		},
	).Return(nil, err).Once()

	return result
}
func (s *PutBasketSuite) DynamoDB(err error) Repository {
	result, e := New(&Config{
		Tables: map[string]string{
			BASKETS_TABLE: BASKETS_TABLE,
		},
	})
	s.NoError(e)

	result.SetDb(s.DB(result.Context(), err))

	return result
}

func (s *PutBasketSuite) TestSuccess() {
	err := s.DynamoDB(nil).PutBasket(s.row)
	s.NoError(err)
}

func (s *PutBasketSuite) TestError() {
	err := s.DynamoDB(errors.New("Invalid Credentials")).PutBasket(s.row)
	s.Error(err)
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestPutBasketSuite(t *testing.T) {
	suite.Run(t, new(PutBasketSuite))
}
