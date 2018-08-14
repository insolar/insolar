package main

import (
	"io/ioutil"
	"testing"
)

func Test_generateForFile(t *testing.T) {
	mode = "wrapper"
	w, err := generateForFile("../testplugins/secondary/main.go")
	if err != nil {
		t.Fatal(err)
	}
	// io.Copy(os.Stdout, w)
	b, err := ioutil.ReadAll(w)
	if err != nil {
		t.Fatal("reading from generated code", err)
	}
	if len(b) == 0 {
		t.Fatal("generator returns zero length code")
	}
}
