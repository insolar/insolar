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

package goplugintestutils

import (
	"go/build"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/insolar/insolar/core/message"
	"github.com/insolar/insolar/inscontext"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/testutils"
	"github.com/pkg/errors"

	"github.com/insolar/insolar/core"
	"github.com/stretchr/testify/assert"
	"github.com/ugorji/go/codec"
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
	ARef         core.RecordRef
	ACode        []byte
	AMachineType core.MachineType
}

// Ref implementation for tests
func (t *TestCodeDescriptor) Ref() *core.RecordRef {
	return &t.ARef
}

// MachineType implementation for tests
func (t *TestCodeDescriptor) MachineType() core.MachineType {
	return t.AMachineType
}

// Code implementation for tests
func (t *TestCodeDescriptor) Code() ([]byte, error) {
	return t.ACode, nil
}

// TestClassDescriptor ...
type TestClassDescriptor struct {
	AM    *TestArtifactManager
	ARef  *core.RecordRef
	ACode *core.RecordRef
}

// HeadRef ...
func (t *TestClassDescriptor) HeadRef() *core.RecordRef {
	return t.ARef
}

// StateID ...
func (t *TestClassDescriptor) StateID() *core.RecordID {
	panic("not implemented")
}

// CodeDescriptor ...≈
func (t *TestClassDescriptor) CodeDescriptor() core.CodeDescriptor {
	res, ok := t.AM.Codes[*t.ACode]
	if !ok {
		return nil
	}
	return res
}

// TestObjectDescriptor implementation for tests
type TestObjectDescriptor struct {
	AM                *TestArtifactManager
	Data              []byte
	Code              *core.RecordRef
	Class             *core.RecordRef
	Delegates         map[core.RecordRef]core.RecordRef
	ChildrenContainer []core.RecordRef
}

// HeadRef implementation for tests
func (t *TestObjectDescriptor) HeadRef() *core.RecordRef {
	panic("not implemented")
}

// StateID implementation for tests
func (t *TestObjectDescriptor) StateID() *core.RecordID {
	panic("not implemented")
}

// Memory implementation for tests
func (t *TestObjectDescriptor) Memory() []byte {
	return t.Data
}

// Children implementation for tests
func (t *TestObjectDescriptor) Children(pulse *core.PulseNumber) (core.RefIterator, error) {
	panic("not implemented")
}

// ClassDescriptor implementation for tests
func (t *TestObjectDescriptor) ClassDescriptor(state *core.RecordRef) (core.ClassDescriptor, error) {
	if t.Class == nil {
		return nil, errors.New("No class")
	}

	res, ok := t.AM.Classes[*t.Class]
	if !ok {
		return nil, errors.New("No class")
	}
	return res, nil
}

// TestArtifactManager implementation for tests
type TestArtifactManager struct {
	Types   []core.MachineType
	Codes   map[core.RecordRef]*TestCodeDescriptor
	Objects map[core.RecordRef]*TestObjectDescriptor
	Classes map[core.RecordRef]*TestClassDescriptor
}

// GetChildren implementation for tests
func (t *TestArtifactManager) GetChildren(ctx core.Context, parent core.RecordRef, pulse *core.PulseNumber) (core.RefIterator, error) {
	panic("implement me")
}

// NewTestArtifactManager implementation for tests
func NewTestArtifactManager() *TestArtifactManager {
	return &TestArtifactManager{
		Codes:   make(map[core.RecordRef]*TestCodeDescriptor),
		Objects: make(map[core.RecordRef]*TestObjectDescriptor),
		Classes: make(map[core.RecordRef]*TestClassDescriptor),
	}
}

// Start implementation for tests
func (t *TestArtifactManager) Start(components core.Components) error { return nil }

// Stop implementation for tests
func (t *TestArtifactManager) Stop() error { return nil }

// GenesisRef implementation for tests
func (t *TestArtifactManager) GenesisRef() *core.RecordRef { return &core.RecordRef{} }

// RegisterRequest implementation for tests
func (t *TestArtifactManager) RegisterRequest(ctx core.Context, message core.Message) (*core.RecordRef, error) {
	nonce := testutils.RandomRef()
	return &nonce, nil
}

// GetClass implementation for tests
func (t *TestArtifactManager) GetClass(ctx core.Context, object core.RecordRef, state *core.RecordRef) (core.ClassDescriptor, error) {
	res, ok := t.Classes[object]
	if !ok {
		return nil, errors.New("No object")
	}
	return res, nil
}

// GetObject implementation for tests
func (t *TestArtifactManager) GetObject(ctx core.Context, object core.RecordRef, state *core.RecordRef) (core.ObjectDescriptor, error) {
	res, ok := t.Objects[object]
	if !ok {
		return nil, errors.New("No object")
	}
	return res, nil
}

// GetDelegate implementation for tests
func (t *TestArtifactManager) GetDelegate(ctx core.Context, head, asClass core.RecordRef) (*core.RecordRef, error) {
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
func (t *TestArtifactManager) DeclareType(ctx core.Context, domain core.RecordRef, request core.RecordRef, typeDec []byte) (*core.RecordRef, error) {
	panic("not implemented")
}

// DeployCode implementation for tests
func (t *TestArtifactManager) DeployCode(ctx core.Context, domain core.RecordRef, request core.RecordRef, code []byte, mt core.MachineType) (*core.RecordRef, error) {
	ref := testutils.RandomRef()

	t.Codes[ref] = &TestCodeDescriptor{
		ARef:         ref,
		ACode:        code,
		AMachineType: core.MachineTypeGoPlugin,
	}
	return &ref, nil
}

// GetCode implementation for tests
func (t *TestArtifactManager) GetCode(ctx core.Context, code core.RecordRef) (core.CodeDescriptor, error) {
	res, ok := t.Codes[code]
	if !ok {
		return nil, errors.New("No code")
	}
	return res, nil
}

// ActivateClass implementation for tests
func (t *TestArtifactManager) ActivateClass(ctx core.Context, domain core.RecordRef, request core.RecordRef, code core.RecordRef) (*core.RecordID, error) {
	t.Classes[request] = &TestClassDescriptor{
		AM:   t,
		ARef: &request,
	}

	id := testutils.RandomID()
	return &id, nil
}

// DeactivateClass implementation for tests
func (t *TestArtifactManager) DeactivateClass(domain core.RecordRef, request core.RecordRef, class core.RecordRef) (*core.RecordID, error) {
	panic("not implemented")
}

// UpdateClass implementation for tests
func (t *TestArtifactManager) UpdateClass(domain core.RecordRef, request core.RecordRef, class core.RecordRef, code core.RecordRef, migrationRefs []core.RecordRef) (*core.RecordID, error) {
	classDesc, ok := t.Classes[class]
	if !ok {
		return nil, errors.New("wrong class")
	}
	classDesc.ACode = &code
	id := testutils.RandomID()
	return &id, nil
}

// ActivateObject implementation for tests
func (t *TestArtifactManager) ActivateObject(domain core.RecordRef, request core.RecordRef, class core.RecordRef, parent core.RecordRef, memory []byte) (*core.RecordID, error) {
	codeRef := t.Classes[class].ACode

	t.Objects[request] = &TestObjectDescriptor{
		AM:        t,
		Data:      memory,
		Code:      codeRef,
		Class:     &class,
		Delegates: make(map[core.RecordRef]core.RecordRef),
	}

	id := testutils.RandomID()
	return &id, nil
}

// ActivateObjectDelegate implementation for tests
func (t *TestArtifactManager) ActivateObjectDelegate(domain, request, class, parent core.RecordRef, memory []byte) (*core.RecordID, error) {
	id, err := t.ActivateObject(domain, request, class, parent, memory)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to generate ref")
	}

	pObj, ok := t.Objects[parent]
	if !ok {
		return nil, errors.New("No parent to inject delegate into")
	}

	pObj.Delegates[class] = request

	return id, nil
}

// DeactivateObject implementation for tests
func (t *TestArtifactManager) DeactivateObject(domain core.RecordRef, request core.RecordRef, obj core.RecordRef) (*core.RecordID, error) {
	panic("not implemented")
}

// UpdateObject implementation for tests
func (t *TestArtifactManager) UpdateObject(domain core.RecordRef, request core.RecordRef, obj core.RecordRef, memory []byte) (*core.RecordID, error) {
	objDesc, ok := t.Objects[obj]
	if !ok {
		return nil, errors.New("No object to update")
	}

	objDesc.Data = memory

	// TODO: return real exact "ref"
	return &core.RecordID{}, nil
}

// CBORMarshal - testing serialize helper
func CBORMarshal(t testing.TB, o interface{}) []byte {
	ch := new(codec.CborHandle)
	var data []byte
	err := codec.NewEncoderBytes(&data, ch).Encode(o)
	assert.NoError(t, err, "Marshal")
	return data
}

// CBORUnMarshal - testing deserialize helper
func CBORUnMarshal(t testing.TB, data []byte) interface{} {
	ch := new(codec.CborHandle)
	var ret interface{}
	err := codec.NewDecoderBytes(data, ch).Decode(&ret)
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
	am core.ArtifactManager,
	domain core.RecordRef,
	request core.RecordRef,
	mtype core.MachineType,
	code []byte,
) (
	typeRef *core.RecordRef,
	codeRef *core.RecordRef,
	classRef *core.RecordRef,
	err error,
) {
	ctx := inscontext.TODO()
	codeRef, err = am.DeployCode(
		ctx, domain, request, code, mtype,
	)
	assert.NoError(t, err, "create code on ledger")

	nonce := testutils.RandomRef()
	classRef, err = am.RegisterRequest(ctx, &message.CallConstructor{ClassRef: nonce})
	assert.NoError(t, err)
	_, err = am.ActivateClass(ctx, domain, *classRef, *codeRef)
	assert.NoError(t, err, "create template for contract data")

	return typeRef, codeRef, classRef, err
}

// ContractsBuilder for tests
type ContractsBuilder struct {
	root string

	ArtifactManager core.ArtifactManager
	IccPath         string
	Classes         map[string]*core.RecordRef
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
		Classes:         make(map[string]*core.RecordRef),
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
	ctx := inscontext.TODO()

	for name := range contracts {
		nonce := testutils.RandomRef()
		class, err := cb.ArtifactManager.RegisterRequest(ctx, &message.CallConstructor{ClassRef: nonce})
		if err != nil {
			return err
		}

		log.Debugf("Registered class %q for contract %q in %q", class.String(), name, cb.root)
		cb.Classes[name] = class
	}

	re := regexp.MustCompile(`package\s+\S+`)
	for name, code := range contracts {
		code = re.ReplaceAllString(code, "package main")
		err := WriteFile(filepath.Join(cb.root, "src/contract", name), "main.go", code)
		if err != nil {
			return err
		}
		err = cb.proxy(name)
		if err != nil {
			return err
		}
		err = cb.wrapper(name)
		if err != nil {
			return err
		}
	}

	for name := range contracts {
		log.Debugf("Building plugin for contract %q in %q", name, cb.root)
		err := cb.plugin(name)
		if err != nil {
			return err
		}
		log.Debugf("Built plugin for contract %q", name)

		pluginBinary, err := ioutil.ReadFile(filepath.Join(cb.root, "plugins", name+".so"))
		if err != nil {
			return err
		}

		code, err := cb.ArtifactManager.DeployCode(
			ctx,
			core.RecordRef{}, core.RecordRef{},
			pluginBinary, core.MachineTypeGoPlugin,
		)
		if err != nil {
			return err
		}
		log.Debugf("Deployed code %q for contract %q in %q", code.String(), name, cb.root)
		cb.Codes[name] = code

		_, err = cb.ArtifactManager.ActivateClass(
			ctx,
			core.RecordRef{}, *cb.Classes[name],
			*code,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func (cb *ContractsBuilder) proxy(name string) error {
	dstDir := filepath.Join(cb.root, "src/github.com/insolar/insolar/application/proxy", name)

	err := os.MkdirAll(dstDir, 0777)
	if err != nil {
		return err
	}

	contractPath := filepath.Join(cb.root, "src/contract", name, "main.go")

	out, err := exec.Command(
		cb.IccPath, "proxy",
		"-o", filepath.Join(dstDir, "main.go"),
		"--code-reference", cb.Classes[name].String(),
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
		return err
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
