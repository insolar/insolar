// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package cryptography

import (
	"crypto"

	component "github.com/insolar/component-manager"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/keystore"
	"github.com/insolar/insolar/platformpolicy"
	"github.com/pkg/errors"
)

type NodeCryptographyService struct {
	KeyStore                   insolar.KeyStore                   `inject:""`
	PlatformCryptographyScheme insolar.PlatformCryptographyScheme `inject:""`
	KeyProcessor               insolar.KeyProcessor               `inject:""`
}

func (cs *NodeCryptographyService) GetPublicKey() (crypto.PublicKey, error) {
	privateKey, err := cs.KeyStore.GetPrivateKey("")
	if err != nil {
		return nil, errors.Wrap(err, "[ Sign ] Failed to get private privateKey")
	}

	return cs.KeyProcessor.ExtractPublicKey(privateKey), nil
}

func (cs *NodeCryptographyService) Sign(payload []byte) (*insolar.Signature, error) {
	privateKey, err := cs.KeyStore.GetPrivateKey("")
	if err != nil {
		return nil, errors.Wrap(err, "[ Sign ] Failed to get private privateKey")
	}

	signer := cs.PlatformCryptographyScheme.DataSigner(privateKey, cs.PlatformCryptographyScheme.IntegrityHasher())
	signature, err := signer.Sign(payload)
	if err != nil {
		return nil, errors.Wrap(err, "[ Sign ] Failed to sign payload")
	}

	return signature, nil
}

func (cs *NodeCryptographyService) Verify(publicKey crypto.PublicKey, signature insolar.Signature, payload []byte) bool {
	return cs.PlatformCryptographyScheme.DataVerifier(publicKey, cs.PlatformCryptographyScheme.IntegrityHasher()).Verify(signature, payload)
}

func NewCryptographyService() insolar.CryptographyService {
	return &NodeCryptographyService{}
}

func NewKeyBoundCryptographyService(privateKey crypto.PrivateKey) insolar.CryptographyService {
	platformCryptographyScheme := platformpolicy.NewPlatformCryptographyScheme()
	keyStore := keystore.NewInplaceKeyStore(privateKey)
	keyProcessor := platformpolicy.NewKeyProcessor()
	cryptographyService := NewCryptographyService()

	cm := component.NewManager(nil)

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

	cm := component.NewManager(nil)

	cm.Register(platformCryptographyScheme, keyStore)
	cm.Inject(cryptographyService, keyProcessor)
	return cryptographyService, nil
}
