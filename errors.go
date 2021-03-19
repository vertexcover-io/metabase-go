package metabase_client

import (
	"fmt"

	"github.com/pkg/errors"
)

type ErrorCode int

const (
	UnknownError        ErrorCode = -99
	NetworkError        ErrorCode = 1
	BadRequest          ErrorCode = 400
	Unauthorized        ErrorCode = 401
	PermissionDenied    ErrorCode = 403
	NotFound            ErrorCode = 404
	InternalServerError ErrorCode = 500
	ServiceUnavialble   ErrorCode = 503
)

type APIError struct {
	error
	context map[string]string
	code    ErrorCode
}

func (e *APIError) Error() string {
	return fmt.Sprintf("Code: %d.Error: %s", e.code, e.error)
}

func (e *APIError) Cause() error {
	return errors.Cause(e.error)
}

func (e *APIError) Format(s fmt.State, verb rune) {
	var f string

	switch verb {
	case 'v':
		if s.Flag('+') {
			f = "%+v"
		} else {
			f = "%v"
		}
	default:
		f = fmt.Sprintf("%c", verb)
	}
	fmt.Fprintf(s, fmt.Sprintf("Code :%%d\nError:%s\nContext: %%+v", f), verb)
}

func NewAPIError(err error, code int, msg string) *APIError {
	return &APIError{
		error: errors.Wrap(err, msg),
		code:  ErrorCode(code),
	}
}

func (e *APIError) WithContext(ctx map[string]string) *APIError {
	e.context = ctx
	return e
}

func (e *APIError) WithField(field string, value string) *APIError {
	e.context[field] = value
	return e
}
