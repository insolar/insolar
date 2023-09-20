package merkle

import (
	"hash"

	"github.com/onrik/gomerkle"
	"github.com/pkg/errors"
)

type tree interface {
	Root() []byte
}

func treeFromHashList(list [][]byte, hasher hash.Hash) (tree, error) {
	mt := gomerkle.NewTree(hasher)
	mt.AddHash(list...)

	if err := mt.Generate(); err != nil {
		return nil, errors.Wrap(err, "[ treeFromHashList ] Failed to generate merkle tree")
	}

	return &mt, nil
}
