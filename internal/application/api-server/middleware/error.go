package middleware

import (
	"log/slog"
	"net/http"

	"github.com/cresendoo/decidash-backend/pkg/errorx"
	"github.com/gin-gonic/gin"
)

var (
	// Database Error
	ErrDatabase Error = NewErrorWithCode("D1", http.StatusInternalServerError, slog.LevelError)
	ErrDBCommit Error = NewErrorWithCode("D2", http.StatusInternalServerError, slog.LevelError)

	// General Error
	ErrUnauthorized   Error = NewErrorWithCode("G1", http.StatusUnauthorized, slog.LevelInfo)
	ErrBadRequest     Error = NewErrorWithCode("G2", http.StatusBadRequest, slog.LevelInfo)
	ErrNotFound       Error = NewErrorWithCode("G3", http.StatusNotFound, slog.LevelInfo)
	ErrInternalServer Error = NewErrorWithCode("G4", http.StatusInternalServerError, slog.LevelError)
)

type Error struct {
	code     string
	status   int
	logLevel slog.Level
}

func (err Error) String() string {
	return err.code
}

func (err *Error) LogValue() slog.Value {
	if err == nil {
		return slog.Value{}
	}
	return slog.GroupValue([]slog.Attr{
		slog.String("code", err.code),
		slog.Int("status", err.status),
		slog.String("log_level", err.logLevel.String()),
	}...)
}

func NewErrorWithCode(code string, status int, logLevel slog.Level) Error {
	return Error{
		code:     code,
		status:   status,
		logLevel: logLevel,
	}
}

func ErrorWithCode(c *gin.Context, err error, code Error) {
	err = errorx.WrapDepth(err, 4)
	val, ok := c.Get("errCtx")
	if ok {
		if errCtx, ok := val.(map[string]any); ok {
			err = errorx.WrapWithData(err, errCtx)
		}
	}
	_ = c.Error(err).SetMeta(code)
}

func ErrCtx(c *gin.Context) (errCtx map[string]any) {
	val, ok := c.Get("errCtx")
	if ok {
		if errCtx, ok = val.(map[string]any); ok {
			return
		}
	}
	errCtx = make(map[string]any)
	c.Set("errCtx", errCtx)
	return
}

func TestError() gin.HandlerFunc {
	return func(c *gin.Context) {
		errCtx := ErrCtx(c)
		errCtx["test key"] = "test value"
		errCtx["test2 key"] = "test2 value"

		err := errorx.New("test error")
		ErrorWithCode(c, err, ErrInternalServer)
	}
}
