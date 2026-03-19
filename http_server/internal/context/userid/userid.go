package userid

import (
	"context"

	"github.com/google/uuid"
)

type key int

const userIDKey key = 0

func NewContext(ctx context.Context, id uuid.UUID) context.Context {
    return context.WithValue(ctx, userIDKey, id)
}

func FromContext(ctx context.Context) (uuid.UUID, bool) {
    userID, ok := ctx.Value(userIDKey).(uuid.UUID)
    return userID, ok
}
