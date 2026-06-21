package checkout

import (
	"testing"

	mocks "github.com/alexroden/checkout-kata-go/internal/mocks/pkg/dynamodb"
	"github.com/alexroden/checkout-kata-go/pkg/dynamodb"
	"github.com/alexroden/checkout-kata-go/pkg/errors"
	"github.com/alexroden/checkout-kata-go/pkg/uuid"
	"github.com/stretchr/testify/suite"
)

type StartSessionSuite struct {
	suite.Suite
}

func (s *StartSessionSuite) SetupTest() {
	isTesting = true
}

func (s *StartSessionSuite) DynamoDB(err error) dynamodb.Repository {
	result := &mocks.MockRepository{}

	result.On(
		"PutBasket",
		&dynamodb.BasketRow{
			Id: uuid.New(isTesting),
		},
	).Return(err).Once()

	return result
}

func (s *StartSessionSuite) Test_Success() {
	c := New(s.DynamoDB(nil))

	id, err := c.StartSession()
	s.NoError(err)

	s.Equal(uuid.New(isTesting), id)
}

func (s *StartSessionSuite) Test_Error() {
	c := New(s.DynamoDB(errors.InternalServerError("something went wrong")))

	_, err := c.StartSession()
	s.Error(err)
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestStartSessionSuite(t *testing.T) {
	suite.Run(t, new(StartSessionSuite))
}
