package checkout

import (
	"testing"

	mocks "github.com/alexroden/checkout-kata-go/internal/mocks/pkg/dynamodb"
	"github.com/alexroden/checkout-kata-go/pkg/dynamodb"
	"github.com/alexroden/checkout-kata-go/pkg/errors"
	"github.com/alexroden/checkout-kata-go/pkg/uuid"
	"github.com/stretchr/testify/suite"
)

type ScanItemSuite struct {
	suite.Suite
	session string
	sku     string
}

func (s *ScanItemSuite) SetupTest() {
	isTesting = true
	s.session = uuid.New(isTesting)
	s.sku = "A"
}

func (s *ScanItemSuite) DynamoDB(err error) dynamodb.Repository {
	result := &mocks.MockRepository{}

	result.On(
		"PutBasketItem",
		&dynamodb.BasketItemRow{
			BasketId: s.session,
			Sku:      s.sku,
			Quantity: 1,
		},
	).Return(err).Once()

	return result
}

func (s *ScanItemSuite) Test_Success() {
	c := New(s.DynamoDB(nil))
	c.SetSession(s.session)

	c.ScanItem(s.sku)
	s.False(c.HasError())
}

func (s *ScanItemSuite) Test_Error() {
	err := errors.InternalServerError("something went wrong")
	c := New(s.DynamoDB(err))
	c.SetSession(s.session)

	c.ScanItem(s.sku)
	s.True(c.HasError())

	s.Equal(err, c.Errors(0))
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestScanItemSuite(t *testing.T) {
	suite.Run(t, new(ScanItemSuite))
}
