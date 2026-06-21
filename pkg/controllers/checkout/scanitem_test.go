package checkout

import (
	"net/http"
	"net/http/httptest"
	"testing"

	mocks "github.com/alexroden/checkout-kata-go/internal/mocks/pkg/checkout"
	"github.com/alexroden/checkout-kata-go/pkg/checkout"
	"github.com/alexroden/checkout-kata-go/pkg/errors"
	"github.com/alexroden/checkout-kata-go/pkg/uuid"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

type ScanItemSuite struct {
	suite.Suite
	w   *httptest.ResponseRecorder
	ctx *gin.Context
}

func (s *ScanItemSuite) SetupTest() {
	gin.SetMode(gin.TestMode)

	s.w = httptest.NewRecorder()

	req := httptest.NewRequest(
		http.MethodPost,
		"/v1/scan-item/A",
		nil,
	)

	ctx, _ := gin.CreateTestContext(s.w)

	req.Header.Add("Session-Id", uuid.NilUUID)
	ctx.Request = req
	ctx.Params = gin.Params{
		{
			Key:   "sku",
			Value: "A",
		},
	}

	s.ctx = ctx
}

func (s *ScanItemSuite) Checkout(err error) checkout.Repository {
	result := &mocks.MockRepository{}

	hasError := err != nil

	result.On(
		"ScanItem",
		"A",
	).Return(uuid.NilUUID, err).Once()

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

func (s *ScanItemSuite) Test_Success() {
	controller := NewCheckoutController(s.Checkout(nil))

	controller.ScanItem(s.ctx)

	s.Equal(http.StatusNoContent, s.w.Code)
}

func (s *ScanItemSuite) Test_Error() {
	controller := NewCheckoutController(s.Checkout(errors.InternalServerError("something went wrong")))

	controller.ScanItem(s.ctx)

	s.Equal(http.StatusInternalServerError, s.w.Code)
	s.JSONEq(`{"error":{"code":"INTERNAL_SERVER_ERROR","message":"something went wrong"}}`, s.w.Body.String())
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestScanItemSuite(t *testing.T) {
	suite.Run(t, new(ScanItemSuite))
}
