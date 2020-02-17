// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package privatekey

import (
	"crypto"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"io/ioutil"
	"path/filepath"

	"github.com/pkg/errors"
)

type keyLoader struct {
	parseFunc func(key []byte) (crypto.PrivateKey, error)
}

func NewLoader() Loader {
	return &keyLoader{
		parseFunc: pemParse,
	}
}

func (p *keyLoader) Load(file string) (crypto.PrivateKey, error) {
	key, err := readJSON(file)
	if err != nil {
		return nil, errors.Wrap(err, "[ Load ] Could't read private key")
	}

	signer, err := p.parseFunc(key)
	if err != nil {
		return nil, errors.Wrap(err, "[ Load ] Could't parse private key")
	}
	return signer, nil
}

// deprecated, todo: use PEM format
func readJSON(path string) ([]byte, error) {
	data, err := ioutil.ReadFile(filepath.Clean(path))
	if err != nil {
		return nil, errors.Wrap(err, "[ read ] couldn't read keys from: "+path)
	}
	var keys map[string]string
	err = json.Unmarshal(data, &keys)
	if err != nil {
		return nil, errors.Wrap(err, "[ read ] failed to parse json.")
	}

	key, ok := keys["private_key"]
	if !ok {
		return nil, errors.Errorf("[ read ] couldn't read keys from: %s", path)
	}

	return []byte(key), nil
}

func pemParse(key []byte) (crypto.PrivateKey, error) {
	block, _ := pem.Decode(key)
	if block == nil {
		return nil, errors.Errorf("[ Parse ] Problems with decoding PEM")
	}

	x509Encoded := block.Bytes
	privateKey, err := x509.ParsePKCS8PrivateKey(x509Encoded)
	if err != nil {
		// try to read old version marshalled with x509.MarshalECPrivateKey()
		privateKey, err = x509.ParseECPrivateKey(x509Encoded)
		if err != nil {
			return nil, errors.Errorf("[ Parse ] Problems with parsing private key")
		}
	}

	return privateKey, nil
}
