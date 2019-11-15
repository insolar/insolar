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

package goplugintestutils

import (
	"context"
	"go/build"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/application"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/api"
	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/logicrunner/artifacts"
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
	root    string
	IccPath string

	pulseAccessor   pulse.Accessor
	artifactManager artifacts.Client
	jetCoordinator  jet.Coordinator

	Prototypes map[string]*insolar.Reference
	Codes      map[string]*insolar.Reference
}

// NewContractBuilder returns a new `ContractsBuilder`, takes in: path to tmp directory,
// artifact manager, ...
func NewContractBuilder(icc string, am artifacts.Client, pa pulse.Accessor, jc jet.Coordinator) *ContractsBuilder {
	tmpDir := insolar.ContractBuildTmpDir("test-")

	cb := &ContractsBuilder{
		root:    tmpDir,
		IccPath: icc,

		pulseAccessor:   pa,
		artifactManager: am,
		jetCoordinator:  jc,

		Prototypes: make(map[string]*insolar.Reference),
		Codes:      make(map[string]*insolar.Reference),
	}
	return cb
}

func (cb *ContractsBuilder) Clean() {
	log.Debugf("Cleaning build directory %q", cb.root)
	err := os.RemoveAll(cb.root)
	if err != nil {
		panic(err)
	}
}

// Build ...
func (cb *ContractsBuilder) Build(ctx context.Context, contracts map[string]string) error {
	logger := inslogger.FromContext(ctx)

	for name := range contracts {
		pulse, err := cb.pulseAccessor.Latest(ctx)
		if err != nil {
			return errors.Wrap(err, "can't get current pulse")
		}

		nonce := gen.Reference()
		request := record.IncomingRequest{
			CallType:  record.CTDeployPrototype,
			Prototype: &nonce,
			Reason:    api.MakeReason(pulse.PulseNumber, []byte(name)),
			APINode:   cb.jetCoordinator.Me(),
		}
		protoID, err := cb.registerRequest(ctx, &request)

		if err != nil {
			return errors.Wrap(err, "[ Build ] Can't RegisterIncomingRequest")
		}
		protoRef := insolar.NewReference(*protoID)
		logger.Debugf("Registered prototype %q for contract %q in %q", protoRef.String(), name, cb.root)
		cb.Prototypes[name] = protoRef
	}

	re := regexp.MustCompile(`package\s+\S+`)
	for name, code := range contracts {
		code = re.ReplaceAllString(code, "package main")
		err := WriteFile(filepath.Join(cb.root, "src/contract", name), "main.go", code)
		if err != nil {
			return errors.Wrap(err, "[ buildPrototypes ] Can't WriteFile")
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
		logger.Debug("Building plugin for contract ", name, " in ", cb.root)

		err := cb.plugin(name)
		if err != nil {
			return errors.Wrap(err, "[ Build ] Can't call plugin")
		}
		logger.Debug("Built plugin for contract ", name)

		pluginBinary, err := ioutil.ReadFile(filepath.Join(cb.root, "plugins", name+".so"))
		if err != nil {
			return errors.Wrap(err, "[ Build ] Can't ReadFile")
		}

		logger.Debug("Deploying code for contract ", name)
		codeID, err := cb.artifactManager.DeployCode(ctx, pluginBinary, insolar.MachineTypeGoPlugin)
		if err != nil {
			return errors.Wrap(err, "[ Build ] DeployCode returns error")
		}

		codeRef := insolar.NewReference(*codeID)

		logger.Debugf("Deployed code %q for contract %q in %q", codeRef.String(), name, cb.root)
		cb.Codes[name] = codeRef

		// FIXME: It's a temporary fix and should not be here. Ii will NOT work properly on production. Remove it ASAP!
		err = cb.artifactManager.ActivatePrototype(
			ctx,
			*cb.Prototypes[name],
			application.GenesisRecord.Ref(), // FIXME: Only bootstrap can do this!
			*codeRef,
			nil,
		)
		if err != nil {
			return errors.Wrap(err, "[ Build ] Can't ActivatePrototype")
		}
	}

	return nil
}

// Using registerRequest without VM is a tmp solution while there is no logic of contract uploading in VM
// Because of this we need copy some logic in test code
func (cb *ContractsBuilder) registerRequest(ctx context.Context, request *record.IncomingRequest) (*insolar.ID, error) {
	var err error
	var lastPulse insolar.PulseNumber

	retries := 5
	logger := inslogger.FromContext(ctx)

	if cb.pulseAccessor == nil {
		return nil, errors.New("No pulse accessor")
	}

	for current := 1; current <= retries; current++ {
		currentPulse, err := cb.pulseAccessor.Latest(ctx)
		if err != nil {
			return nil, errors.Wrap(err, "[ registerRequest ] Can't get latest pulse")
		}

		if currentPulse.PulseNumber == lastPulse {
			logger.Debugf("[ registerRequest ]  wait for pulse change. Current: %d", currentPulse.PulseNumber)
			time.Sleep(100 * time.Millisecond)
			continue
		}
		lastPulse = currentPulse.PulseNumber

		reqInfo, err := cb.artifactManager.RegisterIncomingRequest(ctx, request)
		if err == nil || !strings.Contains(err.Error(), flow.ErrCancelled.Error()) {
			reqID := reqInfo.RequestID
			return &reqID, err
		}

		logger.Debugf("[ registerRequest ] retry. attempt: %d/%d", current, retries)
	}
	return nil, errors.Wrap(err, "flow cancelled, retries exceeded")
}

func (cb *ContractsBuilder) proxy(name string) error {
	root := insolar.RootModuleDir()
	dstDir := filepath.Join(root, "application/proxy", name)

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
		// "-trimpath", // if enable this option, the plugin will not load
		"-mod=vendor",
		"-o", filepath.Join(dstDir, name+".so"),
		filepath.Join(cb.root, "src/contract", name),
	)
	cmd.Env = append(os.Environ(), "GOPATH="+PrependGoPath(cb.root))
	cmd.Env = append(cmd.Env, "GO111MODULE=on")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return errors.Wrap(err, "can't build contract: "+string(out))
	}
	return nil
}
