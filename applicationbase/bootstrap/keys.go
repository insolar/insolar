// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package bootstrap

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar/secrets"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/platformpolicy"
)

func keysToNodeInfo(kp *secrets.KeyPair) nodeInfo {
	return nodeInfo{
		privateKey: kp.Private,
		publicKey:  platformpolicy.MustPublicKeyToString(kp.Public),
	}
}

func keyPairsToNodeInfo(kp ...*secrets.KeyPair) []nodeInfo {
	nodes := make([]nodeInfo, 0, len(kp))
	for _, p := range kp {
		nodes = append(nodes, keysToNodeInfo(p))
	}
	return nodes
}

func createKeysInDir(
	ctx context.Context,
	dir string,
	keyFilenameFormat string,
	nodes []Node,
	reuse bool,
) ([]nodeInfo, error) {
	amount := len(nodes)

	// XXX: Hack: works only for generated files with keyFilenameFormat
	// TODO: reconsider this option implementation - (INS-2473) - @nordicdyno 16.May.2019
	if reuse {
		pairs, err := secrets.ReadKeysFromDir(dir)
		if err != nil {
			return nil, err
		}
		if len(pairs) != amount {
			return nil, errors.New(fmt.Sprintf("[ uploadKeys ] amount of discoveryNodes != amount of files in directory: %d != %d", len(pairs), amount))
		}
		return keyPairsToNodeInfo(pairs...), nil
	}

	nodeInfos := make([]nodeInfo, 0, amount)
	for i, n := range nodes {
		keyname := fmt.Sprintf(keyFilenameFormat, i+certNamesStartFrom)
		if len(n.KeyName) > 0 {
			keyname = n.KeyName
		}

		pair, err := secrets.GenerateKeyPair()

		if err != nil {
			return nil, errors.Wrap(err, "[ createKeysInDir ] couldn't generate keys")
		}

		ks := platformpolicy.NewKeyProcessor()
		privKeyStr, err := ks.ExportPrivateKeyPEM(pair.Private)
		if err != nil {
			return nil, errors.Wrap(err, "[ createKeysInDir ] couldn't export private key")
		}

		pubKeyStr, err := ks.ExportPublicKeyPEM(pair.Public)
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

		inslogger.FromContext(ctx).Info("Genesis write key " + filepath.Join(dir, keyname))
		err = makeFileWithDir(dir, keyname, result)
		if err != nil {
			return nil, errors.Wrap(err, "[ createKeysInDir ] couldn't write keys to file")
		}

		p := keysToNodeInfo(pair)
		p.role = n.Role
		p.certName = n.CertName
		nodeInfos = append(nodeInfos, p)
	}

	return nodeInfos, nil
}

// makeFileWithDir saves content into file with `name` in directory `dir`.
// Creates directory if needed as well as file
func makeFileWithDir(dir string, name string, content []byte) error {
	err := os.MkdirAll(dir, 0700)
	if err != nil {
		return err
	}
	file := filepath.Join(dir, name)
	return ioutil.WriteFile(file, content, 0600)
}
