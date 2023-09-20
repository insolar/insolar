package hash

import (
	"hash"
)

type hashWrapper struct {
	hash    hash.Hash
	sumFunc func([]byte) []byte
}

func (h *hashWrapper) Write(p []byte) (n int, err error) {
	return h.hash.Write(p)
}

func (h *hashWrapper) Sum(b []byte) []byte {
	return h.hash.Sum(b)
}

func (h *hashWrapper) Reset() {
	h.hash.Reset()
}

func (h *hashWrapper) Size() int {
	return h.hash.Size()
}

func (h *hashWrapper) BlockSize() int {
	return h.hash.BlockSize()
}

func (h *hashWrapper) Hash(b []byte) []byte {
	return h.sumFunc(b)
}
