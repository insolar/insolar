package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnvFile(t *testing.T) {
	err := addValueToEnvFile("INS_LATEST_VER", "v1.1.0")
	assert.NoError(t, err)
}
