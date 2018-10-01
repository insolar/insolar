package pulsar

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSingAndVerify(t *testing.T) {
	assertObj := assert.New(t)
	privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	publicKey, _ := ExportPublicKey(&privateKey.PublicKey)

	signature, err := singData(privateKey, "This is the message to be signed and verified!")
	assertObj.NoError(err)

	checkSignature, err := checkPayloadSignature(&Payload{PublicKey: publicKey, Signature: signature, Body: "This is the message to be signed and verified!"})

	assertObj.Equal(true, checkSignature)
}
