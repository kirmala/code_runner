package middleware

import (
	"fmt"

	"github.com/labstack/echo/v5"
)

type ErrFromPanic struct {
	panicValue any
}

func (e ErrFromPanic) Error() string {
	return fmt.Sprintf("panic occured: %v", e.panicValue)
}

func Recover(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c *echo.Context) (err error) {
		defer func() {
			if r := recover(); r != nil {
				err = ErrFromPanic{panicValue: r}
			}
		}()

		return next(c)
	}
}