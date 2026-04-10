package middleware

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/labstack/echo/v5"
)

// Logger is a middleware that logs request/response metadata:
func Logger(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c *echo.Context) error {
		start := time.Now()
		err := next(c)
		latency := time.Since(start)

		attrs := []slog.Attr{
			slog.String("method", c.Request().Method),
			slog.String("path", c.Request().URL.Path),
			slog.Duration("latency", latency),
		}

		errDto := mapError(err)
		if err != nil {
			attrs = append(attrs,
				slog.Int("error_code", errDto.Code),
				slog.Any("error", err),
			)
		}

		ctx := c.Request().Context()

		if errDto.Code == http.StatusInternalServerError {
			slog.LogAttrs(ctx, slog.LevelError, "request failed", attrs...)
		} else {
			slog.LogAttrs(ctx, slog.LevelInfo, "request handled", attrs...)
		}

		return err
	}
}
