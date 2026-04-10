package middleware

import (
	"strconv"
	"time"

	"github.com/kirmala/code_runner/http_server/internal/metrics"
	"github.com/labstack/echo/v5"
)

const appName = "code_runner"

func Metrics(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c *echo.Context) error {
		start := time.Now()

		err := next(c)

		resp, _ := echo.UnwrapResponse(c.Response())

		method := c.Request().Method
		path := c.Path()
		statusGroup := resp.Status/100

		metrics.HttpRequests.WithLabelValues(
			appName,
			method,
			path,
			strconv.Itoa(statusGroup),
		).Inc()

		metrics.HttpDuration.WithLabelValues(
			appName,
			method,
			path,
		).Observe(time.Since(start).Seconds())

		return err
	}
}
