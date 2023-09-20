package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRandTraceID(t *testing.T) {
	traceID := RandTraceID()
	require.NotEmpty(t, traceID)
}
