package middleware

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/kirmala/code_runner/http_server/api"
	"github.com/kirmala/code_runner/http_server/api/dto"
	"github.com/kirmala/code_runner/http_server/repository"
	"github.com/kirmala/code_runner/http_server/service"
	"github.com/labstack/echo/v5"
)


func mapError(c *echo.Context, err error) {

	if err == nil {
		return
	}

	var (
		notFound repository.ErrNotFound
		conflict repository.ErrConflict
		badReq   api.ErrBadRequest
		unauthenticated service.ErrUnauthenticated
	)

	switch {
	case errors.As(err, &notFound):
		msg := dto.Error{Error: notFound.Error()}
		writeError(c, msg, http.StatusNotFound)
	case errors.As(err, &conflict):
		msg := dto.Error{Error: conflict.Error()}
		writeError(c, msg, http.StatusConflict)
	case errors.As(err, &badReq):
		msg := dto.Error{Error: badReq.Error()}
		writeError(c, msg, http.StatusBadRequest)
	case errors.As(err, &unauthenticated):
		msg := dto.Error{Error: "Unauthorized"}
		writeError(c, msg, http.StatusUnauthorized)
	default:
		msg := dto.Error{Error: "Internal Server Error"}
		writeError(c, msg, http.StatusInternalServerError)
		log.Printf("Internal server error: %v", err)
	}
}

func writeError(c *echo.Context, msg dto.Error, errorStatus int) {
	if errorStatus < 400 || errorStatus >= 600 {
		panic(fmt.Sprintf("invalid error status: %d", errorStatus))
	}

	err := c.JSON(errorStatus, msg)
	if err != nil {
		panic(fmt.Sprintf("failed to write error response: %v", err))
	}
}

func ServeErrors(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c *echo.Context) error {
		err := next(c)
		if err != nil {
			mapError(c, err)
		}
		return err
	}
}