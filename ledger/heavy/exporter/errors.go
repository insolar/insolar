package exporter

import (
	"errors"
)

var (
	ErrNilCount                = errors.New("count can't be 0")
	ErrNotFinalPulseData       = errors.New("trying to get a non-finalized pulse data")
	ErrDeprecatedClientVersion = errors.New("version of the observer is outdated, please upgrade this client")
)

var RateLimitExceededMsg = "rate limit exceeded, please retry later"
