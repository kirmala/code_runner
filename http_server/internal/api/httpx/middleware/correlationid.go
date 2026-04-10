package middleware

import (
	"log/slog"

	"github.com/google/uuid"
	"github.com/kirmala/code_runner/http_server/pkg/correlationid"
	"github.com/labstack/echo/v5"
	slogctx "github.com/veqryn/slog-context"
)

const (
	correlationHeader = "X-Correlation-ID" // X-Correlation-ID is a common header used to pass a unique identifier for each request
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

		ctx = slogctx.Prepend(ctx, slog.String("correlation_id", correlationID))

		ctx = correlationid.NewContext(ctx, correlationID)

		r = r.WithContext(ctx)
		c.SetRequest(r)

		c.Response().Header().Set(correlationHeader, correlationID)

		return next(c)
	}

}
