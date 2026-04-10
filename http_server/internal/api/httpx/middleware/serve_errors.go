package middleware

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/kirmala/code_runner/http_server/internal/api"
	"github.com/kirmala/code_runner/http_server/internal/api/dto"
	"github.com/kirmala/code_runner/http_server/internal/repository"
	"github.com/kirmala/code_runner/http_server/internal/service"
	"github.com/labstack/echo/v5"
)

func mapError(err error) dto.Error {

	if err == nil {
		return dto.Error{}
	}

	var (
		notFound        repository.ErrNotFound
		conflict        repository.ErrConflict
		badReq          api.ErrBadRequest
		unauthenticated service.ErrUnauthenticated
	)

	switch {
	case errors.As(err, &notFound):
		return dto.Error{Error: notFound.Error(), Code: http.StatusNotFound}
	case errors.As(err, &conflict):
		return dto.Error{Error: conflict.Error(), Code: http.StatusConflict}
	case errors.As(err, &badReq):
		return dto.Error{Error: badReq.Error(), Code: http.StatusBadRequest}
	case errors.As(err, &unauthenticated):
		return dto.Error{Error: "Unauthorized", Code: http.StatusUnauthorized}
	default:
		return dto.Error{Error: "Internal Server Error", Code: http.StatusInternalServerError}
	}
}

func writeError(c *echo.Context, errDto dto.Error) error {
	if errDto.Code < 400 || errDto.Code >= 600 {
		panic(fmt.Sprintf("invalid error status: %d", errDto.Code))
	}

	return c.JSON(errDto.Code, errDto)
}

func ServeErrors(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c *echo.Context) error {
		err := next(c)
		if err != nil {
			errDto := mapError(err)
			return writeError(c, errDto)
		}

		return nil
	}
}
