//
// Copyright 2019 Insolar Technologies GmbH
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
//

package genesis

import (
	"context"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/internal/ledger/artifact"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/logicrunner/goplugin/preprocessor"
)

var (
	contractSources = insolar.RootModule + "/application/contract"
	proxySources    = insolar.RootModule + "/application/proxy"
)

// PrependGoPath prepends `path` to GOPATH environment variable
// accounting for possibly for default value. Returns new value.
// NOTE: that environment is not changed
func PrependGoPath(path string) string {
	return path + string(os.PathListSeparator) + goPATH()
}

// WriteFile dumps `text` into file named `name` into directory `dir`.
// Creates directory if needed as well as file
func WriteFile(dir string, name string, text string) error {
	err := os.MkdirAll(dir, 0775)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filepath.Join(dir, name), []byte(text), 0644)
}

func OpenFile(dir string, name string) (*os.File, error) {
	err := os.MkdirAll(dir, 0775)
	if err != nil {
		return nil, err
	}
	return os.OpenFile(filepath.Join(dir, name), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
}

// ContractsBuilder for tests
type ContractsBuilder struct {
	root string

	Prototypes map[string]*insolar.Reference
	Codes      map[string]*insolar.Reference

	genesisRef      insolar.Reference
	artifactManager artifact.Manager
}

// NewContractBuilder returns a new `ContractsBuilder`,
// requires initialized artifact manager.
func NewContractBuilder(genesisRef insolar.Reference, am artifact.Manager) *ContractsBuilder {
	tmpDir, err := ioutil.TempDir("", "test-")
	if err != nil {
		return nil
	}

	cb := &ContractsBuilder{
		root:       tmpDir,
		Prototypes: make(map[string]*insolar.Reference),
		Codes:      make(map[string]*insolar.Reference),

		genesisRef:      genesisRef,
		artifactManager: am,
	}
	return cb
}

// Clean deletes tmp directory used for contracts building
func (cb *ContractsBuilder) Clean() {
	log.Debugf("Cleaning build directory %q", cb.root)
	err := os.RemoveAll(cb.root)
	if err != nil {
		panic(err)
	}
}

// Build ...
func (cb *ContractsBuilder) Build(ctx context.Context, contracts map[string]*preprocessor.ParsedFile, domain *insolar.ID) error {

	domainRef := insolar.NewReference(*domain, *domain)

	for name := range contracts {
		protoID, err := cb.artifactManager.RegisterRequest(
			ctx,
			*domainRef,
			&message.Parcel{
				Msg: &message.GenesisRequest{
					Name: name + "_proto",
				},
			})
		if err != nil {
			return errors.Wrap(err, "[ Build ] Can't RegisterRequest for contract")
		}
		cb.Prototypes[name] = insolar.NewReference(*domain, *protoID)
	}

	for name, code := range contracts {
		code.ChangePackageToMain()

		ctr, err := OpenFile(filepath.Join(cb.root, "src/contract", name), "main.go")
		if err != nil {
			return errors.Wrap(err, "[ Build ] Can't open contract file")
		}
		err = code.Write(ctr)
		ctr.Close()
		if err != nil {
			return errors.Wrap(err, "[ Build ] Can't WriteFile")
		}

		proxy, err := OpenFile(filepath.Join(cb.root, "src/github.com/insolar/insolar/application/proxy", name), "main.go")
		if err != nil {
			return errors.Wrap(err, "[ Build ] Can't open proxy file")
		}
		err = code.WriteProxy(cb.Prototypes[name].String(), proxy)
		proxy.Close()
		if err != nil {
			return errors.Wrap(err, "[ Build ] Can't write proxy")
		}

		wrp, err := OpenFile(filepath.Join(cb.root, "src/contract", name), "main_wrapper.go")
		if err != nil {
			return errors.Wrap(err, "[ Build ] Can't open wrapper file")
		}
		err = code.WriteWrapper(wrp)
		wrp.Close()
		if err != nil {
			return errors.Wrap(err, "[ Build ] Can't write wrapper")
		}
	}

	for name := range contracts {
		log.Debugf("Building plugin for contract %q in %q", name, cb.root)
		err := cb.plugin(name)
		if err != nil {
			return errors.Wrap(err, "[ Build ] Can't call plugin")
		}

		pluginBinary, err := ioutil.ReadFile(filepath.Join(cb.root, "plugins", name+".so"))
		if err != nil {
			return errors.Wrap(err, "[ Build ] Can't ReadFile")
		}
		codeReq, err := cb.artifactManager.RegisterRequest(
			ctx,
			*domainRef,
			&message.Parcel{
				Msg: &message.GenesisRequest{Name: name + "_code"},
			},
		)
		if err != nil {
			return errors.Wrapf(err, "[ Build ] Can't RegisterRequest for code '%v'", name)
		}

		log.Debugf("Deploying code for contract %q", name)
		codeID, err := cb.artifactManager.DeployCode(
			ctx,
			*domainRef, *insolar.NewReference(*domain, *codeReq),
			pluginBinary, insolar.MachineTypeGoPlugin,
		)
		if err != nil {
			return errors.Wrapf(err, "[ Build ] Can't DeployCode for code '%v", name)
		}

		codeRef := insolar.NewReference(*domain, *codeID)
		_, err = cb.artifactManager.RegisterResult(ctx, *domainRef, *codeRef, nil)
		if err != nil {
			return errors.Wrapf(err, "[ Build ] Can't SetRecord for code '%v'", name)
		}

		log.Debugf("Deployed code %q for contract %q in %q", codeRef.String(), name, cb.root)
		cb.Codes[name] = codeRef

		_, err = cb.artifactManager.ActivatePrototype(
			ctx,
			*domainRef,
			*cb.Prototypes[name],
			cb.genesisRef,
			*codeRef,
			nil,
		)
		if err != nil {
			return errors.Wrapf(err, "[ Build ] Can't ActivatePrototypef for code '%v'", name)
		}

		_, err = cb.artifactManager.RegisterResult(ctx, *domainRef, *cb.Prototypes[name], nil)
		if err != nil {
			return errors.Wrapf(err, "[ Build ] Can't RegisterResult of prototype for code '%v'", name)
		}
	}

	return nil
}

// compile plugin
func (cb *ContractsBuilder) plugin(name string) error {
	dstDir := filepath.Join(cb.root, "plugins")

	err := os.MkdirAll(dstDir, 0777)
	if err != nil {
		return errors.Wrap(err, "[ plugin ]")
	}

	cmd := exec.Command(
		"go",
		"build",
		"-buildmode=plugin",
		"-o", filepath.Join(dstDir, name+".so"),
		filepath.Join(cb.root, "src/contract", name),
	)
	cmd.Env = append(os.Environ(), "GOPATH="+PrependGoPath(cb.root))
	out, err := cmd.CombinedOutput()
	if err != nil {
		return errors.Wrapf(err, "can't build contract: %v", string(out))
	}
	return nil
}
