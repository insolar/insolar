package pulsar

import (
	"testing"

	ecdsa_helper "github.com/insolar/insolar/crypto_helpers/ecdsa"
	"github.com/stretchr/testify/assert"
)

func TestSingAndVerify(t *testing.T) {
	assertObj := assert.New(t)
	privateKey, err := ecdsa_helper.GeneratePrivateKey()
	assert.NoError(t, err)
	publicKey, err := ecdsa_helper.ExportPublicKey(&privateKey.PublicKey)
	assert.NoError(t, err)

	signature, err := singData(privateKey, "This is the message to be signed and verified!")
	assertObj.NoError(err)

	checkSignature, err := checkPayloadSignature(&Payload{PublicKey: publicKey, Signature: signature, Body: "This is the message to be signed and verified!"})

	assertObj.Equal(true, checkSignature)
}
