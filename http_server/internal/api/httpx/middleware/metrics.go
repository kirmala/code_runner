package middleware

import (
	"github.com/kirmala/code_runner/http_server/internal/metrics"
	"github.com/labstack/echo/v5"
)

func Metrics(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c *echo.Context) error {
		metrics.HTTPRequestsTotal.Add(1)
		return next(c)
	}
}