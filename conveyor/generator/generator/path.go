package generator

import (
	"path"
	"strings"
)

func (g *Generator) sourceFile(dir string, file string) string {
	return path.Join(g.path, dir, file)
}

func (g *Generator) generatedFile(file string) string {
	dir, file := path.Split(file)
	file = file[0:len(file)-3] + "_generated.go"
	return path.Join(dir, file)
}

func (g *Generator) importPath(dir string) string {
	return path.Join(g.base, g.path, dir)
}

func (g *Generator) modulePath(dir string) string {
	fullPath := path.Join(g.base, g.path, dir)
	pathItems := strings.Split(fullPath, "/")
	if pathItems[len(pathItems)-1] == "" && len(pathItems) > 2 {
		return pathItems[len(pathItems)-2]
	}
	return pathItems[len(pathItems)-1]
}
