package dynamodb

import (
	"context"
	"errors"
	"testing"

	mocks "github.com/alexroden/checkout-kata-go/internal/mocks/pkg/repositories"
	"github.com/alexroden/checkout-kata-go/pkg/repositories"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/stretchr/testify/suite"
)

type GetSpecialsSuite struct {
	suite.Suite
	skus  []string
	rows  []*SpecialRow
	items []map[string]types.AttributeValue
}

func (s *GetSpecialsSuite) SetupTest() {
	s.skus = []string{"A"}

	s.rows = []*SpecialRow{
		{
			Sku:      s.skus[0],
			Quantity: 3,
			Price:    28,
		},
	}

	row, err := s.rows[0].Marshal()
	s.NoError(err)

	s.items = []map[string]types.AttributeValue{row}
}

func (s *GetSpecialsSuite) DB(ctx context.Context, err error) repositories.DynamoDBAPI {
	result := &mocks.MockDynamoDBAPI{}

	result.On(
		"BatchGetItem",
		ctx,
		&dynamodb.BatchGetItemInput{
			RequestItems: map[string]types.KeysAndAttributes{
				SPECIALS_TABLE: {
					Keys: []map[string]types.AttributeValue{
						{
							"sku": &types.AttributeValueMemberS{Value: s.skus[0]},
						},
					},
				},
			},
		},
	).Return(
		&dynamodb.BatchGetItemOutput{
			Responses: map[string][]map[string]types.AttributeValue{
				SPECIALS_TABLE: s.items,
			},
		},
		err,
	).Once()

	return result
}
func (s *GetSpecialsSuite) DynamoDB(err error) Repository {
	result, e := New(&Config{
		Tables: map[string]string{
			SPECIALS_TABLE: SPECIALS_TABLE,
		},
	})
	s.NoError(e)

	result.SetDb(s.DB(result.Context(), err))

	return result
}

func (s *GetSpecialsSuite) TestSuccess() {
	res, err := s.DynamoDB(nil).GetProducts(s.skus...)
	s.NoError(err)

	s.IsType([]*ProductRow{}, res)
}

func (s *GetSpecialsSuite) TestError() {
	_, err := s.DynamoDB(errors.New("Invalid Credentials")).GetProducts(s.skus...)
	s.Error(err)
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestGetSpecialsSuite(t *testing.T) {
	suite.Run(t, new(GetSpecialsSuite))
}
