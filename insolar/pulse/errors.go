package pulse

import (
	"github.com/pkg/errors"
)

var (
	// ErrNotFound is returned when value was not found.
	ErrNotFound = errors.New("pulse not found")
	// ErrBadPulse is returned when appended Pulse is less than the latest.
	ErrBadPulse = errors.New("pulse should be greater than the latest")
)
