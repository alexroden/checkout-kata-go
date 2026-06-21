package controllers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	mocks "github.com/alexroden/checkout-kata-go/internal/mocks/pkg/dynamodb"
	"github.com/alexroden/checkout-kata-go/pkg/dynamodb"
	"github.com/alexroden/checkout-kata-go/pkg/errors"
	"github.com/alexroden/checkout-kata-go/pkg/models"
	"github.com/alexroden/checkout-kata-go/pkg/uuid"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

type apiResponse[T any] struct {
	Data T `json:"data"`
}

type CheckoutSuite struct {
	suite.Suite
	w   *httptest.ResponseRecorder
	ctx *gin.Context
}

func (s *CheckoutSuite) SetupTest() {
	isTesting = true
	gin.SetMode(gin.TestMode)

	s.w = httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(s.w)
	s.ctx = ctx
}

func (s *CheckoutSuite) DynamoDB(err error) dynamodb.Repository {
	result := &mocks.MockRepository{}

	result.On(
		"PutBasket",
		&dynamodb.BasketRow{
			Id: uuid.New(isTesting),
		},
	).Return(err).Once()

	return result
}

func (s *CheckoutSuite) Test_Success() {
	controller := NewCheckoutController(s.DynamoDB(nil))

	controller.StartSession(s.ctx)

	s.Equal(http.StatusCreated, s.w.Code)

	var resp apiResponse[models.Basket]
	s.NoError(json.Unmarshal(s.w.Body.Bytes(), &resp))

	s.NotEmpty(resp.Data.Id)
}

func (s *CheckoutSuite) Test_Error() {
	controller := NewCheckoutController(s.DynamoDB(errors.InternalServerError("something went wrong")))

	controller.StartSession(s.ctx)

	s.Equal(http.StatusInternalServerError, s.w.Code)
	s.JSONEq(`{"error":{"code":"INTERNAL_SERVER_ERROR","message":"something went wrong"}}`, s.w.Body.String())
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestCheckoutSuite(t *testing.T) {
	suite.Run(t, new(CheckoutSuite))
}
