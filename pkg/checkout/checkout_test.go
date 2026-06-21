package checkout

import (
	"testing"

	mocks "github.com/alexroden/checkout-kata-go/internal/mocks/pkg/dynamodb"
	"github.com/alexroden/checkout-kata-go/pkg/dynamodb"
	"github.com/alexroden/checkout-kata-go/pkg/errors"
	"github.com/alexroden/checkout-kata-go/pkg/uuid"
	"github.com/stretchr/testify/suite"
)

type CheckoutSuite struct {
	suite.Suite
	session string
	errors  []error
}

func (s *CheckoutSuite) SetupTest() {
	isTesting = true
	s.session = uuid.New(isTesting)
	s.errors = append(s.errors, errors.InternalServerError("something went wrong"))
}

func (s *CheckoutSuite) DynamoDB() dynamodb.Repository {
	result := &mocks.MockRepository{}

	return result
}

func (s *CheckoutSuite) Test_New() {
	c := New(s.DynamoDB())

	s.IsType(&Checkout{}, c)
}

func (s *CheckoutSuite) Test_SetSession() {
	c := New(s.DynamoDB())
	c.SetSession(s.session)

	s.NotEmpty(c.Session(), s.session)
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestCheckoutSuite(t *testing.T) {
	suite.Run(t, new(CheckoutSuite))
}
