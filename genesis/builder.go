/*
 *    Copyright 2018 Insolar
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package genesis

import (
	"context"
	"go/build"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/message"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/testutils"
	"github.com/pkg/errors"
)

// PrependGoPath prepends `path` to GOPATH environment variable
// accounting for possibly for default value. Returns new value.
// NOTE: that environment is not changed
func PrependGoPath(path string) string {
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		gopath = build.Default.GOPATH
	}

	return path + string(os.PathListSeparator) + gopath
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

// ContractsBuilder for tests
type ContractsBuilder struct {
	root string

	ArtifactManager core.ArtifactManager
	IccPath         string
	Prototypes      map[string]*core.RecordRef
	Codes           map[string]*core.RecordRef
}

// NewContractBuilder returns a new `ContractsBuilder`, takes in: path to tmp directory,
// artifact manager, ...
func NewContractBuilder(am core.ArtifactManager, icc string) *ContractsBuilder {
	tmpDir, err := ioutil.TempDir("", "test-")
	if err != nil {
		return nil
	}

	cb := &ContractsBuilder{
		root:            tmpDir,
		Prototypes:      make(map[string]*core.RecordRef),
		Codes:           make(map[string]*core.RecordRef),
		ArtifactManager: am,
		IccPath:         icc}
	return cb
}

// Clean deletes tmp directory used for contracts building
func (cb *ContractsBuilder) Clean() {
	log.Debugf("Cleaning build directory %q", cb.root)
	err := os.RemoveAll(cb.root) // nolint: errcheck
	if err != nil {
		panic(err)
	}
}

// Build ...
func (cb *ContractsBuilder) Build(contracts map[string]string) error {
	ctx := context.TODO()

	for name := range contracts {
		nonce := testutils.RandomRef()
		protoID, err := cb.ArtifactManager.RegisterRequest(
			ctx, &message.Parcel{Msg: &message.CallConstructor{PrototypeRef: nonce}},
		)
		if err != nil {
			return errors.Wrap(err, "[ Build ] Can't RegisterRequest")
		}

		protoRef := core.RecordRef{}
		protoRef.SetRecord(*protoID)
		log.Debugf("Registered prototype %q for contract %q in %q", protoRef.String(), name, cb.root)
		cb.Prototypes[name] = &protoRef
	}

	re := regexp.MustCompile(`package\s+\S+`)
	for name, code := range contracts {
		//preprocessor.ParseFile()
		code = re.ReplaceAllString(code, "package main")
		err := WriteFile(filepath.Join(cb.root, "src/contract", name), "main.go", code)
		if err != nil {
			return errors.Wrap(err, "[ Build ] Can't WriteFile")
		}
		err = cb.proxy(name)
		if err != nil {
			return errors.Wrap(err, "[ Build ] Can't call proxy")
		}
		err = cb.wrapper(name)
		if err != nil {
			return errors.Wrap(err, "[ Build ] Can't call wrapper")
		}
	}

	for name := range contracts {
		log.Debugf("Building plugin for contract %q in %q", name, cb.root)
		err := cb.plugin(name)
		if err != nil {
			return errors.Wrap(err, "[ Build ] Can't call plugin")
		}
		log.Debugf("Built plugin for contract %q", name)

		pluginBinary, err := ioutil.ReadFile(filepath.Join(cb.root, "plugins", name+".so"))
		if err != nil {
			return errors.Wrap(err, "[ Build ] Can't ReadFile")
		}

		log.Debugf("Deploying code for contract %q", name)
		codeID, err := cb.ArtifactManager.DeployCode(
			ctx,
			core.RecordRef{}, core.RecordRef{},
			pluginBinary, core.MachineTypeGoPlugin,
		)
		codeRef := &core.RecordRef{}
		codeRef.SetRecord(*codeID)
		if err != nil {
			return errors.Wrap(err, "[ Build ] Can't SetRecord")
		}
		log.Debugf("Deployed code %q for contract %q in %q", codeRef.String(), name, cb.root)
		cb.Codes[name] = codeRef

		// FIXME: It's a temporary fix and should not be here. Ii will NOT work properly on production. Remove it ASAP!
		_, err = cb.ArtifactManager.ActivatePrototype(
			ctx,
			core.RecordRef{},
			*cb.Prototypes[name],
			*cb.ArtifactManager.GenesisRef(), // FIXME: Only bootstrap can do this!
			*codeRef,
			nil,
		)
		if err != nil {
			return errors.Wrap(err, "[ Build ] Can't ActivatePrototype")
		}
	}

	return nil
}

func (cb *ContractsBuilder) proxy(name string) error {
	dstDir := filepath.Join(cb.root, "src/github.com/insolar/insolar/application/proxy", name)

	err := os.MkdirAll(dstDir, 0777)
	if err != nil {
		return errors.Wrap(err, "[ proxy ]")
	}

	contractPath := filepath.Join(cb.root, "src/contract", name, "main.go")

	out, err := exec.Command(
		cb.IccPath, "proxy",
		"-o", filepath.Join(dstDir, "main.go"),
		"--code-reference", cb.Prototypes[name].String(),
		contractPath,
	).CombinedOutput()
	if err != nil {
		return errors.Wrap(err, "can't generate proxy: "+string(out))
	}
	return nil
}

func (cb *ContractsBuilder) wrapper(name string) error {
	contractPath := filepath.Join(cb.root, "src/contract", name, "main.go")
	wrapperPath := filepath.Join(cb.root, "src/contract", name, "main_wrapper.go")

	out, err := exec.Command(cb.IccPath, "wrapper", "-o", wrapperPath, contractPath).CombinedOutput()
	if err != nil {
		return errors.Wrap(err, "can't generate wrapper for contract '"+name+"': "+string(out))
	}
	return nil
}

// Plugin ...
func (cb *ContractsBuilder) plugin(name string) error {
	dstDir := filepath.Join(cb.root, "plugins")

	err := os.MkdirAll(dstDir, 0777)
	if err != nil {
		return errors.Wrap(err, "[ plugin ]")
	}

	cmd := exec.Command(
		"go", "build",
		"-buildmode=plugin",
		"-o", filepath.Join(dstDir, name+".so"),
		filepath.Join(cb.root, "src/contract", name),
	)
	cmd.Env = append(os.Environ(), "GOPATH="+PrependGoPath(cb.root))
	out, err := cmd.CombinedOutput()
	if err != nil {
		return errors.Wrap(err, "can't build contract: "+string(out))
	}
	return nil
}
