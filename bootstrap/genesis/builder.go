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
	"go/build"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/internal/ledger/artifact"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/logicrunner/preprocessor"
	"github.com/pkg/errors"
)

var (
	contractSources = insolar.RootModule + "/application/contract"
	proxySources    = insolar.RootModule + "/application/proxy"
)

// prototypes holds name -> code reference pair
type prototypes map[string]*insolar.Reference

// contractsBuilder for tests
type contractsBuilder struct {
	root            string
	prototypes      prototypes
	artifactManager artifact.Manager
}

// newContractBuilder returns a new `contractsBuilder`,
// requires initialized artifact manager.
func newContractBuilder(am artifact.Manager) *contractsBuilder {
	tmpDir, err := ioutil.TempDir("", "test-")
	if err != nil {
		return nil
	}

	cb := &contractsBuilder{
		root:            tmpDir,
		prototypes:      make(map[string]*insolar.Reference),
		artifactManager: am,
	}
	return cb
}

// clean deletes tmp directory used for contracts building
func (cb *contractsBuilder) clean() {
	log.Debugf("Cleaning build directory %q", cb.root)
	err := os.RemoveAll(cb.root)
	if err != nil {
		panic(err)
	}
}

func (cb *contractsBuilder) buildPrototypes(ctx context.Context, rootDomainID *insolar.ID) (prototypes, error) {
	inslog := inslogger.FromContext(ctx)
	inslog.Info("[ buildSmartContracts ] building contracts:", contractNames)
	contracts, err := parseContracts()
	if err != nil {
		return nil, errors.Wrap(err, "[ buildSmartContracts ] failed to get contracts map")
	}

	inslog.Info("[ buildSmartContracts ] Start building contracts ...")
	err = cb.build(ctx, contracts, rootDomainID)
	if err != nil {
		return nil, errors.Wrap(err, "[ buildSmartContracts ] couldn't build contracts")
	}
	inslog.Info("[ buildSmartContracts ] Stop building contracts ...")

	return cb.prototypes, nil
}

// buildPrototypes ...
func (cb *contractsBuilder) build(ctx context.Context, contracts map[string]*preprocessor.ParsedFile, domain *insolar.ID) error {

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
			return errors.Wrap(err, "[ buildPrototypes ] Can't RegisterRequest for contract")
		}
		cb.prototypes[name] = insolar.NewReference(*domain, *protoID)
	}

	for name, code := range contracts {
		code.ChangePackageToMain()

		ctr, err := createFileInDir(filepath.Join(cb.root, "src/contract", name), "main.go")
		if err != nil {
			return errors.Wrap(err, "[ buildPrototypes ] Can't open contract file")
		}
		err = code.Write(ctr)
		ctr.Close()
		if err != nil {
			return errors.Wrap(err, "[ buildPrototypes ] Can't makeFileWithDir")
		}

		proxyPath := filepath.Join(cb.root, "src", proxySources, name)
		proxy, err := createFileInDir(proxyPath, "main.go")
		if err != nil {
			return errors.Wrap(err, "[ buildPrototypes ] Can't open proxy file")
		}
		err = code.WriteProxy(cb.prototypes[name].String(), proxy)
		proxy.Close()
		if err != nil {
			return errors.Wrap(err, "[ buildPrototypes ] Can't write proxy")
		}

		wrp, err := createFileInDir(filepath.Join(cb.root, "src/contract", name), "main_wrapper.go")
		if err != nil {
			return errors.Wrap(err, "[ buildPrototypes ] Can't open wrapper file")
		}
		err = code.WriteWrapper(wrp, "main")
		wrp.Close()
		if err != nil {
			return errors.Wrap(err, "[ buildPrototypes ] Can't write wrapper")
		}
	}

	for name := range contracts {
		log.Debugf("Building plugin for contract %q in %q", name, cb.root)
		err := cb.plugin(name)
		if err != nil {
			return errors.Wrap(err, "[ buildPrototypes ] Can't call plugin")
		}

		pluginBinary, err := ioutil.ReadFile(filepath.Join(cb.root, "plugins", name+".so"))
		if err != nil {
			return errors.Wrap(err, "[ buildPrototypes ] Can't ReadFile")
		}
		codeReq, err := cb.artifactManager.RegisterRequest(
			ctx,
			*domainRef,
			&message.Parcel{
				Msg: &message.GenesisRequest{Name: name + "_code"},
			},
		)
		if err != nil {
			return errors.Wrapf(err, "[ buildPrototypes ] Can't RegisterRequest for code '%v'", name)
		}

		log.Debugf("Deploying code for contract %q", name)
		codeID, err := cb.artifactManager.DeployCode(
			ctx,
			*domainRef, *insolar.NewReference(*domain, *codeReq),
			pluginBinary, insolar.MachineTypeGoPlugin,
		)
		if err != nil {
			return errors.Wrapf(err, "[ buildPrototypes ] Can't DeployCode for code '%v", name)
		}

		codeRef := insolar.NewReference(*domain, *codeID)
		_, err = cb.artifactManager.RegisterResult(ctx, *domainRef, *codeRef, nil)
		if err != nil {
			return errors.Wrapf(err, "[ buildPrototypes ] Can't SetRecord for code '%v'", name)
		}

		log.Debugf("Deployed code %q for contract %q in %q", codeRef.String(), name, cb.root)

		_, err = cb.artifactManager.ActivatePrototype(
			ctx,
			*domainRef,
			*cb.prototypes[name],
			insolar.GenesisRecord.Ref(),
			*codeRef,
			nil,
		)
		if err != nil {
			return errors.Wrapf(err, "[ buildPrototypes ] Can't ActivatePrototypef for code '%v'", name)
		}

		_, err = cb.artifactManager.RegisterResult(ctx, *domainRef, *cb.prototypes[name], nil)
		if err != nil {
			return errors.Wrapf(err, "[ buildPrototypes ] Can't RegisterResult of prototype for code '%v'", name)
		}
	}

	return nil
}

// compile plugin
func (cb *contractsBuilder) plugin(name string) error {
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
	cmd.Env = append(os.Environ(), "GOPATH="+prependGoPath(cb.root))
	out, err := cmd.CombinedOutput()
	if err != nil {
		return errors.Wrapf(err, "can't build contract: %v", string(out))
	}
	return nil
}

func goPATH() string {
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		gopath = build.Default.GOPATH
	}
	return gopath
}

func getContractPath(name string) (string, error) {
	contractDir := filepath.Join(goPATH(), "src", contractSources)
	contractFile := name + ".go"
	return filepath.Join(contractDir, name, contractFile), nil
}

func parseContracts() (map[string]*preprocessor.ParsedFile, error) {
	contracts := make(map[string]*preprocessor.ParsedFile)
	for _, name := range contractNames {
		contractPath, err := getContractPath(name)
		if err != nil {
			return nil, errors.Wrap(err, "[ contractsMap ] couldn't get path to contracts: ")
		}
		parsed, err := preprocessor.ParseFile(contractPath, insolar.MachineTypeGoPlugin)
		if err != nil {
			return nil, errors.Wrapf(err, "[ contractsMap ] couldn't read contract: %v", contractPath)
		}
		contracts[name] = parsed
	}
	return contracts, nil
}

// prependGoPath prepends `path` to GOPATH environment variable
// accounting for possibly for default value. Returns new value.
// NOTE: that environment is not changed
func prependGoPath(path string) string {
	return path + string(os.PathListSeparator) + goPATH()
}

// makeFileWithDir dumps data into file in provided directory, creates directory if it does not exist.
func makeFileWithDir(dir string, name string, data []byte) error {
	err := os.MkdirAll(dir, 0775)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filepath.Join(dir, name), data, 0644)
}

// createFileInDir opens file in provided directory, creates directory if it does not exist.
func createFileInDir(dir string, name string) (*os.File, error) {
	err := os.MkdirAll(dir, 0775)
	if err != nil {
		return nil, err
	}
	return os.OpenFile(filepath.Join(dir, name), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
}
