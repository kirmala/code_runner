package rabbitmq

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/streadway/amqp"
	slogctx "github.com/veqryn/slog-context"
)

const CorrelationKey = "correlation_id"

// CorrelationIdMiddleware is middleware guraties unique correlation id for each request
func CorrelationIdMiddleware(next Handler) Handler {
	return func(ctx context.Context, d amqp.Delivery) error {
		id := d.CorrelationId
		if id == "" {
			id = uuid.NewString()
		}
		ctx = slogctx.Prepend(ctx, slog.String(CorrelationKey, id))
		return next(ctx, d)
	}
}