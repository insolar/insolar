// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

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

	err = MoveFile(tmpF.Name(), outFileName)
	if err != nil {
		log.Fatalf("failed move file from %v to %v: %s", tmpF.Name(), outFileName, err)
	}

	cwd, _ := os.Getwd()
	_, _ = fmt.Fprintf(os.Stderr, "Generated file: %v (%v)", outFileName, cwd)
}

/*
   GoLang: os.Rename() give error "invalid cross-device link" for Docker container with Volumes.
   MoveFile(source, destination) will work moving file between folders
*/
func MoveFile(sourcePath, destPath string) error {
	inputFile, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("Couldn't open source file: %s", err)
	}
	outputFile, err := os.Create(destPath)
	if err != nil {
		inputFile.Close()
		return fmt.Errorf("Couldn't open dest file: %s", err)
	}
	defer outputFile.Close()
	_, err = io.Copy(outputFile, inputFile)
	inputFile.Close()
	if err != nil {
		return fmt.Errorf("Writing to output file failed: %s", err)
	}
	// The copy was successful, so now delete the original file
	err = os.Remove(sourcePath)
	if err != nil {
		return fmt.Errorf("Failed removing original file: %s", err)
	}
	return nil
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
