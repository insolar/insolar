/*
 *    Copyright 2019 Insolar
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

package cryptography

import (
	"crypto"

	"github.com/insolar/insolar/component"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/keystore"
	"github.com/insolar/insolar/platformpolicy"
	"github.com/pkg/errors"
)

type nodeCryptographyService struct {
	KeyStore                   core.KeyStore                   `inject:""`
	PlatformCryptographyScheme core.PlatformCryptographyScheme `inject:""`
	KeyProcessor               core.KeyProcessor               `inject:""`
}

func (cs *nodeCryptographyService) GetPublicKey() (crypto.PublicKey, error) {
	privateKey, err := cs.KeyStore.GetPrivateKey("")
	if err != nil {
		return nil, errors.Wrap(err, "[ Sign ] Failed to get private privateKey")
	}

	return cs.KeyProcessor.ExtractPublicKey(privateKey), nil
}

func (cs *nodeCryptographyService) Sign(payload []byte) (*core.Signature, error) {
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

func (cs *nodeCryptographyService) Verify(publicKey crypto.PublicKey, signature core.Signature, payload []byte) bool {
	return cs.PlatformCryptographyScheme.Verifier(publicKey).Verify(signature, payload)
}

func NewCryptographyService() core.CryptographyService {
	return &nodeCryptographyService{}
}

type inPlaceKeyStore struct {
	privateKey crypto.PrivateKey
}

func (ipks *inPlaceKeyStore) GetPrivateKey(string) (crypto.PrivateKey, error) {
	return ipks.privateKey, nil
}

func NewKeyBoundCryptographyService(privateKey crypto.PrivateKey) core.CryptographyService {
	platformCryptographyScheme := platformpolicy.NewPlatformCryptographyScheme()
	keyStore := &inPlaceKeyStore{privateKey: privateKey}
	keyProcessor := platformpolicy.NewKeyProcessor()
	cryptographyService := NewCryptographyService()

	cm := component.Manager{}

	cm.Register(platformCryptographyScheme)
	cm.Inject(keyStore, cryptographyService, keyProcessor)
	return cryptographyService
}

func NewStorageBoundCryptographyService(path string) (core.CryptographyService, error) {
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
