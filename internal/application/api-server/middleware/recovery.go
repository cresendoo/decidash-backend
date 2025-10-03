package middleware

import (
	"errors"
	"log/slog"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"

	"github.com/cresendoo/decidash-backend/pkg/errorx"
	"github.com/gin-gonic/gin"
)

// Support : gin recovery
// See https://github.com/gin-gonic/gin/blob/a889c58de78711cb9b53de6cfcc9272c8518c729/recovery.go#L51
func GinRecovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if recoverErr := recover(); recoverErr != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				var brokenPipe bool
				if ne, ok := recoverErr.(*net.OpError); ok {
					var se *os.SyscallError
					if errors.As(ne, &se) {
						seStr := strings.ToLower(se.Error())
						if strings.Contains(seStr, "broken pipe") ||
							strings.Contains(seStr, "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				headers := strings.Split(string(httpRequest), "\r\n")
				for idx, header := range headers {
					current := strings.Split(header, ":")
					if current[0] == "Authorization" {
						headers[idx] = current[0] + ": *"
					}
				}

				l := Logger(c)

				l = l.With(
					slog.String("http_request", string(httpRequest)),
					slog.String("method", c.Request.Method),
					slog.String("url", c.Request.URL.String()),
				)
				if len(headers) > 0 {
					headersToStr := strings.Join(headers, "\r\n")
					l = l.With(slog.String("header", headersToStr))
				}

				var err error
				if handleErr, ok := recoverErr.(error); ok {
					err = errorx.Wrap(handleErr)
				} else {
					err = errorx.New("panic recovered")
				}

				if brokenPipe {
					l.Error("broken pipe", "error", err)
				} else {
					l.Error("panic recovered", "error", err)
				}

				if brokenPipe {
					// If the connection is dead, we can't write a status to it.
					c.Error(err) //nolint: errcheck
					c.Abort()
				} else {
					c.AbortWithStatus(http.StatusInternalServerError)
				}
			}
		}()
		c.Next()
	}
}
