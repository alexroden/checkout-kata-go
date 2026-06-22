package checkout

import (
	"testing"

	mocks "github.com/alexroden/checkout-kata-go/internal/mocks/pkg/dynamodb"
	"github.com/alexroden/checkout-kata-go/pkg/dynamodb"
	"github.com/alexroden/checkout-kata-go/pkg/errors"
	"github.com/alexroden/checkout-kata-go/pkg/uuid"
	"github.com/stretchr/testify/suite"
)

type GetTotalPriceSuite struct {
	suite.Suite
	session string
}

func (s *GetTotalPriceSuite) SetupTest() {
	s.session = uuid.NilUUID
}

func (s *GetTotalPriceSuite) DynamoDB(
	err error,
	isMockingProducts, isMockingSpecials bool,
) dynamodb.Repository {
	result := &mocks.MockRepository{}

	var e error
	if isMockingProducts {
		e = err
		err = nil
	}

	result.On(
		"GetBasketItems",
		s.session,
	).Return(
		[]*dynamodb.BasketItemRow{
			{
				Sku:      "A",
				Quantity: 3,
			},
			{
				Sku:      "B",
				Quantity: 1,
			},
		},
		err,
	).Once()

	if isMockingProducts {
		if isMockingSpecials {
			err = e
			e = nil
		}

		skus := []string{"A", "B"}
		result.On(
			"GetProducts",
			skus,
		).Return(
			[]*dynamodb.ProductRow{
				{
					Sku:       "A",
					UnitPrice: 10,
				},
				{
					Sku:       "B",
					UnitPrice: 12,
				},
			},
			e,
		).Once()

		if isMockingSpecials {
			result.On(
				"GetSpecials",
				skus,
			).Return(
				[]*dynamodb.SpecialRow{
					{
						Sku:      "A",
						Quantity: 3,
						Price:    28,
					},
				},
				err,
			).Once()
		}
	}

	return result
}

func (s *GetTotalPriceSuite) Test_Success() {
	c := New(s.DynamoDB(nil, true, true))
	c.SetSession(s.session)

	total := c.GetTotalPrice()
	s.False(c.HasError())

	s.Equal(40, total)
}

func (s *GetTotalPriceSuite) Test_BasketItems_Error() {
	err := errors.InternalServerError("something went wrong")
	c := New(s.DynamoDB(err, false, false))
	c.SetSession(s.session)

	c.GetTotalPrice()
	s.True(c.HasError())

	s.Equal(err, c.Errors(0))
}

func (s *GetTotalPriceSuite) Test_Products_Error() {
	err := errors.InternalServerError("something went wrong")
	c := New(s.DynamoDB(err, true, false))
	c.SetSession(s.session)

	c.GetTotalPrice()
	s.True(c.HasError())

	s.Equal(err, c.Errors(0))
}

func (s *GetTotalPriceSuite) Test_Specials_Error() {
	err := errors.InternalServerError("something went wrong")
	c := New(s.DynamoDB(err, true, true))
	c.SetSession(s.session)

	c.GetTotalPrice()
	s.True(c.HasError())

	s.Equal(err, c.Errors(0))
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestGetTotalPriceSuite(t *testing.T) {
	suite.Run(t, new(GetTotalPriceSuite))
}
