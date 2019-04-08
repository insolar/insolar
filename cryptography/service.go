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

package cryptography

import (
	"github.com/pkg/errors"

	"github.com/insolar/insolar/component"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/keystore"
	"github.com/insolar/insolar/platformpolicy"
	"github.com/insolar/insolar/platformpolicy/keys"
)

type nodeCryptographyService struct {
	KeyStore                   insolar.KeyStore                   `inject:""`
	PlatformCryptographyScheme insolar.PlatformCryptographyScheme `inject:""`
	KeyProcessor               insolar.KeyProcessor               `inject:""`
}

func (cs *nodeCryptographyService) GetPublicKey() (keys.PublicKey, error) {
	privateKey, err := cs.KeyStore.GetPrivateKey("")
	if err != nil {
		return nil, errors.Wrap(err, "[ Sign ] Failed to get private privateKey")
	}

	return cs.KeyProcessor.ExtractPublicKey(privateKey), nil
}

func (cs *nodeCryptographyService) Sign(payload []byte) (*insolar.Signature, error) {
	privateKey, err := cs.KeyStore.GetPrivateKey("")
	if err != nil {
		return nil, errors.Wrap(err, "[ Sign ] Failed to get private privateKey")
	}

	signer := cs.PlatformCryptographyScheme.Signer(privateKey)
	signature, err := signer.Sign(payload)
	if err != nil {
		return nil, errors.Wrap(err, "[ Sign ] Failed to sign payload")
	}

	return signature, nil
}

func (cs *nodeCryptographyService) Verify(publicKey keys.PublicKey, signature insolar.Signature, payload []byte) bool {
	return cs.PlatformCryptographyScheme.Verifier(publicKey).Verify(signature, payload)
}

func NewCryptographyService() insolar.CryptographyService {
	return &nodeCryptographyService{}
}

type inPlaceKeyStore struct {
	privateKey keys.PrivateKey
}

func (ipks *inPlaceKeyStore) GetPrivateKey(string) (keys.PrivateKey, error) {
	return ipks.privateKey, nil
}

func NewKeyBoundCryptographyService(privateKey keys.PrivateKey) insolar.CryptographyService {
	platformCryptographyScheme := platformpolicy.NewPlatformCryptographyScheme()
	keyStore := &inPlaceKeyStore{privateKey: privateKey}
	keyProcessor := platformpolicy.NewKeyProcessor()
	cryptographyService := NewCryptographyService()

	cm := component.Manager{}

	cm.Register(platformCryptographyScheme)
	cm.Inject(keyStore, cryptographyService, keyProcessor)
	return cryptographyService
}

func NewStorageBoundCryptographyService(path string) (insolar.CryptographyService, error) {
	platformCryptographyScheme := platformpolicy.NewPlatformCryptographyScheme()
	keyStore, err := keystore.NewKeyStore(path)
	if err != nil {
		return nil, errors.Wrap(err, "[ NewStorageBoundCryptographyService ] Failed to create KeyStore")
	}
	keyProcessor := platformpolicy.NewKeyProcessor()
	cryptographyService := NewCryptographyService()

	cm := component.Manager{}

	cm.Register(platformCryptographyScheme, keyStore)
	cm.Inject(cryptographyService, keyProcessor)
	return cryptographyService, nil
}
