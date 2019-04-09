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

package sign

import (
	"github.com/insolar/insolar/platformpolicy/customcrypto/secp256k1"

	"github.com/insolar/insolar/insolar"
)

type mySignerWrapper struct {
	privateKey *secp256k1.PrivateKey
	hasher     insolar.Hasher
}

func (sw *mySignerWrapper) Sign(data []byte) (*insolar.Signature, error) {
	return &insolar.Signature{}, nil
}

type myVerifyWrapper struct {
	publicKey *secp256k1.PublicKey
	hasher    insolar.Hasher
}

func (sw *myVerifyWrapper) Verify(signature insolar.Signature, data []byte) bool {
	return false
}
