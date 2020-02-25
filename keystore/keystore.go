// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package keystore

import (
	"context"
	"crypto"

	component "github.com/insolar/component-manager"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/keystore/internal/privatekey"
	"github.com/pkg/errors"
)

type keyStore struct {
	Loader privatekey.Loader `inject:""`
	file   string
}

func (ks *keyStore) GetPrivateKey(identifier string) (crypto.PrivateKey, error) {
	return ks.Loader.Load(ks.file)
}

func (ks *keyStore) Start(ctx context.Context) error {
	// TODO: ugly hack; do proper checks
	if _, err := ks.GetPrivateKey(""); err != nil {
		return errors.Wrap(err, "[ Start ] Failed to start keyStore")
	}

	return nil
}

type cachedKeyStore struct {
	keyStore insolar.KeyStore

	privateKey crypto.PrivateKey
}

func (ks *cachedKeyStore) getCachedPrivateKey() crypto.PublicKey {
	if ks.privateKey != nil {
		return ks.privateKey
	}
	return nil
}

func (ks *cachedKeyStore) loadPrivateKey(identifier string) (crypto.PrivateKey, error) {
	privateKey, err := ks.keyStore.GetPrivateKey(identifier)
	if err != nil {
		return nil, errors.Wrap(err, "[ loadPrivateKey ] Can't GetPrivateKey")
	}

	ks.privateKey = privateKey
	return privateKey, nil
}

func (ks *cachedKeyStore) GetPrivateKey(_ string) (crypto.PrivateKey, error) {
	privateKey := ks.getCachedPrivateKey()

	return privateKey, nil
}

func (ks *cachedKeyStore) Start(ctx context.Context) error {
	// TODO: ugly hack; do proper checks
	if _, err := ks.loadPrivateKey(""); err != nil {
		return errors.Wrap(err, "[ Start ] Failed to start keyStore")
	}

	return nil
}

func NewKeyStore(path string) (insolar.KeyStore, error) {
	keyStore := &keyStore{
		file: path,
	}

	cachedKeyStore := &cachedKeyStore{
		keyStore: keyStore,
	}

	manager := component.NewManager(nil)
	manager.Inject(
		cachedKeyStore,
		keyStore,
		privatekey.NewLoader(),
	)

	if err := manager.Start(context.Background()); err != nil {
		return nil, errors.Wrap(err, "[ NewKeyStore ] Failed to create keyStore")
	}

	return cachedKeyStore, nil
}

type inPlaceKeyStore struct {
	privateKey crypto.PrivateKey
}

func (ipks *inPlaceKeyStore) GetPrivateKey(string) (crypto.PrivateKey, error) {
	return ipks.privateKey, nil
}

func NewInplaceKeyStore(privateKey crypto.PrivateKey) insolar.KeyStore {
	return &inPlaceKeyStore{privateKey: privateKey}
}
