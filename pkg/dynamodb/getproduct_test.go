package dynamodb

import (
	"context"
	"errors"
	"testing"

	mocks "github.com/alexroden/checkout-kata-go/internal/mocks/pkg/repositories"
	"github.com/alexroden/checkout-kata-go/pkg/repositories"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/stretchr/testify/suite"
)

type GetProductSuite struct {
	suite.Suite
	sku  string
	row  *ProductRow
	item map[string]types.AttributeValue
}

func (s *GetProductSuite) SetupTest() {
	s.sku = "A"

	s.row = &ProductRow{
		Sku:       s.sku,
		UnitPrice: 10,
	}

	row, err := s.row.Marshal()
	s.NoError(err)

	s.item = row
}

func (s *GetProductSuite) DB(ctx context.Context, err error) repositories.DynamoDBAPI {
	result := &mocks.MockDynamoDBAPI{}

	result.On(
		"GetItem",
		ctx,
		&dynamodb.GetItemInput{
			TableName: aws.String(PRODUCTS_TABLE),
			Key: map[string]types.AttributeValue{
				"sku": &types.AttributeValueMemberS{Value: s.sku},
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
func (s *GetProductSuite) DynamoDB(err error) Repository {
	result, e := New(&Config{
		Tables: map[string]string{
			PRODUCTS_TABLE: PRODUCTS_TABLE,
		},
	})
	s.NoError(e)

	result.SetDb(s.DB(result.Context(), err))

	return result
}

func (s *GetProductSuite) TestSuccess() {
	res, err := s.DynamoDB(nil).GetProduct(s.sku)
	s.NoError(err)

	s.IsType(&ProductRow{}, res)
}

func (s *GetProductSuite) TestError() {
	_, err := s.DynamoDB(errors.New("Invalid Credentials")).GetProduct(s.sku)
	s.Error(err)
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestGetProductSuite(t *testing.T) {
	suite.Run(t, new(GetProductSuite))
}
