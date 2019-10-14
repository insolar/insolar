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

package integration

import (
	"encoding/json"
	"io/ioutil"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/platformpolicy"
)

func loadMemberKeys(keysPath string) (*User, error) {
	text, err := ioutil.ReadFile(keysPath)
	if err != nil {
		return nil, errors.Wrapf(err, "[ loadMemberKeys ] could't load member keys")
	}

	var data map[string]string
	err = json.Unmarshal(text, &data)
	if err != nil {
		return nil, errors.Wrapf(err, "[ loadMemberKeys ] could't unmarshal member keys")
	}
	if data["private_key"] == "" || data["public_key"] == "" {
		return nil, errors.New("[ loadMemberKeys ] could't find any keys")
	}

	return &User{
		PrivateKey: data["private_key"],
		PublicKey:  data["public_key"],
	}, nil
}

type User struct {
	Reference        insolar.Reference
	PrivateKey       string
	PublicKey        string
	MigrationAddress string
}

func NewUserWithKeys() (*User, error) {
	ks := platformpolicy.NewKeyProcessor()

	privateKey, err := ks.GeneratePrivateKey()
	if err != nil {
		return nil, err
	}
	privateKeyString, err := ks.ExportPrivateKeyPEM(privateKey)
	if err != nil {
		return nil, err
	}

	publicKey := ks.ExtractPublicKey(privateKey)
	publicKeyString, err := ks.ExportPublicKeyPEM(publicKey)
	if err != nil {
		return nil, err
	}

	return &User{
		PrivateKey: string(privateKeyString),
		PublicKey:  string(publicKeyString),
	}, nil
}
