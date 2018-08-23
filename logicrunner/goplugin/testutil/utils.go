package testutil

import (
	"go/build"
	"os"
)

// ChangeGoPath prepends `path` to GOPATH environment variable
// accounting for possibly for default value. Returns original
// value of the enviroment variable, don't forget to restore
// it with defer:
//    defer os.Setenv("GOPATH", origGoPath)
func ChangeGoPath(path string) (string, error) {
	gopathOrigEnv := os.Getenv("GOPATH")
	gopath := gopathOrigEnv
	if gopath == "" {
		gopath = build.Default.GOPATH
	}

	err := os.Setenv("GOPATH", path+":"+gopath)
	if err != nil {
		return "", err
	}
	return gopathOrigEnv, nil
}

// WriteFile dumps `text` into file named `name` into directory `dir`.
// Creates directory if needed as well as file
func WriteFile(dir string, name string, text string) error {
	err := os.MkdirAll(dir, 0777)
	if err != nil {
		return err
	}

	fh, err := os.OpenFile(dir+"/"+name, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}

	_, err = fh.WriteString(text)
	if err != nil {
		return err
	}

	err = fh.Close()
	if err != nil {
		return err
	}

	return nil
}
