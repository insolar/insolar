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
	"testing"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/logicrunner/artifacts"
	"github.com/insolar/insolar/testutils"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
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

// TestCodeDescriptor implementation for tests
type TestCodeDescriptor struct {
	ARef         insolar.Reference
	ACode        []byte
	AMachineType insolar.MachineType
}

// Ref implementation for tests
func (t *TestCodeDescriptor) Ref() *insolar.Reference {
	return &t.ARef
}

// MachineType implementation for tests
func (t *TestCodeDescriptor) MachineType() insolar.MachineType {
	return t.AMachineType
}

// Code implementation for tests
func (t *TestCodeDescriptor) Code() ([]byte, error) {
	return t.ACode, nil
}

// TestObjectDescriptor implementation for tests
type TestObjectDescriptor struct {
	AM                *TestArtifactManager
	ARef              *insolar.Reference
	Data              []byte
	State             *insolar.ID
	PrototypeRef      *insolar.Reference
	Delegates         map[insolar.Reference]insolar.Reference
	ChildrenContainer []insolar.Reference
}

func (t *TestObjectDescriptor) HasPendingRequests() bool {
	panic("implement me")
}

// Parent implementation for tests
func (t *TestObjectDescriptor) Parent() *insolar.Reference {
	panic("implement me")
}

// ChildPointer implementation for tests
func (t *TestObjectDescriptor) ChildPointer() *insolar.ID {
	panic("not implemented")
}

// HeadRef implementation for tests
func (t *TestObjectDescriptor) HeadRef() *insolar.Reference {
	return t.ARef
}

// StateID implementation for tests
func (t *TestObjectDescriptor) StateID() *insolar.ID {
	return t.State
}

// Memory implementation for tests
func (t *TestObjectDescriptor) Memory() []byte {
	return t.Data
}

// IsPrototype implementation for tests
func (t *TestObjectDescriptor) IsPrototype() bool {
	return false
}

// Prototype implementation for tests
func (t *TestObjectDescriptor) Prototype() (*insolar.Reference, error) {
	if t.PrototypeRef == nil {
		panic("No prototype")
	}
	return t.PrototypeRef, nil
}

// Code implementation for tests
func (t *TestObjectDescriptor) Code() (*insolar.Reference, error) {
	if t.PrototypeRef == nil {
		panic("No code")
	}
	return t.PrototypeRef, nil
}

// TestArtifactManager implementation for tests
type TestArtifactManager struct {
	Types      []insolar.MachineType
	Codes      map[insolar.Reference]*TestCodeDescriptor
	Objects    map[insolar.Reference]*TestObjectDescriptor
	Prototypes map[insolar.Reference]*TestObjectDescriptor
}

func (t *TestArtifactManager) GetPendingRequest(ctx context.Context, objectID insolar.ID) (insolar.Parcel, error) {
	panic("implement me")
}

func (t *TestArtifactManager) HasPendingRequests(ctx context.Context, object insolar.Reference) (bool, error) {
	panic("implement me")
}

// State implementation for tests
func (t *TestArtifactManager) State() ([]byte, error) {
	panic("implement me")
}

// GetChildren implementation for tests
func (t *TestArtifactManager) GetChildren(ctx context.Context, parent insolar.Reference, pulse *insolar.PulseNumber) (artifacts.RefIterator, error) {
	panic("implement me")
}

// NewTestArtifactManager implementation for tests
func NewTestArtifactManager() *TestArtifactManager {
	return &TestArtifactManager{
		Codes:      make(map[insolar.Reference]*TestCodeDescriptor),
		Objects:    make(map[insolar.Reference]*TestObjectDescriptor),
		Prototypes: make(map[insolar.Reference]*TestObjectDescriptor),
	}
}

// RegisterRequest implementation for tests
func (t *TestArtifactManager) RegisterRequest(ctx context.Context, obj insolar.Reference, parcel insolar.Parcel) (*insolar.ID, error) {
	nonce := testutils.RandomID()
	return &nonce, nil
}

// RegisterResult saves VM method call result.
func (t *TestArtifactManager) RegisterResult(
	ctx context.Context, object, request insolar.Reference, payload []byte,
) (*insolar.ID, error) {
	panic("implement me")
}

// GetObject implementation for tests
func (t *TestArtifactManager) GetObject(ctx context.Context, object insolar.Reference, state *insolar.ID, approved bool) (artifacts.ObjectDescriptor, error) {
	res, ok := t.Objects[object]
	if !ok {
		return nil, errors.New("No object")
	}
	return res, nil
}

// GetDelegate implementation for tests
func (t *TestArtifactManager) GetDelegate(ctx context.Context, head, asClass insolar.Reference) (*insolar.Reference, error) {
	obj, ok := t.Objects[head]
	if !ok {
		return nil, errors.New("No object")
	}

	res, ok := obj.Delegates[asClass]
	if !ok {
		return nil, errors.New("No delegate")
	}

	return &res, nil
}

// DeclareType implementation for tests
func (t *TestArtifactManager) DeclareType(ctx context.Context, domain insolar.Reference, request insolar.Reference, typeDec []byte) (*insolar.ID, error) {
	panic("not implemented")
}

// DeployCode implementation for tests
func (t *TestArtifactManager) DeployCode(ctx context.Context, domain insolar.Reference, request insolar.Reference, code []byte, mt insolar.MachineType) (*insolar.ID, error) {
	ref := testutils.RandomRef()

	t.Codes[ref] = &TestCodeDescriptor{
		ARef:         ref,
		ACode:        code,
		AMachineType: insolar.MachineTypeGoPlugin,
	}
	id := ref.Record()
	return id, nil
}

// GetCode implementation for tests
func (t *TestArtifactManager) GetCode(ctx context.Context, code insolar.Reference) (artifacts.CodeDescriptor, error) {
	res, ok := t.Codes[code]
	if !ok {
		return nil, errors.New("No code")
	}
	return res, nil
}

// ActivatePrototype implementation for tests
func (t *TestArtifactManager) ActivatePrototype(
	ctx context.Context,
	domain, request, parent, code insolar.Reference,
	memory []byte,
) (artifacts.ObjectDescriptor, error) {
	id := testutils.RandomID()

	t.Prototypes[request] = &TestObjectDescriptor{
		AM:           t,
		ARef:         &request,
		Data:         memory,
		State:        &id,
		PrototypeRef: &code,
		Delegates:    make(map[insolar.Reference]insolar.Reference),
	}

	return t.Objects[request], nil
}

// ActivateObject implementation for tests
func (t *TestArtifactManager) ActivateObject(
	ctx context.Context,
	domain, request, parent, prototype insolar.Reference,
	asDelegate bool,
	memory []byte,
) (artifacts.ObjectDescriptor, error) {
	id := testutils.RandomID()

	t.Objects[request] = &TestObjectDescriptor{
		AM:           t,
		ARef:         &request,
		Data:         memory,
		State:        &id,
		PrototypeRef: &prototype,
		Delegates:    make(map[insolar.Reference]insolar.Reference),
	}
	if asDelegate {
		pObj, ok := t.Objects[parent]
		if !ok {
			return nil, errors.New("No parent to inject delegate into")
		}

		pObj.Delegates[prototype] = request
	}

	return t.Objects[request], nil
}

// DeactivateObject implementation for tests
func (t *TestArtifactManager) DeactivateObject(
	ctx context.Context,
	domain insolar.Reference, request insolar.Reference, obj artifacts.ObjectDescriptor,
) (*insolar.ID, error) {
	panic("not implemented")
}

// UpdatePrototype implementation for tests
func (t *TestArtifactManager) UpdatePrototype(
	ctx context.Context,
	domain insolar.Reference,
	request insolar.Reference,
	object artifacts.ObjectDescriptor,
	memory []byte,
	code *insolar.Reference,
) (artifacts.ObjectDescriptor, error) {
	objDesc, ok := t.Prototypes[*object.HeadRef()]
	if !ok {
		return nil, errors.New("No object to update")
	}

	objDesc.Data = memory

	// TODO: return real exact "ref"
	return objDesc, nil
}

// UpdateObject implementation for tests
func (t *TestArtifactManager) UpdateObject(
	ctx context.Context,
	domain insolar.Reference,
	request insolar.Reference,
	object artifacts.ObjectDescriptor,
	memory []byte,
) (artifacts.ObjectDescriptor, error) {
	objDesc, ok := t.Objects[*object.HeadRef()]
	if !ok {
		return nil, errors.New("No object to update")
	}

	objDesc.Data = memory

	// TODO: return real exact "ref"
	return objDesc, nil
}

// RegisterValidation implementation for tests
func (t *TestArtifactManager) RegisterValidation(
	ctx context.Context,
	object insolar.Reference,
	state insolar.ID,
	isValid bool,
	validationMessages []insolar.Message,
) error {
	panic("implement me")
}

// CBORMarshal - testing serialize helper
func CBORMarshal(t testing.TB, o interface{}) []byte {
	data, err := insolar.Serialize(o)
	assert.NoError(t, err, "Marshal")
	return data
}

// CBORUnMarshal - testing deserialize helper
func CBORUnMarshal(t testing.TB, data []byte) interface{} {
	var ret interface{}
	err := insolar.Deserialize(data, &ret)
	assert.NoError(t, err, "serialise")
	return ret
}

// CBORUnMarshalToSlice - wrapper for CBORUnMarshal, expects slice
func CBORUnMarshalToSlice(t testing.TB, in []byte) []interface{} {
	r := CBORUnMarshal(t, in)
	assert.IsType(t, []interface{}{}, r)
	return r.([]interface{})
}

// AMPublishCode publishes code on ledger
func AMPublishCode(
	t testing.TB,
	am artifacts.Client,
	domain insolar.Reference,
	request insolar.Reference,
	mtype insolar.MachineType,
	code []byte,
) (
	typeRef *insolar.Reference,
	codeRef *insolar.Reference,
	protoRef *insolar.Reference,
	err error,
) {
	ctx := context.TODO()
	codeID, err := am.DeployCode(
		ctx, domain, request, code, mtype,
	)
	assert.NoError(t, err, "create code on ledger")
	codeRef = &insolar.Reference{}
	codeRef.SetRecord(*codeID)

	nonce := testutils.RandomRef()
	protoID, err := am.RegisterRequest(ctx, insolar.GenesisRecord.Ref(), &message.Parcel{Msg: &message.CallConstructor{PrototypeRef: nonce}})
	assert.NoError(t, err)
	protoRef = &insolar.Reference{}
	protoRef.SetRecord(*protoID)
	_, err = am.ActivatePrototype(ctx, domain, *protoRef, insolar.GenesisRecord.Ref(), *codeRef, nil)
	assert.NoError(t, err, "create template for contract data")

	return typeRef, codeRef, protoRef, err
}

// ContractsBuilder for tests
type ContractsBuilder struct {
	root string

	ArtifactManager artifacts.Client
	IccPath         string
	Prototypes      map[string]*insolar.Reference
	Codes           map[string]*insolar.Reference
}

// NewContractBuilder returns a new `ContractsBuilder`, takes in: path to tmp directory,
// artifact manager, ...
func NewContractBuilder(am artifacts.Client, icc string) *ContractsBuilder {
	tmpDir, err := ioutil.TempDir("", "test-")
	if err != nil {
		return nil
	}

	cb := &ContractsBuilder{
		root:            tmpDir,
		Prototypes:      make(map[string]*insolar.Reference),
		Codes:           make(map[string]*insolar.Reference),
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
			ctx, insolar.GenesisRecord.Ref(), &message.Parcel{Msg: &message.CallConstructor{PrototypeRef: nonce}},
		)
		if err != nil {
			return errors.Wrap(err, "[ Build ] Can't RegisterRequest")
		}

		protoRef := insolar.Reference{}
		protoRef.SetRecord(*protoID)
		log.Debugf("Registered prototype %q for contract %q in %q", protoRef.String(), name, cb.root)
		cb.Prototypes[name] = &protoRef
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
		nonce := testutils.RandomRef()
		codeReq, err := cb.ArtifactManager.RegisterRequest(
			ctx, insolar.GenesisRecord.Ref(), &message.Parcel{Msg: &message.CallConstructor{PrototypeRef: nonce}},
		)
		if err != nil {
			return errors.Wrap(err, "[ Build ] Can't RegisterRequest")
		}

		log.Debugf("Deploying code for contract %q", name)
		codeID, err := cb.ArtifactManager.DeployCode(
			ctx,
			insolar.Reference{}, *insolar.NewReference(insolar.ID{}, *codeReq),
			pluginBinary, insolar.MachineTypeGoPlugin,
		)
		codeRef := &insolar.Reference{}
		codeRef.SetRecord(*codeID)
		if err != nil {
			return errors.Wrap(err, "[ Build ] Can't SetRecord")
		}
		log.Debugf("Deployed code %q for contract %q in %q", codeRef.String(), name, cb.root)
		cb.Codes[name] = codeRef

		// FIXME: It's a temporary fix and should not be here. Ii will NOT work properly on production. Remove it ASAP!
		_, err = cb.ArtifactManager.ActivatePrototype(
			ctx,
			insolar.Reference{},
			*cb.Prototypes[name],
			insolar.GenesisRecord.Ref(), // FIXME: Only bootstrap can do this!
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
