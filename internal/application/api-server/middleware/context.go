package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"strings"

	"github.com/gin-gonic/gin"
)

// Context Key
const (
	ctxLogger = "ctx_logger"
	ctxResult = "ctx_result"
	ctxBody   = "ctx_body"
	ctxError  = "ctx_error"
)

func SetRequsetLogger(l *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(ctxLogger, l.With(slog.String("request_id", RequestID(c))))
	}
}

func Logger(c *gin.Context) *slog.Logger {
	return c.MustGet(ctxLogger).(*slog.Logger)
}

func GetBody(c *gin.Context) *bytes.Buffer {
	cloneBody(c)
	if body, ok := c.Get(ctxBody); ok {
		buf := new(bytes.Buffer)
		_ = json.Compact(buf, body.([]byte))
		return buf
	}
	return nil
}

func cloneBody(c *gin.Context) {
	if _, ok := c.Get(ctxBody); ok {
		return
	}
	contentType := c.Request.Header.Get("Content-Type")
	if strings.Contains(contentType, "application/json") {
		body, _ := io.ReadAll(c.Request.Body)
		c.Request.Body.Close()
		c.Set(ctxBody, body)
		buf := new(bytes.Buffer)
		_ = json.Compact(buf, body)
		c.Request.Body = io.NopCloser(buf)
	}
}
