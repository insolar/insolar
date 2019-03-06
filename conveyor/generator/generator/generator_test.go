/*
 *    Copyright 2019 Insolar Technologies
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package generator

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

/*func readTestCode(t *testing.T) (*ast.File, string) {
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
*/
