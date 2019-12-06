package api

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Error struct {
	HttpStatus int    `json:"-"`
	Reason     string `json:"error"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("Error<%d %s>", e.HttpStatus, e.Reason)
}

func NewInvalidParamError(reason string) *Error {
	return &Error{http.StatusBadRequest, reason}
}

func NewInternalServerError(reason string) *Error {
	return &Error{http.StatusInternalServerError, reason}
}

func NewNotFoundError(reason string) *Error {
	return &Error{http.StatusNotFound, reason}
}

func ErrorHandle(err error, c echo.Context) {
	switch e := err.(type) {
	case *Error:
		code := e.HttpStatus
		c.JSON(code, e)
	default:
		c.Echo().DefaultHTTPErrorHandler(err, c)
	}
}
