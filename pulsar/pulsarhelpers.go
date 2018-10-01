package pulsar

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/gob"
	"errors"
	"math/big"
	"sort"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/crypto_helpers"
	ecdsa_helper "github.com/insolar/insolar/crypto_helpers/ecdsa"
	"github.com/insolar/insolar/ledger/hash"
)

func checkPayloadSignature(request *Payload) (bool, error) {
	return checkSignature(request.Body, request.PublicKey, request.Signature)
}

func checkSignature(data interface{}, pub string, signature []byte) (bool, error) {
	b := bytes.Buffer{}
	e := gob.NewEncoder(&b)
	err := e.Encode(data)
	if err != nil {
		return false, err
	}

	r := big.Int{}
	s := big.Int{}
	sigLen := len(signature)
	r.SetBytes(signature[:(sigLen / 2)])
	s.SetBytes(signature[(sigLen / 2):])

	hash := crypto_helpers.MakeSha3Hash(b.Bytes())
	publicKey, err := ecdsa_helper.ImportPublicKey(pub)
	if err != nil {
		return false, err
	}

	return ecdsa.Verify(publicKey, hash[:], &r, &s), nil
}

func singData(privateKey *ecdsa.PrivateKey, data interface{}) ([]byte, error) {
	b := bytes.Buffer{}
	e := gob.NewEncoder(&b)
	err := e.Encode(data)
	if err != nil {
		return nil, err
	}

	hash := crypto_helpers.MakeSha3Hash(b.Bytes())
	r, s, err := ecdsa.Sign(rand.Reader, privateKey, hash[:])
	if err != nil {
		return nil, err
	}

	return append(r.Bytes(), s.Bytes()...), nil
}

func selectByEntropy(entropy core.Entropy, values []string, count int) ([]string, error) { // nolint: megacheck
	type idxHash struct {
		idx  int
		hash []byte
	}

	if len(values) < count {
		return nil, errors.New("count value should be less than values size")
	}

	hashes := make([]*idxHash, 0, len(values))
	for i, value := range values {
		h := hash.NewSHA3()
		_, err := h.Write(entropy[:])
		if err != nil {
			return nil, err
		}
		_, err = h.Write([]byte(value))
		if err != nil {
			return nil, err
		}
		hashes = append(hashes, &idxHash{
			idx:  i,
			hash: h.Sum(nil),
		})
	}

	sort.SliceStable(hashes, func(i, j int) bool { return bytes.Compare(hashes[i].hash, hashes[j].hash) < 0 })

	selected := make([]string, 0, count)
	for i := 0; i < count; i++ {
		selected = append(selected, values[hashes[i].idx])
	}
	return selected, nil
}
