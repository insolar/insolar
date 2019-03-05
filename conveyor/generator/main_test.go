package main

import (
	"fmt"
	"os/exec"
	"testing"
	"github.com/insolar/insolar/conveyor/generator/generator"
	"io/ioutil"
	"strings"
	"github.com/stretchr/testify/require"
)

func Test_Main(t *testing.T) {
	g := generator.NewGenerator(
		"github.com/insolar/insolar/conveyor/generator/",
		"state_machines",
		"matrix/matrix.go")
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
	err = g.GenMatrix()
	require.NoError(t, err)

	out, err := exec.Command("go", "test", "-tags=with_generated", "./state_machine_test.go").CombinedOutput()
	require.NoError(t, err)
	fmt.Println(string(out))
}
