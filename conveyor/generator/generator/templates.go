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
	"io"
	"io/ioutil"
	"path"
	"text/template"
)

type tplParams struct {
	Machine stateMachine
	Handler *handler
}

var (
	funcMap = template.FuncMap{
		"inc": func(i int) int {
			return i + 1
		},
		"handlerExists": func(x *handler) bool {
			return x != nil
		},
		"gtNull": func(i int) bool {
			return i > 0
		},
		"params": func(m stateMachine, h *handler) tplParams {
			return tplParams{
				Machine: m,
				Handler: h,
			}
		},
	}
)

func (p *Parser) Generate(w io.Writer) {
	file := path.Join(p.generator.fullPathToInsolar, stateMachineTemplate)
	tplBody, err := ioutil.ReadFile(file)
	checkErr(err)
	err = template.Must(template.New("stateMachineTpl").Funcs(funcMap).
		Parse(string(tplBody))).
		Execute(w, p)
	checkErr(err)
}
