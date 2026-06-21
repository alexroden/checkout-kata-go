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

type StartSessionSuite struct {
	suite.Suite
	w   *httptest.ResponseRecorder
	ctx *gin.Context
}

func (s *StartSessionSuite) SetupTest() {
	gin.SetMode(gin.TestMode)

	req := httptest.NewRequest(
		http.MethodPost,
		"/v1/start-session",
		nil,
	)

	ctx, _ := gin.CreateTestContext(s.w)
	ctx.Request = req

	s.ctx = ctx
}

func (s *StartSessionSuite) Checkout(err error) checkout.Repository {
	result := &mocks.MockRepository{}

	result.On(
		"StartSession",
	).Return(uuid.NilUUID, err).Once()

	return result
}

func (s *StartSessionSuite) Test_Success() {
	controller := NewCheckoutController(s.Checkout(nil))

	controller.StartSession(s.ctx)

	s.Equal(http.StatusCreated, s.w.Code)

	var resp apiResponse[models.Basket]
	s.NoError(json.Unmarshal(s.w.Body.Bytes(), &resp))

	s.NotEmpty(resp.Data.Id)
}

func (s *StartSessionSuite) Test_Error() {
	controller := NewCheckoutController(s.Checkout(errors.InternalServerError("something went wrong")))

	controller.StartSession(s.ctx)

	s.Equal(http.StatusInternalServerError, s.w.Code)
	s.JSONEq(`{"error":{"code":"INTERNAL_SERVER_ERROR","message":"something went wrong"}}`, s.w.Body.String())
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestStartSessionSuite(t *testing.T) {
	suite.Run(t, new(StartSessionSuite))
}
