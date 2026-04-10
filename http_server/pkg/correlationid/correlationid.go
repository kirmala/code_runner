package correlationid

import "context"

type key int

// correlationIDKey is the key used to store the correlation ID in the context
const correlationIDKey key = 0

// FromContext retrives correlation id from context
func FromContext(ctx context.Context) (string, bool) {
    correaltionID, ok := ctx.Value(correlationIDKey).(string)
    return correaltionID, ok
} 

// NewContext returns new context that carries provided CorrealtionID value
func NewContext(ctx context.Context, correlationID string) context.Context {
    return context.WithValue(ctx, correlationIDKey, correlationID)
}
