package main

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"path/filepath"

	"github.com/insolar/insolar/api/requester"
	"github.com/insolar/insolar/platformpolicy"
	xcrypto "github.com/insolar/x-crypto"
	xecdsa "github.com/insolar/x-crypto/ecdsa"
	xrand "github.com/insolar/x-crypto/rand"
	xx509 "github.com/insolar/x-crypto/x509"
	"github.com/pkg/errors"
)

type memberKeys struct {
	Private string `json:"private_key"`
	Public  string `json:"public_key"`
}

func readRequestParams(path string) (*requester.Request, error) {

	fileParams := &requester.Request{}
	err := readFile(path, fileParams)
	if err != nil {
		return nil, errors.Wrap(err, "[ readRequestParams ] ")
	}

	return fileParams, nil
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

func check(msg string, err error) {
	if err != nil {
		fmt.Println(msg, err)
		os.Exit(1)
	}
}

func execute(url string, keys memberKeys, request requester.Request) (string, error) {
	seed, err := requester.GetSeed(url)
	check("[Execute]", err)
	request.Params.PublicKey = keys.Public
	request.Params.Seed = seed
	dataToSign, err := json.Marshal(request)
	check("[Execute]", err)
	signature, err := sign(keys.Private, dataToSign)
	check("[Execute]", err)

	fmt.Println("Request: " + string(dataToSign))

	body, err := requester.GetResponseBodyContract(url+"/call", request, signature)
	check("[Execute]", err)
	return string(body), nil
}

func sign(privateKeyPem string, data []byte) (string, error) {
	hash := sha256.Sum256(data)

	privateKey, curveName, err := importPrivateKeyPEM([]byte(privateKeyPem))
	if err != nil {
		panic(err)
	}
	fmt.Println(curveName)
	switch curveName {
	case "P-256":
		ks := platformpolicy.NewKeyProcessor()
		privateKey, err := ks.ImportPrivateKeyPEM([]byte(privateKeyPem))
		if err != nil {
			panic(err)
		}
		r, s, err := ecdsa.Sign(rand.Reader, privateKey.(*ecdsa.PrivateKey), hash[:])
		if err != nil {
			panic(err)
		}

		return pointsToDER(r, s), nil

	case "P-256K":
		r, s, err := xecdsa.Sign(xrand.Reader, privateKey.(*xecdsa.PrivateKey), hash[:])
		if err != nil {
			panic(err)
		}

		return pointsToDER(r, s), nil
	}
	return "", errors.New("Undefined keys curve name")
}

func pointsToDER(r, s *big.Int) string {
	prefixPoint := func(b []byte) []byte {
		if len(b) == 0 {
			b = []byte{0x00}
		}
		if b[0]&0x80 != 0 {
			paddedBytes := make([]byte, len(b)+1)
			copy(paddedBytes[1:], b)
			b = paddedBytes
		}
		return b
	}

	rb := prefixPoint(r.Bytes())
	sb := prefixPoint(s.Bytes())

	// DER encoding:
	// 0x30 + z + 0x02 + len(rb) + rb + 0x02 + len(sb) + sb
	length := 2 + len(rb) + 2 + len(sb)

	der := append([]byte{0x30, byte(length), 0x02, byte(len(rb))}, rb...)
	der = append(der, 0x02, byte(len(sb)))
	der = append(der, sb...)

	return base64.StdEncoding.EncodeToString(der)
}

func importPrivateKeyPEM(pemEncoded []byte) (xcrypto.PrivateKey, string, error) {
	block, _ := pem.Decode(pemEncoded)
	if block == nil {
		return nil, "", fmt.Errorf("[ ImportPrivateKey ] Problems with decoding. Key - %v", pemEncoded)
	}
	x509Encoded := block.Bytes
	privateKey, err := xx509.ParseECPrivateKey(x509Encoded)

	if err != nil {
		return nil, "", fmt.Errorf("[ ImportPrivateKey ] Problems with parsing. Key - %v", pemEncoded)
	}
	return privateKey, privateKey.Curve.Params().Name, nil
}
