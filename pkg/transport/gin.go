package transport

import (
	"github.com/gin-gonic/gin"
)

func GinServer(
	endpoint EndPoint,
	decode func(c *gin.Context) (interface{}, error),
	encode func(c *gin.Context, response interface{}),
	encodeError func(c *gin.Context, err error)) func(c *gin.Context) {
	return func(c *gin.Context) {
		request, err := decode(c)
		if err != nil {
			encodeError(c, err)
			return
		}

		response, err := endpoint(c.Request.Context(), request)
		if err != nil {
			encodeError(c, err)
			return
		}

		encode(c, response)
	}
}
