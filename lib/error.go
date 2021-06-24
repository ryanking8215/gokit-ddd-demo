package lib

import "fmt"

const (
	ErrInvalidArgument   = 1
	ErrNotFound          = 2
	ErrAlreadyExists     = 3
	ErrPermissionDenied  = 4
	ErrResourceExhausted = 5
	ErrUnavailable       = 6
	ErrDataLoss          = 7
	ErrBusy              = 8
	ErrAborted           = 9
	ErrInternal          = 10
)

var errorStrings = []string{"", "invalid argument", "not found", "aleady exists", "permission denied", "resource exhausted", "unavailable", "data loss", "busy", "aborted", "internal"}

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"msg"`
}

func NewError(c int, m string) Error {
	if m == "" {
		m = errorStrings[c]
	}
	return Error{c, m}
}

func (e Error) Error() string {
	return fmt.Sprintf("<Err(%d,%s)>", e.Code, e.Message)
}
