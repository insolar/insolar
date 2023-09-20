package zlogadapter

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func stripPackageName(packageName string) string {
	result := strings.TrimPrefix(packageName, insolarPrefix)
	i := strings.Index(result, ".")
	if result == packageName || i == -1 {
		return result
	}
	return result[:i]
}

// beware to adding lines in this test (test output depend on test code offset!)
func TestLog_getCallInfo(t *testing.T) {
	expectedLine := 27 // should be equal of line number where getCallInfo is called
	info := getCallInfo(1)

	assert.Contains(t, info.fileName, "log/zlogadapter/sourceinfo_test.go:")
	assert.Equal(t, "TestLog_getCallInfo", info.funcName)
	assert.Equal(t, expectedLine, info.line)
}

func TestLog_stripPackageName(t *testing.T) {
	tests := map[string]struct {
		packageName string
		result      string
	}{
		"insolar":    {"github.com/insolar/insolar/mypackage", "mypackage"},
		"thirdParty": {"github.com/stretchr/testify/assert", "github.com/stretchr/testify/assert"},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, test.result, stripPackageName(test.packageName))
		})
	}
}
