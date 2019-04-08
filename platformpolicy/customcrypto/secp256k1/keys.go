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

package secp256k1

import (
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/platformpolicy"
)

type keyProcessor struct {
}

func NewKeyProcessor() insolar.KeyProcessor {
	return nil
}

func (kp *keyProcessor) GeneratePrivateKey() (platformpolicy.PrivateKey, error) {
	return nil, nil
}

func (*keyProcessor) ExtractPublicKey(privateKey platformpolicy.PrivateKey) platformpolicy.PublicKey {
	return nil
}

func (*keyProcessor) ImportPublicKeyPEM(pemEncoded []byte) (platformpolicy.PublicKey, error) {
	return nil, nil
}

func (*keyProcessor) ImportPrivateKeyPEM(pemEncoded []byte) (platformpolicy.PrivateKey, error) {
	return nil, nil
}

func (*keyProcessor) ExportPublicKeyPEM(publicKey platformpolicy.PublicKey) ([]byte, error) {
	return nil, nil
}

func (*keyProcessor) ExportPrivateKeyPEM(privateKey platformpolicy.PrivateKey) ([]byte, error) {
	return nil, nil
}

func (kp *keyProcessor) ExportPublicKeyBinary(publicKey platformpolicy.PublicKey) ([]byte, error) {
	return nil, nil
}

func (kp *keyProcessor) ImportPublicKeyBinary(data []byte) (platformpolicy.PublicKey, error) {
	return nil, nil
}
