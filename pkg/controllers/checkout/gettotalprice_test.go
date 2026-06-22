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

type GetTotalPriceSuite struct {
	suite.Suite
	w   *httptest.ResponseRecorder
	ctx *gin.Context
}

func (s *GetTotalPriceSuite) SetupTest() {
	gin.SetMode(gin.TestMode)

	s.w = httptest.NewRecorder()
	req := httptest.NewRequest(
		http.MethodGet,
		"/v1/total",
		nil,
	)

	ctx, _ := gin.CreateTestContext(s.w)

	req.Header.Add("Session-Id", uuid.NilUUID)
	ctx.Request = req

	s.ctx = ctx
}

func (s *GetTotalPriceSuite) Checkout(err error) checkout.Repository {
	result := &mocks.MockRepository{}

	hasError := err != nil

	result.On(
		"SetSession",
		uuid.NilUUID,
	).Once()

	result.On(
		"GetTotalPrice",
	).Return(10).Once()

	result.On(
		"HasError",
	).Return(hasError).Once()
	if hasError {
		result.On(
			"Errors",
			0,
		).Return(err).Once()
	}

	return result
}

func (s *GetTotalPriceSuite) Test_Success() {
	controller := NewCheckoutController(s.Checkout(nil))

	controller.GetTotalPrice(s.ctx)

	s.Equal(http.StatusOK, s.w.Code)

	var resp apiResponse[models.BasketTotal]
	s.NoError(json.Unmarshal(s.w.Body.Bytes(), &resp))

	s.NotEmpty(resp.Data.Total)
}

func (s *GetTotalPriceSuite) Test_Error() {
	controller := NewCheckoutController(s.Checkout(errors.InternalServerError("something went wrong")))

	controller.GetTotalPrice(s.ctx)

	s.Equal(http.StatusInternalServerError, s.w.Code)
	s.JSONEq(`{"error":{"code":"INTERNAL_SERVER_ERROR","message":"something went wrong"}}`, s.w.Body.String())
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestGetTotalPriceSuite(t *testing.T) {
	suite.Run(t, new(GetTotalPriceSuite))
}
