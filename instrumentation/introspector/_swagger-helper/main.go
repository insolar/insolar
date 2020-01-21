// Copyright 2020 Insolar Network Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// swagger helper generator:
// encodes <name>.swagger.json files to literals with name <name>Swagger in swagger_const_gen.go file.

package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/insolar/insolar/log"
)

var (
	pkgName     = "introspector"
	outFileName = "swagger_const_gen.go"
	suffix      = ".swagger.json"
)

var inDir = flag.String("in", ".", "directory with swagger files")

func main() {
	flag.Parse()

	tmpF, err := ioutil.TempFile("", "swghelper_*.go")
	if err != nil {
		log.Fatal("failed open tmp file:", err)
	}

	files, err := ioutil.ReadDir(*inDir)
	if err != nil {
		log.Fatal("filed to read current directory", err)
	}

	sw := &strictWriter{w: tmpF}
	sw.writeString(fmt.Sprintf("package %v\n", pkgName))
	sw.writeString(preambula)
	sw.writeString("const (\n")
	for _, info := range files {
		if strings.HasSuffix(info.Name(), suffix) {
			name := strings.TrimSuffix(info.Name(), suffix)
			sw.writeString(name + "Swagger = `")
			filePath := path.Join(*inDir, info.Name())
			f, err := os.Open(filePath)
			if err != nil {
				log.Fatalf("failed to read file %v: %s", filePath, err)
			}
			sw.write(f)
			sw.writeString("`\n")
		}
	}
	sw.writeString(")\n")

	err = os.Rename(tmpF.Name(), outFileName)
	if err != nil {
		log.Fatalf("failed move file from %v to %v: %s", tmpF.Name(), outFileName, err)
	}

	cwd, _ := os.Getwd()
	_, _ = fmt.Fprintf(os.Stderr, "Generated file: %v (%v)", outFileName, cwd)
}

type strictWriter struct {
	w io.Writer
}

func (sw *strictWriter) writeString(s string) {
	_, err := sw.w.Write([]byte(s))
	if err != nil {
		panic(err)
	}
}

func (sw *strictWriter) write(r io.Reader) {
	_, err := io.Copy(sw.w, r)
	if err != nil {
		panic(err)
	}
}

var preambula = `
/*
DO NOT EDIT!
This code was generated automatically using _swagger-helper
*/

`
