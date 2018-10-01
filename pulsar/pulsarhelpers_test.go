package pulsar

import (
	"crypto/ecdsa"
	"crypto/rand"
	"testing"

	ecdsa_helper "github.com/insolar/insolar/crypto_helpers/ecdsa"
	"github.com/stretchr/testify/assert"
)

func TestSingAndVerify(t *testing.T) {
	assertObj := assert.New(t)
	privateKey, _ := ecdsa.GenerateKey(ecdsa_helper.GetCurve(), rand.Reader)
	publicKey, _ := ecdsa_helper.ExportPublicKey(&privateKey.PublicKey)

	signature, err := singData(privateKey, "This is the message to be signed and verified!")
	assertObj.NoError(err)

	checkSignature, err := checkPayloadSignature(&Payload{PublicKey: publicKey, Signature: signature, Body: "This is the message to be signed and verified!"})

	assertObj.Equal(true, checkSignature)
}
