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

package main

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"strings"
	"testing"

	"github.com/insolar/insolar/conveyor/generator/generator"
	"github.com/stretchr/testify/require"
)

func Test_Main(t *testing.T) {
	g := generator.NewGenerator("conveyor/generator/state_machines/",
		"conveyor/generator/matrix/matrix.go")
	files, err := ioutil.ReadDir("state_machines/")
	require.NoError(t, err)
	for _, file := range files {
		if file.IsDir() {
			dirName := file.Name()
			files, err := ioutil.ReadDir("state_machines/" + dirName)
			require.NoError(t, err)
			for _, file := range files {
				if !strings.HasSuffix(file.Name(), "generated.go") {
					g.ParseFile(dirName, file.Name())
				}
			}
			continue
		}
		if !strings.HasSuffix(file.Name(), "generated.go") {
			g.ParseFile("", file.Name())
		}
	}
	g.GenMatrix()

	out, err := exec.Command("go", "test", "-tags=with_generated", "./state_machine_test.go").CombinedOutput()
	fmt.Println(err)
	require.NoError(t, err)
	fmt.Println(string(out))
}
