package response

import (
	"net/http"

	"github.com/alexroden/checkout-kata-go/pkg/errors"
	"github.com/gin-gonic/gin"
)

func HandleErrorResponse(
	c *gin.Context,
	err error,
) {
	statusCode := http.StatusInternalServerError
	switch e := err.(type) {
	case *errors.Error:
		statusCode = e.StatusCode()
	}

	c.JSON(statusCode, gin.H{
		"error": err,
	})
}
