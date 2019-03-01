package db

type Scope byte

func (s Scope) Bytes() []byte {
	return []byte{byte(s)}
}

type Key interface {
	Scope() Scope
	Key() []byte
}

type DB interface {
	Get(key Key) (value []byte, err error)
	Set(key Key, value []byte) error
}

const (
	ScopePulse Scope = 1
)
