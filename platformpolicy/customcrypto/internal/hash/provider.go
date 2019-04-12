//
// Copyright 2019 Insolar Technologies GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package hash

import (
	"golang.org/x/crypto/sha3"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/platformpolicy/algorithmprovider"
)

type sha3Provider struct{}

func NewSHA3Provider() algorithmprovider.HashAlgorithmProvider {
	return &sha3Provider{}
}

func (*sha3Provider) Hash224bits() insolar.Hasher {
	return &hashWrapper{
		hash: sha3.New224(),
		sumFunc: func(b []byte) []byte {
			s := sha3.Sum224(b)
			return s[:]
		},
	}
}

func (*sha3Provider) Hash256bits() insolar.Hasher {
	return &hashWrapper{
		hash: sha3.New256(),
		sumFunc: func(b []byte) []byte {
			s := sha3.Sum256(b)
			return s[:]
		},
	}
}

func (*sha3Provider) Hash512bits() insolar.Hasher {
	return &hashWrapper{
		hash: sha3.New512(),
		sumFunc: func(b []byte) []byte {
			s := sha3.Sum512(b)
			return s[:]
		},
	}
}
