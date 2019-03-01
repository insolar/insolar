package db

type Scope byte

func (s Scope) Bytes() []byte {
	return []byte{byte(s)}
}

type Key interface {
	Scope() Scope
	Key() []byte
}

const (
	ScopePulse Scope = 1
)
