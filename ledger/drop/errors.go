package drop

import (
	"github.com/pkg/errors"
)

var (
	// ErrNotFound is returned when value was not found.
	ErrNotFound = errors.New("value not found")

	// ErrOverride is returned if something tries to update existing record.
	ErrOverride = errors.New("records override is forbidden")
)
