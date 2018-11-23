package hack

// FIXME: exists only for temporary hack (INS-744)

import (
	"context"
)

type skipValidationKey struct{}

func SkipValidation(ctx context.Context) bool {
	if valBool, ok := ctx.Value(skipValidationKey{}).(bool); ok {
		return valBool
	}
	return false
}

func SetSkipValidation(ctx context.Context, skip bool) context.Context {
	return context.WithValue(ctx, skipValidationKey{}, skip)
}
