package main

import (
	"encoding/json"
	"os"
	"path"
	"text/template"

	"github.com/davecgh/go-spew/spew"

	"github.com/insolar/insolar/logicrunner/preprocessor"
)


const DomainName = "domain"

func generateDomainBase(list preprocessor.ContractList, basePath string) error {
	spew.Dump(list)

	dPath := path.Join(basePath, DomainName)

	_, err := os.Stat(dPath)
	if err == nil {
		os.RemoveAll(dPath)
	} else {
		if _, ok := err.(*os.PathError); !ok {
		checkError(err)
		}
	}


	j, err := json.Marshal(list)
	checkError(err)

	err = os.Mkdir(dPath, 0755)
	checkError(err)

	f, err := os.OpenFile(path.Join(dPath, DomainName+".go"), os.O_CREATE + os.O_WRONLY, 0644)
	checkError(err)
	defer f.Close()

	t := template.Must(template.New("Domain").Parse(domainTemplate))
	checkError(t.Execute(f, string(j)))

	return nil
}


const  domainTemplate =
`package domain

import (
//	"github.com/insolar/insolar/application/appfoundation"
//	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/logicrunner/builtin/foundation"
)


type Domain struct {
	foundation.BaseContract
}


func (d *Domain)GetSchema() (string, error) {
	schema := ` + "`{{.}}`" + `
	return schema, nil
}

`
