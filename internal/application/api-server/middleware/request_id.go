package middleware

import (
	"fmt"
	"sync/atomic"

	"github.com/cresendoo/decidash-backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

const requestIDKey = "reqID"

var _reqid uint64

func RequestID(c *gin.Context) string {
	return c.GetString(requestIDKey)
}

func SetRequestID(c *gin.Context) {
	c.Set(
		requestIDKey,
		fmt.Sprintf("%s-%d", utils.ProcessUID(), nextRequestID()),
	)
	c.Next()
}

func nextRequestID() uint64 {
	return atomic.AddUint64(&_reqid, 1)
}
