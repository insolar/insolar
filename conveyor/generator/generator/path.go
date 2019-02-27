package generator

import (
	"path"
	"strings"
)

const projectPath = "github.com/insolar/insolar"
const basePath = "conveyor/generator/state_machines"

func sourceFile(dir string, file string) string {
	return path.Join(basePath, dir, file)
}

func generatedFile(file string) string {
	dir, file := path.Split(file)
	file = file[0:len(file)-3] + "_generated.go"
	return path.Join(dir, file)
}

func importPath(dir string) string {
	return path.Join(projectPath, basePath, dir)
}

func modulePath(dir string) string {
	fullPath := path.Join(projectPath, basePath, dir)
	pathItems := strings.Split(fullPath, "/")
	if pathItems[len(pathItems)-1] == "" && len(pathItems) > 2 {
		return pathItems[len(pathItems)-2]
	}
	return pathItems[len(pathItems)-1]
}