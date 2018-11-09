package messagebus

import (
	"github.com/pkg/errors"
)

var (
	// ErrNoReply is returned from player when there is no stored reply for provided message.
	ErrNoReply = errors.New("no such reply")
)
