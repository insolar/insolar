package main

import (
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	insJose "github.com/insolar/go-jose"
	"github.com/insolar/insolar/platformpolicy"
	xecdsa "github.com/insolar/x-crypto/ecdsa"
	"github.com/insolar/x-crypto/x509"
	"github.com/pkg/errors"
	"github.com/square/go-jose"
)

type DataToSign struct {
	Reference string `json:"reference"`
	Method    string `json:"method"`
	Params    string `json:"params"`
	Seed      string `json:"seed"`
}

type memberKeys struct {
	Private string `json:"private_key"`
	Public  string `json:"public_key"`
}

func createSignedData(keys *memberKeys, datas *DataToSign) (string, string, error) {

	payload, err := json.Marshal(datas)
	if err != nil {
		return "", "", err
	}
	privateKey, err := importPrivateKeyPEM([]byte(keys.Private))

	switch privateKey.Curve.Params().Name {

	case "P-256K":
		pub := insJose.JSONWebKey{Key: privateKey.Public()}
		pubjs, err := pub.MarshalJSON()
		if err != nil {
			return "", "", err
		}
		signer, err := insJose.NewSigner(insJose.SigningKey{Algorithm: insJose.ES256K, Key: privateKey}, nil)
		if err != nil {
			return "", "", err
		}
		object, err := signer.Sign(payload)
		if err != nil {
			return "", "", err
		}
		compactSerialized, err := object.CompactSerialize()
		if err != nil {
			return "", "", err
		}
		return compactSerialized, string(pubjs), nil

	default:
		keyProcessor := platformpolicy.NewKeyProcessor()
		privateKey, err := keyProcessor.ImportPrivateKeyPEM([]byte(keys.Private))
		publicKey := keyProcessor.ExtractPublicKey(privateKey)

		pub := jose.JSONWebKey{Key: publicKey}
		pubjs, err := pub.MarshalJSON()
		if err != nil {
			return "", "", err
		}
		signer, err := jose.NewSigner(jose.SigningKey{Algorithm: jose.ES256, Key: privateKey}, nil)
		if err != nil {
			return "", "", err
		}
		object, err := signer.Sign(payload)
		if err != nil {
			return "", "", err
		}
		compactSerialized, err := object.CompactSerialize()
		if err != nil {
			return "", "", err
		}
		return compactSerialized, string(pubjs), nil
	}
}

func getResponse(body []byte) (*response, error) {
	res := &response{}
	err := json.Unmarshal(body, &res)
	if err != nil {
		return nil, errors.Wrap(err, "[ getResponse ] problems with unmarshal response")
	}
	return res, nil
}

func importPrivateKeyPEM(pemEncoded []byte) (*xecdsa.PrivateKey, error) {
	block, _ := pem.Decode(pemEncoded)
	if block == nil {
		return nil, fmt.Errorf("[ ImportPrivateKey ] Problems with decoding. Key - %v", pemEncoded)
	}
	x509Encoded := block.Bytes
	privateKey, err := x509.ParseECPrivateKey(x509Encoded)
	if err != nil {
		return nil, fmt.Errorf("[ ImportPrivateKey ] Problems with parsing. Key - %v", pemEncoded)
	}
	return privateKey, nil
}

func readRequestParams(path string) (*DataToSign, error) {
	type DataToSignFile struct {
		Reference string      `json:"reference"`
		Method    string      `json:"method"`
		Params    interface{} `json:"params"`
	}

	fileParams := &DataToSignFile{}
	err := readFile(path, fileParams)
	if err != nil {
		return nil, errors.Wrap(err, "[ readRequestParams ] ")
	}
	parameters, err := json.Marshal(fileParams.Params)
	if err != nil {
		return nil, errors.Wrap(err, "[ readRequestParams ] ")
	}
	params := &DataToSign{
		Method:    fileParams.Method,
		Reference: fileParams.Reference,
		Params:    string(parameters),
	}

	return params, nil
}

func readFile(path string, configType interface{}) error {
	var rawConf []byte
	var err error
	if path == "-" {
		rawConf, err = ioutil.ReadAll(os.Stdin)
	} else {
		rawConf, err = ioutil.ReadFile(filepath.Clean(path))
	}
	if err != nil {
		return errors.Wrap(err, "[ readFile ] Problem with reading config")
	}

	err = json.Unmarshal(rawConf, &configType)
	if err != nil {
		return errors.Wrap(err, "[ readFile ] Problem with unmarshaling config")
	}

	return nil
}
