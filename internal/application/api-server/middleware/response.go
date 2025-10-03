package middleware

import (
	"context"
	"io"
	"log/slog"
	"net/http"

	"github.com/cresendoo/decidash-backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

func Response(c *gin.Context, r any) {
	c.Set(ctxResult, r)
}

func ResponseHandler(c *gin.Context) {
	body := GetBody(c)
	c.Next()

	errs := c.Errors.ByType(gin.ErrorTypePrivate)
	for _, err := range errs {
		req := c.Request.Clone(context.Background())
		req.RemoteAddr = c.ClientIP()
		if body != nil {
			req.Body = io.NopCloser(body)
		}

		l := Logger(c).With(
			slog.String("host_name", utils.Hostname()),
			slog.Any("http_request", HTTPRequest{R: req}),
			slog.Any("http_error", err.Meta.(Error)),
			slog.Any("error", err.Err),
		)

		if gin.IsDebugging() {
			l.Error(err.Error())
			continue
		}

		switch err.Meta.(Error).logLevel {
		case slog.LevelDebug:
			l.Debug(err.Error())
		case slog.LevelInfo:
			l.Info(err.Error())
		case slog.LevelWarn:
			l.Warn(err.Error())
		case slog.LevelError:
			l.Error(err.Error())
		}
	}

	if len(errs) != 0 {
		err := errs[0]
		code, ok := err.Meta.(Error)
		if !ok {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		if code.code != "" {
			c.AbortWithStatusJSON(code.status, gin.H{"code": code.code})
			return
		}
		c.AbortWithStatus(code.status)
	} else {
		if r, ok := c.Get(ctxResult); ok {
			c.JSON(http.StatusOK, r)
		}
	}
}

type HTTPRequest struct {
	R *http.Request
}

func (req HTTPRequest) LogValue() slog.Value {
	if req.R == nil {
		return slog.Value{}
	}
	return slog.GroupValue(
		slog.String("method", req.R.Method),
		slog.String("url", req.R.URL.String()),
	)
}
