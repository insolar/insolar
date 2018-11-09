/*
 *    Copyright 2018 Insolar
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package merkle

import (
	"bytes"

	"github.com/cbergoon/merkletree"
	"github.com/insolar/insolar/cryptohelpers/hash"
)

type tree interface {
	MerkleRoot() []byte
}

type treeNode struct {
	content []byte
}

func (t *treeNode) CalculateHash() ([]byte, error) {
	return hash.IntegrityHasher().Hash(t.content), nil
}

func (t *treeNode) Equals(other merkletree.Content) (bool, error) {
	equal := bytes.Equal(t.content, other.(*treeNode).content)
	return equal, nil
}

func fromList(list [][]byte) tree {
	var result []merkletree.Content

	for _, content := range list {
		result = append(result, &treeNode{content: content})
	}

	mt, err := merkletree.NewTree(result)
	if err != nil {
		panic(err.Error())
	}

	return mt
}
