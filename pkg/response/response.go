package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandleResponse(
	c *gin.Context,
	statusCode int,
	data any,
	headers http.Header,
) {
	for k, v := range headers {
		c.Header(k, v[0])
	}

	c.JSON(statusCode, gin.H{
		"data": data,
	})
}
