package checkout

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	mocks "github.com/alexroden/checkout-kata-go/internal/mocks/pkg/checkout"
	"github.com/alexroden/checkout-kata-go/pkg/checkout"
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
	gin.SetMode(gin.TestMode)

	s.w = httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(s.w)
	s.ctx = ctx
}

func (s *CheckoutSuite) Checkout(err error) checkout.Repository {
	result := &mocks.MockRepository{}

	result.On(
		"StartSession",
	).Return(uuid.NilUUID, err).Once()

	return result
}

func (s *CheckoutSuite) Test_StartSession_Success() {
	controller := NewCheckoutController(s.Checkout(nil))

	controller.StartSession(s.ctx)

	s.Equal(http.StatusCreated, s.w.Code)

	var resp apiResponse[models.Basket]
	s.NoError(json.Unmarshal(s.w.Body.Bytes(), &resp))

	s.NotEmpty(resp.Data.Id)
}

func (s *CheckoutSuite) Test_StartSession_Error() {
	controller := NewCheckoutController(s.Checkout(errors.InternalServerError("something went wrong")))

	controller.StartSession(s.ctx)

	s.Equal(http.StatusInternalServerError, s.w.Code)
	s.JSONEq(`{"error":{"code":"INTERNAL_SERVER_ERROR","message":"something went wrong"}}`, s.w.Body.String())
}

func (s *CheckoutSuite) Test_ScanItem_Success() {
	controller := NewCheckoutController(s.Checkout(nil))

	s.ctx.Header("Session-Id", uuid.NilUUID)

	controller.ScanItem(s.ctx)

	s.Equal(http.StatusNoContent, s.w.Code)
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestCheckoutSuite(t *testing.T) {
	suite.Run(t, new(CheckoutSuite))
}
