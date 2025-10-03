package middleware

import (
	"log/slog"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func RequestLog(c *gin.Context) {
	start := time.Now()
	path := c.Request.URL.Path
	query := c.Request.URL.RawQuery
	c.Next()

	end := time.Now()
	latency := end.Sub(start)
	if !strings.Contains(path, "health_check") {
		fields := []slog.Attr{
			slog.String("rid", RequestID(c)),
			slog.Int("status", c.Writer.Status()),	
			slog.String("method", c.Request.Method),
			slog.String("path", path),
			slog.String("query", query),
			slog.String("ip", c.ClientIP()),
			slog.String("user-agent", c.Request.UserAgent()),
			slog.Int64("latency_ms", latency.Milliseconds()),
		}
		if deviceID := c.Request.Header.Get("x-device-id"); deviceID != "" {
			fields = append(fields, slog.String("device_id", deviceID))
		}

		if gin.IsDebugging() {
			for k, v := range c.Request.Header {
				Logger(c).Debug("HEADER", slog.Any(k, v))
			}
			if body := GetBody(c); body != nil && len(body.Bytes()) != 0 {
				if gin.IsDebugging() {
					Logger(c).Debug("BODY", slog.Any("body", body.Bytes()))
				}
			}
		}

		Logger(c).LogAttrs(c, slog.LevelInfo, path, fields...)
	}
}
