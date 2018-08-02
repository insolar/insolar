package main

import (
	"io"
	"os"
	"plugin"
	"reflect"
	"testing"
)

func Test_generateForFile(t *testing.T) {
	w := generateForFile("../testplugins/secondary/main.go")
	io.Copy(os.Stdout, w)
}

// just to find methods
func TestConfigLoad(t *testing.T) {
	path, _ := os.Getwd()
	plugin, err := plugin.Open(path + "/../testplugins/secondary.so")
	if err != nil {
		t.Fatal(err)
	}

	hw, err := plugin.Lookup("INSEXPORT")

	r := reflect.ValueOf(hw)
	m := r.MethodByName("INSMETHOD__Hello")
	ret := m.Call([]reflect.Value{})
	t.Logf("%+v", ret)
}
