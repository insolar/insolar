package hack

// FIXME: exists only for temporary hack (INS-744)

import (
	"context"
)

type skipValidationKey struct{}

func SkipValidation(ctx context.Context) bool {
	val := ctx.Value(skipValidationKey{})
	if val == nil {
		return false
	}
	return val.(bool)
}

func SetSkipValidation(ctx context.Context, skip bool) context.Context {
	return context.WithValue(ctx, skipValidationKey{}, skip)
}
