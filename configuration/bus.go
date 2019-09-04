package configuration

import (
	"time"
)

// Bus holds some timeout options
type Bus struct {
	ReplyTimeout time.Duration
}

func NewBus() Bus {
	return Bus{
		ReplyTimeout: 15 * time.Second,
	}
}
