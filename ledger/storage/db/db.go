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
	Get(Key) ([]byte, error)
	Set(Key, []byte) error
}

const (
	ScopePulse Scope = 1
)
