// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

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
