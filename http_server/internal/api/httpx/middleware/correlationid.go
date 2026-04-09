package middleware

import (
	"log/slog"

	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
	slogctx "github.com/veqryn/slog-context"
)

const (
	correlationHeader = "X-Correlation-ID" // X-Correlation-ID is a common header used to pass a unique identifier for each request
	correlationKey    = "correlation_id"   // correlationKey is the key used to store the correlation ID in the Echo context
)

// CorrelationID is an Echo middleware that generates a unique correlation ID for each incoming request.
func CorrelationID(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c *echo.Context) error {
		r := c.Request()
		ctx := r.Context()

		correlationID := r.Header.Get(correlationHeader)

		if correlationID == "" {
			correlationID = uuid.New().String()
		}

		ctx = slogctx.Prepend(ctx, slog.String(correlationKey, correlationID))

		r = r.WithContext(ctx)
		c.SetRequest(r)

		c.Response().Header().Set(correlationHeader, correlationID)

		return next(c)
	}

}
