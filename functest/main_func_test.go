package functest

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	os.Exit(TestsMainWrapper(m))
}
