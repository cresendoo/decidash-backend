package apiserver

import (
	"github.com/cresendoo/decidash-backend/internal/application/api-server/middleware"
)

var (
	// response
	response      = middleware.Response
	ErrorWithCode = middleware.ErrorWithCode

	// Context
	Logger = middleware.Logger

	// Error
	ErrCtx = middleware.ErrCtx

	// Database Error
	ErrDatabase Error = middleware.ErrDatabase
	ErrDBCommit Error = middleware.ErrDBCommit

	// General Error
	ErrUnauthorized   Error = middleware.ErrUnauthorized
	ErrBadRequest     Error = middleware.ErrBadRequest
	ErrNotFound       Error = middleware.ErrNotFound
	ErrInternalServer Error = middleware.ErrInternalServer
)

type (
	Error = middleware.Error
)
