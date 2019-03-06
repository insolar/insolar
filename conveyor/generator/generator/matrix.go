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

import (
	"bufio"
	"io/ioutil"
	"os"
	"path"
	"text/template"
)

var (
	matrixFuncMap = template.FuncMap{
		"isNull": func(i int) bool {
			return i == 0
		},
	}
)

func (g *Generator) getImports() []string {
	keys := make([]string, len(g.imports))
	i := 0
	for k := range g.imports {
		keys[i] = k
		i++
	}
	return keys
}

func (g *Generator) GenMatrix () {
	fileName := path.Join(g.fullPathToInsolar, matrixTemplate)
	tplBody, err := ioutil.ReadFile(fileName)
	checkErr(err)

	file, err := os.Create(path.Join(g.fullPathToInsolar, g.pathToMatrixFile))
	checkErr(err)
	defer file.Close()
	out := bufio.NewWriter(file)

	err = template.Must(template.New("matrixTmpl").Funcs(matrixFuncMap).
		Parse(string(tplBody))).Execute(out, struct{
		Imports []string
		Machines []*stateMachine
	}{
		Imports: g.getImports(),
		Machines: g.stateMachines,

	})
	checkErr(err)
	out.Flush()
}


