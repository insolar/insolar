package simplerequester

import (
	"encoding/json"
	"encoding/pem"
	"fmt"
	"github.com/insolar/go-jose"
	"github.com/insolar/x-crypto/ecdsa"
	"github.com/insolar/x-crypto/x509"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
	"path/filepath"
)

type SignedData struct {
	JWK string `json:"jwk"`
	JWS string `json:"jws"`
}

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

func createSignedData(privateKey *ecdsa.PrivateKey, datas *DataToSign) (string, string, error) {
	pub := jose.JSONWebKey{Key: privateKey.Public()}
	pubjs, err := pub.MarshalJSON()
	if err != nil {
		return "", "", err
	}
	signer, err := jose.NewSigner(jose.SigningKey{Algorithm: jose.ES256K, Key: privateKey}, nil)
	if err != nil {
		return "", "", err
	}

	payload, err := json.Marshal(datas)
	if err != nil {
		return "", "", err
	}
	object, err := signer.Sign(payload)
	if err != nil {
		return "", "", err
	}
	compactSerialized, _ := object.CompactSerialize()

	return compactSerialized, string(pubjs), nil

}

func getResponse(body []byte) (*response, error) {
	res := &response{}
	err := json.Unmarshal(body, &res)
	if err != nil {
		return nil, errors.Wrap(err, "[ getResponse ] problems with unmarshal response")
	}
	return res, nil
}

func importPrivateKeyPEM(pemEncoded []byte) (*ecdsa.PrivateKey, error) {
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

func ReadRequestParams(path string) (*DataToSign, error) {
	rConfig := &DataToSign{}
	err := readFile(path, rConfig)
	if err != nil {
		return nil, errors.Wrap(err, "[ readRequestParams ] ")
	}

	return rConfig, nil
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
