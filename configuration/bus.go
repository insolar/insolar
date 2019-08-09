package configuration

import (
	"time"
)

type Bus struct {
	ReplyTimeout time.Duration
}

func NewBus() Bus {
	return Bus{
		ReplyTimeout: 15 * time.Second,
	}
}
