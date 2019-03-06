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
	"io/ioutil"
	"strings"

	"github.com/insolar/insolar/conveyor/generator/generator"
	"github.com/insolar/insolar/log"
)

func main() {
	g := generator.NewGenerator( "conveyor/generator/state_machines/",
		"conveyor/generator/matrix/matrix.go")
	files, err := ioutil.ReadDir("conveyor/generator/state_machines/")
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		if file.IsDir() {
			dirName := file.Name()
			files, err := ioutil.ReadDir("conveyor/generator/state_machines/" + dirName)
			if err != nil {
				log.Fatal(err)
			}
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
}
