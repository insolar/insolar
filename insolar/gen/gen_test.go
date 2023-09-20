package gen

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGen_Signature(t *testing.T) {
	zero := Signature(0)
	assert.NotNil(t, zero)
	assert.Equal(t, 0, len(zero))

	one := Signature(1)
	assert.NotNil(t, one)
	assert.Equal(t, 1, len(one))

	case256 := Signature(256)
	assert.NotNil(t, case256)
	assert.Equal(t, 256, len(case256))

	negative := Signature(-1)
	assert.Nil(t, negative)
}
