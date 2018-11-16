package functest

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInsgorundReload(t *testing.T) {
	_, err := signedRequest(&root, "DumpAllUsers")
	assert.NoError(t, err)

	stopInsgorund()
	err = startInsgorund()
	assert.NoError(t, err)

	_, err = signedRequest(&root, "DumpAllUsers")
	assert.NoError(t, err)
}
