package rabbitmq

import (
	"context"
	"log/slog"
	"time"

	"github.com/streadway/amqp"
)

// LoggerMiddleware is middleware used to log duration of request and error if it happens
func LoggerMiddleware(next Handler) Handler {
	return func(ctx context.Context, d amqp.Delivery) error {
		start := time.Now()
		slog.InfoContext(ctx, "task started")
		err := next(ctx, d)
		latency := time.Since(start)
		if err != nil {
			slog.ErrorContext(ctx, "task failed", slog.Any("error", err), slog.Duration("latency", latency))
		} else {
			slog.InfoContext(ctx, "task completed successfully", slog.Duration("latency", latency))
		}
		return err
	}
}