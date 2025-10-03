package errorx

import (
	"errors"
	"fmt"
	"log/slog"
	"strconv"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New(message string) *customError {
	return &customError{
		err:   errors.New(message),
		stack: callers(3),
		data:  make(map[string]any),
	}
}

func Errorf(format string, args ...any) *customError {
	return &customError{
		err:   fmt.Errorf(format, args...),
		stack: callers(3),
		data:  make(map[string]any),
	}
}

type CustomError interface {
	error
	With(key string, value interface{}) CustomError
	WithData(map[string]any) CustomError
	Cause() error
	Unwrap() error
	MarshalLogObject(zapcore.ObjectEncoder) error
}

type customError struct {
	err   error
	stack *stack
	data  map[string]any
}

func (e *customError) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if e.err == nil {
		return nil
	}
	switch tErr := e.err.(type) {
	case *customError:
		zap.Object("error", tErr).AddTo(enc)
	case customErrors:
		zap.Array("errors", tErr).AddTo(enc)
	default:
		enc.AddString("message", e.Error())
	}
	enc.AddString("stack", fmt.Sprintf("%+v", e.stack))

	if len(e.data) > 0 {
		for k, v := range e.data {
			switch value := v.(type) {
			case interface{ Int64() int64 }:
				enc.AddString(k, strconv.FormatInt(value.Int64(), 10))
			case interface{ Int() int }:
				enc.AddInt(k, value.Int())
			case interface{ String() string }:
				enc.AddString(k, value.String())
			default:
				if err := enc.AddReflected(k, v); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (e *customError) Error() string {
	return e.err.Error()
}

func (e *customError) Cause() error {
	return e.err
}

func (e *customError) Unwrap() error {
	return e.err
}

func (e *customError) With(key string, data any) CustomError {
	e.data[key] = data
	return e
}

func (e *customError) WithData(data map[string]any) CustomError {
	for k, v := range data {
		e.data[k] = v
	}
	return e
}

func (e *customError) StackTrace() StackTrace {
	return e.stack.StackTrace()
}

func (e *customError) LogValue() slog.Value {
	if e.err == nil {
		return slog.Value{}
	}

	attrs := make([]slog.Attr, 0, len(e.data)+2)
	attrs = append(attrs, slog.String("message", e.Error()))
	attrs = append(attrs, slog.String("stack", fmt.Sprintf("%+v", e.stack)))

	for k, v := range e.data {
		if k == "user_id" {
			if userID, ok := v.(int64); ok {
				attrs = append(attrs, slog.String(k, strconv.FormatInt(userID, 10)))
				continue
			}
		}
		switch value := v.(type) {
		case interface{ Int64() int64 }:
			attrs = append(attrs, slog.String(k, strconv.FormatInt(value.Int64(), 10)))
		case interface{ Int() int }:
			attrs = append(attrs, slog.Int(k, value.Int()))
		case interface{ String() string }:
			attrs = append(attrs, slog.String(k, value.String()))
		default:
			attrs = append(attrs, slog.Any(k, v))
		}
	}

	return slog.GroupValue(attrs...)
}

type customErrors []CustomError

func (e customErrors) MarshalLogArray(enc zapcore.ArrayEncoder) error {
	for _, err := range e {
		if err == nil {
			continue
		}
		if err := enc.AppendObject(err); err != nil {
			return err
		}
	}
	return nil
}

func (e customErrors) Error() string {
	var b []byte
	for i, err := range e {
		if i > 0 {
			b = append(b, '\n')
		}
		b = append(b, err.Error()...)
	}
	return string(b)
}

func (e customErrors) Unwrap() []error {
	errors := make([]error, 0, len(e))
	for i := range e {
		errors = append(errors, e[i].Unwrap())
	}
	return errors
}
