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

package keystore

import (
	"context"
	"crypto"

	"github.com/insolar/insolar/component"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/keystore/internal/privatekey"
	"github.com/pkg/errors"
)

type keyStore struct {
	Loader privatekey.Loader `inject:""`
	path   string
}

func (ks *keyStore) GetPrivateKey(identifier string) (crypto.PrivateKey, error) {
	return ks.Loader.Load(ks.path)
}

func (ks *keyStore) Start(ctx context.Context) error {
	// TODO: ugly hack; do proper checks
	if _, err := ks.GetPrivateKey(""); err != nil {
		return errors.Wrap(err, "[ Start ] Failed to start keyStore")
	}

	return nil
}

type cachedKeyStore struct {
	keyStore core.KeyStore

	privateKey crypto.PrivateKey
}

func (ks *cachedKeyStore) getCachedPrivateKey(identifier string) crypto.PublicKey {
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

func (ks *cachedKeyStore) GetPrivateKey(identifier string) (crypto.PrivateKey, error) {
	privateKey := ks.getCachedPrivateKey(identifier)

	return privateKey, nil
}

func (ks *cachedKeyStore) Start(ctx context.Context) error {
	// TODO: ugly hack; do proper checks
	if _, err := ks.loadPrivateKey(""); err != nil {
		return errors.Wrap(err, "[ Start ] Failed to start keyStore")
	}

	return nil
}

func NewKeyStore(path string) (core.KeyStore, error) {
	keyStore := &keyStore{
		path: path,
	}

	cachedKeyStore := &cachedKeyStore{
		keyStore: keyStore,
	}

	manager := component.Manager{}
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
