package v1

import (
	"errors"
	"github.com/labstack/echo/v4"
)

var (
	ErrInvalidAuthHeader = errors.New("invalid auth header")
	ErrCannotParseToken  = errors.New("cannot parse token")
)

func newErrorResponse(c echo.Context, errStatus int, message string) {
	err := errors.New(message)
	var HTTPError *echo.HTTPError
	ok := errors.As(err, &HTTPError)
	if !ok {
		report := echo.NewHTTPError(errStatus, err.Error())
		_ = c.JSON(errStatus, report)
	}
	c.Error(errors.New("internal server error"))
}
