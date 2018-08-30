package buitin

import (
	"testing"

	"io/ioutil"
	"os"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/ledger"
	"github.com/stretchr/testify/assert"
)

func TNewBuiltIn(t *testing.T, dir string) *BuiltIn {
	l, err := ledger.NewLedger(configuration.Ledger{
		DataDirectory: dir,
	})
	assert.NoError(t, err)
	am := l.GetManager()
	assert.Equal(t, true, am != nil)

	bi := NewBuiltIn(am, nil)

	assert.Equal(t, true, bi != nil)

	return bi
}

func TestBareHelloworld(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "test-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir) // nolint: errcheck

	bi := TNewBuiltIn(t, tmpDir)
	am := bi.AM
	//hwRef:= helloworld.CodeRef()
	//hw := helloworld.NewHelloWorld()

	am.SetArchPref([]core.MachineType{0})
}
