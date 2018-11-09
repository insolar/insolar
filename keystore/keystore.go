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
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/keystore/internal/privatekey"
	"github.com/pkg/errors"
)

type keyStore struct {
	Loader privatekey.Loader `inject:""`
	path   string
}

func (ks *keyStore) GetPrivateKey(identifier string) (crypto.PrivateKey, error) {
	file := privatekey.KeyFile(ks.path)
	return ks.Loader.Load(file)
}

func (ks *keyStore) Start(ctx context.Context) error {
	// TODO: ugly hack; do proper checks
	if _, err := ks.GetPrivateKey(""); err != nil {
		return errors.Wrap(err, "[ Start ] Failed to start keyStore")
	}

	return nil
}

