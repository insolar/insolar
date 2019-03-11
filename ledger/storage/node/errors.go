package node

import (
	"github.com/pkg/errors"
)

var (
	// ErrOverride is returned when trying to set nodes for non-empty pulse.
	ErrOverride = errors.New("node override is forbidden")
	// ErrNoNodes is returned when nodes for specified criteria could not be found.
	ErrNoNodes = errors.New("matching nodes not found")
)
