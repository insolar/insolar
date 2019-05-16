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

package genesis

import (
	"bytes"
	"context"
	"crypto"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"path/filepath"

	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/platformpolicy"
	"github.com/pkg/errors"
)

func getKeysFromFile(file string) (crypto.PrivateKey, string, error) {
	absPath, err := filepath.Abs(file)
	if err != nil {
		return nil, "", errors.Wrap(err, "[ getKeyFromFile ] couldn't get abs path")
	}
	b, err := ioutil.ReadFile(absPath)
	if err != nil {
		return nil, "", errors.Wrap(err, "[ getKeyFromFile ] couldn't read keys file "+absPath)
	}
	return getKeys(bytes.NewReader(b))
}

func getKeys(r io.Reader) (crypto.PrivateKey, string, error) {
	var keys map[string]string
	err := json.NewDecoder(r).Decode(&keys)
	if err != nil {
		return nil, "", errors.Wrapf(err, "[ getKeys ] couldn't unmarshal keys data")
	}
	if keys["private_key"] == "" {
		return nil, "", errors.New("[ getKeys ] empty private key")
	}
	if keys["public_key"] == "" {
		return nil, "", errors.New("[ getKeys ] empty public key")
	}

	kp := platformpolicy.NewKeyProcessor()
	key, err := kp.ImportPrivateKeyPEM([]byte(keys["private_key"]))
	if err != nil {
		return nil, "", errors.Wrapf(err, "[ getKeys ] couldn't import private key")
	}

	return key, mustNormalizePublicKey(keys["public_key"]), nil
}

func mustNormalizePublicKey(s string) string {
	ks := platformpolicy.NewKeyProcessor()
	pubKey, err := ks.ImportPublicKeyPEM([]byte(s))
	if err != nil {
		panic(err)
	}
	return string(mustPublicKeyToBytes(pubKey))
}

func mustPublicKeyToBytes(key crypto.PublicKey) []byte {
	ks := platformpolicy.NewKeyProcessor()
	b, err := ks.ExportPublicKeyPEM(key)
	if err != nil {
		panic(err)
	}
	return b
}

func readKeysFromDir(dir string, amount int) ([]nodeInfo, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, errors.Wrap(err, "[ uploadKeys ] can't read dir")
	}

	nodes := make([]nodeInfo, amount)
	if len(files) != amount {
		return nil, errors.New(fmt.Sprintf("[ uploadKeys ] amount of nodes != amount of files in directory: %d != %d", len(files), amount))
	}

	for i, f := range files {
		privKey, nodePubKey, err := getKeysFromFile(filepath.Join(dir, f.Name()))
		if err != nil {
			return nil, errors.Wrap(err, "[ uploadKeys ] can't get keys from file")
		}

		nodes[i].publicKey = nodePubKey
		nodes[i].privateKey = privKey
	}

	return nodes, nil
}

func createKeysInDir(
	ctx context.Context,
	dir string,
	keyFilenameFormat string,
	discoveryNodes []Node,
	reuse bool,
) ([]nodeInfo, error) {
	amount := len(discoveryNodes)

	// XXX: Hack: works only for generated files by keyFilenameFormat
	// TODO: reconsider this option implementation - (INS-2473) - @nordicdyno 16.May.2019
	if reuse {
		return readKeysFromDir(dir, amount)
	}

	nodes := make([]nodeInfo, amount)
	for i, dn := range discoveryNodes {
		name := fmt.Sprintf(keyFilenameFormat, i+1)
		if len(dn.KeyName) > 0 {
			name = dn.KeyName
		}

		ks := platformpolicy.NewKeyProcessor()

		privKey, err := ks.GeneratePrivateKey()
		if err != nil {
			return nil, errors.Wrap(err, "[ createKeysInDir ] couldn't generate private key")
		}

		privKeyStr, err := ks.ExportPrivateKeyPEM(privKey)
		if err != nil {
			return nil, errors.Wrap(err, "[ createKeysInDir ] couldn't export private key")
		}

		pubKey := ks.ExtractPublicKey(privKey)
		pubKeyStr, err := ks.ExportPublicKeyPEM(pubKey)
		if err != nil {
			return nil, errors.Wrap(err, "[ createKeysInDir ] couldn't export public key")
		}

		result, err := json.MarshalIndent(map[string]interface{}{
			"private_key": string(privKeyStr),
			"public_key":  string(pubKeyStr),
		}, "", "    ")
		if err != nil {
			return nil, errors.Wrap(err, "[ createKeysInDir ] couldn't marshal keys")
		}

		inslogger.FromContext(ctx).Info("Genesis write key " + filepath.Join(dir, name))
		err = makeFileWithDir(dir, name, result)
		if err != nil {
			return nil, errors.Wrap(err, "[ createKeysInDir ] couldn't write keys to file")
		}

		nodes[i].publicKey = string(mustPublicKeyToBytes(pubKey))
		nodes[i].privateKey = privKey
	}

	return nodes, nil
}
