package db

type Scope byte

type Key interface {
	Scope() Scope
	Key() []byte
}

const (
	ScopePulse Scope = 1
)
