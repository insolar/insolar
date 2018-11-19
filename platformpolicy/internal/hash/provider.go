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

package hash

import (
	"github.com/insolar/insolar/core"
	"golang.org/x/crypto/sha3"
)

type sha3Provider struct{}

func NewSHA3Provider() AlgorithmProvider {
	return &sha3Provider{}
}

func (*sha3Provider) Hash224bits() core.Hasher {
	return &hashWrapper{
		hash: sha3.New224(),
		sumFunc: func(b []byte) []byte {
			s := sha3.Sum224(b)
			return s[:]
		},
	}
}

func (*sha3Provider) Hash256bits() core.Hasher {
	return &hashWrapper{
		hash: sha3.New256(),
		sumFunc: func(b []byte) []byte {
			s := sha3.Sum256(b)
			return s[:]
		},
	}
}

func (*sha3Provider) Hash512bits() core.Hasher {
	return &hashWrapper{
		hash: sha3.New512(),
		sumFunc: func(b []byte) []byte {
			s := sha3.Sum512(b)
			return s[:]
		},
	}
}
