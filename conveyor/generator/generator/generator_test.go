package generator

import (
	"testing"
	"io/ioutil"
	"github.com/stretchr/testify/assert"
	"os"
	"go/token"
	"go/parser"
	"go/ast"
)

const testCode = `
package sample
// conveyer: state_machine
type TestStateMachine interface {
	Init(input Event) (*Payload, common.ElState, error)
	StateFirst() common.ElState
	TransitFirstSecond(input Event, payload *Payload) (*Payload, common.ElState, error)
	MigrateFirst(input Event, payload *Payload) (*Payload, common.ElState, error)
	ErrorFirst(input Event, payload *Payload, err error) (*Payload, common.ElState)
	StateSecond() common.ElState
	TransitSecondThird(input Event, payload *Payload) (*Payload, common.ElState, error)
	MigrateSecond(input Event, payload *Payload) (*Payload, common.ElState, error)
	ErrorSecond(input Event, payload *Payload, err error) (*Payload, common.ElState)
}
`

func readTestCode(t *testing.T) (*ast.File, string) {
	tmpFile, err := ioutil.TempFile("", "test_")
	assert.NoError(t, err)
	_, err = tmpFile.Write([]byte(testCode))
	assert.NoError(t, err)
	err = tmpFile.Close()
	assert.NoError(t, err)
	defer os.Remove(tmpFile.Name()) // nolint: errcheck

	fSet := token.NewFileSet()
	node, err := parser.ParseFile(fSet, tmpFile.Name(), nil, parser.ParseComments)
	assert.NoError(t, err)
	return node, tmpFile.Name()
}

func testGenerator(t *testing.T) *Generator {
	node, fileName := readTestCode(t)
	return &Generator{
		sourceFilename: fileName,
		sourceCode: []byte(testCode),
		sourceNode: node,
	}
}

func TestGenerator_findEachStateMachine(t *testing.T) {
	g := testGenerator(t)
	g.findEachStateMachine()
	assert.Equal(t, 1, len(g.stateMachines))
	assert.Equal(t, 2, len(g.stateMachines[0].States))
}
