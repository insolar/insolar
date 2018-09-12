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

package testutil

import (
	"crypto/rand"
	"go/build"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"

	"github.com/pkg/errors"

	"testing"

	"github.com/insolar/insolar/core"
	"github.com/stretchr/testify/assert"
	"github.com/ugorji/go/codec"
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

// TestCodeDescriptor implementation for tests
type TestCodeDescriptor struct {
	ARef         *core.RecordRef
	ACode        []byte
	AMachineType core.MachineType
}

// Ref implementation for tests
func (t *TestCodeDescriptor) Ref() *core.RecordRef {
	return t.ARef
}

// MachineType implementation for tests
func (t *TestCodeDescriptor) MachineType() (core.MachineType, error) {
	return t.AMachineType, nil
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

// StateRef ...
func (t *TestClassDescriptor) StateRef() *core.RecordRef {
	panic("not implemented")
}

// CodeDescriptor ...
func (t *TestClassDescriptor) CodeDescriptor() (core.CodeDescriptor, error) {
	res, ok := t.AM.Codes[*t.ACode]
	if !ok {
		return nil, errors.New("No code")
	}
	return res, nil
}

// TestObjectDescriptor implementation for tests
type TestObjectDescriptor struct {
	AM   *TestArtifactManager
	Data []byte
	Code *core.RecordRef
}

// HeadRef implementation for tests
func (t *TestObjectDescriptor) HeadRef() *core.RecordRef {
	panic("not implemented")
}

// StateRef implementation for tests
func (t *TestObjectDescriptor) StateRef() *core.RecordRef {
	panic("not implemented")
}

// Memory implementation for tests
func (t *TestObjectDescriptor) Memory() ([]byte, error) {
	return t.Data, nil
}

// CodeDescriptor implementation for tests
func (t *TestObjectDescriptor) CodeDescriptor() (core.CodeDescriptor, error) {
	if t.Code == nil {
		return nil, errors.New("No code")
	}

	res, ok := t.AM.Codes[*t.Code]
	if !ok {
		return nil, errors.New("No code")
	}
	return res, nil
}

// ClassDescriptor implementation for tests
func (t *TestObjectDescriptor) ClassDescriptor() (core.ClassDescriptor, error) {
	panic("not implemented")
}

// TestArtifactManager implementation for tests
type TestArtifactManager struct {
	Types   []core.MachineType
	Codes   map[core.RecordRef]*TestCodeDescriptor
	Objects map[core.RecordRef]*TestObjectDescriptor
	Classes map[core.RecordRef]*TestClassDescriptor
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

// RootRef implementation for tests
func (t *TestArtifactManager) RootRef() *core.RecordRef { return &core.RecordRef{} }

// SetArchPref implementation for tests
func (t *TestArtifactManager) SetArchPref(pref []core.MachineType) {
	t.Types = pref
}

// GetExactObj implementation for tests
func (t *TestArtifactManager) GetExactObj(class core.RecordRef, object core.RecordRef) ([]byte, []byte, error) {
	panic("not implemented")
}

// GetLatestClass implementation for tests
func (t *TestArtifactManager) GetLatestClass(object core.RecordRef) (core.ClassDescriptor, error) {
	res, ok := t.Classes[object]
	if !ok {
		return nil, errors.New("No object")
	}
	return res, nil
}

// GetLatestObj implementation for tests
func (t *TestArtifactManager) GetLatestObj(object core.RecordRef) (core.ObjectDescriptor, error) {
	res, ok := t.Objects[object]
	if !ok {
		return nil, errors.New("No object")
	}
	return res, nil
}

// GetObjChildren implementation for tests
func (t *TestArtifactManager) GetObjChildren(head core.RecordRef) (core.RefIterator, error) {
	panic("not implemented")
}

// GetObjDelegate implementation for tests
func (t *TestArtifactManager) GetObjDelegate(head, asClass core.RecordRef) (*core.RecordRef, error) {
	panic("not implemented")
}

// DeclareType implementation for tests
func (t *TestArtifactManager) DeclareType(domain core.RecordRef, request core.RecordRef, typeDec []byte) (*core.RecordRef, error) {
	panic("not implemented")
}

// DeployCode implementation for tests
func (t *TestArtifactManager) DeployCode(domain core.RecordRef, request core.RecordRef, codeMap map[core.MachineType][]byte) (*core.RecordRef, error) {
	ref, err := randomRef()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to generate ref")
	}
	t.Codes[*ref] = &TestCodeDescriptor{
		ARef:         ref,
		ACode:        codeMap[core.MachineTypeGoPlugin],
		AMachineType: core.MachineTypeGoPlugin,
	}
	return ref, nil
}

// GetCode implementation for tests
func (t *TestArtifactManager) GetCode(code core.RecordRef) (core.CodeDescriptor, error) {
	res, ok := t.Codes[code]
	if !ok {
		return nil, errors.New("No code")
	}
	return res, nil
}

// ActivateClass implementation for tests
func (t *TestArtifactManager) ActivateClass(domain core.RecordRef, request core.RecordRef) (*core.RecordRef, error) {
	ref, err := randomRef()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to generate ref")
	}
	t.Classes[*ref] = &TestClassDescriptor{
		AM:   t,
		ARef: ref,
	}
	return ref, nil
}

// DeactivateClass implementation for tests
func (t *TestArtifactManager) DeactivateClass(domain core.RecordRef, request core.RecordRef, class core.RecordRef) (*core.RecordRef, error) {
	panic("not implemented")
}

// UpdateClass implementation for tests
func (t *TestArtifactManager) UpdateClass(domain core.RecordRef, request core.RecordRef, class core.RecordRef, code core.RecordRef, migrationRefs []core.RecordRef) (*core.RecordRef, error) {
	classDesc, ok := t.Classes[class]
	if !ok {
		return nil, errors.New("wrong class")
	}
	classDesc.ACode = &code
	return randomRef()
}

func randomRef() (*core.RecordRef, error) {
	b := make([]byte, 64)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	ref := core.RecordRef{}
	copy(ref[:], b[0:64])
	return &ref, nil
}

// ActivateObj implementation for tests
func (t *TestArtifactManager) ActivateObj(domain core.RecordRef, request core.RecordRef, class core.RecordRef, parent core.RecordRef, memory []byte) (*core.RecordRef, error) {
	codeRef := t.Classes[class].ACode

	ref, err := randomRef()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to generate ref")
	}

	t.Objects[*ref] = &TestObjectDescriptor{
		AM:   t,
		Data: memory,
		Code: codeRef,
	}
	return ref, nil
}

// ActivateObjDelegate implementation for tests
func (t *TestArtifactManager) ActivateObjDelegate(domain, request, class, parent core.RecordRef, memory []byte) (*core.RecordRef, error) {
	return t.ActivateObj(domain, request, class, parent, memory)
}

// DeactivateObj implementation for tests
func (t *TestArtifactManager) DeactivateObj(domain core.RecordRef, request core.RecordRef, obj core.RecordRef) (*core.RecordRef, error) {
	panic("not implemented")
}

// UpdateObj implementation for tests
func (t *TestArtifactManager) UpdateObj(domain core.RecordRef, request core.RecordRef, obj core.RecordRef, memory []byte) (*core.RecordRef, error) {
	_, ok := t.Objects[obj]
	if !ok {
		return nil, errors.New("No object to update")
	}
	// TODO: return real exact "ref"
	return &core.RecordRef{}, nil
}

// CBORMarshal - testing serialize helper
func CBORMarshal(t *testing.T, o interface{}) []byte {
	ch := new(codec.CborHandle)
	var data []byte
	err := codec.NewEncoderBytes(&data, ch).Encode(o)
	assert.NoError(t, err, "Marshal")
	return data
}

// CBORUnMarshal - testing deserialize helper
func CBORUnMarshal(t *testing.T, data []byte) interface{} {
	ch := new(codec.CborHandle)
	var ret interface{}
	err := codec.NewDecoderBytes(data, ch).Decode(&ret)
	assert.NoError(t, err, "serialise")
	return ret
}

// AMPublishCode publishes code on ledger
func AMPublishCode(
	t *testing.T,
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
	codeRef, err = am.DeployCode(
		domain, request, map[core.MachineType][]byte{mtype: code},
	)
	assert.NoError(t, err, "create code on ledger")

	classRef, err = am.ActivateClass(domain, request)
	assert.NoError(t, err, "create template for contract data")
	_, err = am.UpdateClass(domain, request, *classRef, *codeRef, nil)
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
	return &ContractsBuilder{ArtifactManager: am, IccPath: icc}
}

// Build ...
func (cb *ContractsBuilder) Build(contracts map[string]string) error {
	tmpDir, err := ioutil.TempDir("", "test-")
	if err != nil {
		return err
	}
	cb.root = tmpDir
	defer os.RemoveAll(cb.root) // nolint: errcheck

	cb.Classes = make(map[string]*core.RecordRef)
	for name := range contracts {
		class, err := cb.ArtifactManager.ActivateClass(
			core.RecordRef{}, core.RecordRef{},
		)
		if err != nil {
			return err
		}

		cb.Classes[name] = class
	}

	re := regexp.MustCompile("package\\s+\\S+")
	for name, code := range contracts {
		code = re.ReplaceAllString(code, "package main")
		err := WriteFile(cb.root+"/src/contract/"+name+"/", "main.go", code)
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

	cb.Codes = make(map[string]*core.RecordRef)
	for name := range contracts {
		err := cb.plugin(name)
		if err != nil {
			return err
		}

		pluginBinary, err := ioutil.ReadFile(cb.root + "/plugins/" + name + ".so")
		if err != nil {
			return err
		}

		code, err := cb.ArtifactManager.DeployCode(
			core.RecordRef{}, core.RecordRef{},
			map[core.MachineType][]byte{core.MachineTypeGoPlugin: pluginBinary},
		)
		if err != nil {
			return err
		}
		cb.Codes[name] = code

		_, err = cb.ArtifactManager.UpdateClass(
			core.RecordRef{}, core.RecordRef{},
			*cb.Classes[name],
			*code,
			[]core.RecordRef{},
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func (cb *ContractsBuilder) proxy(name string) error {
	dstDir := cb.root + "/src/contract-proxy/" + name

	err := os.MkdirAll(dstDir, 0777)
	if err != nil {
		return err
	}

	contractPath := cb.root + "/src/contract/" + name + "/main.go"

	out, err := exec.Command(
		cb.IccPath, "proxy",
		"-o", dstDir+"/main.go",
		"--code-reference", cb.Classes[name].String(),
		contractPath,
	).CombinedOutput()
	if err != nil {
		return errors.Wrap(err, "can't generate proxy: "+string(out))
	}
	return nil
}

func (cb *ContractsBuilder) wrapper(name string) error {
	contractPath := cb.root + "/src/contract/" + name + "/main.go"
	wrapperPath := cb.root + "/src/contract/" + name + "/main_wrapper.go"

	out, err := exec.Command(cb.IccPath, "wrapper", "-o", wrapperPath, contractPath).CombinedOutput()
	if err != nil {
		return errors.Wrap(err, "can't generate wrapper for contract '"+name+"': "+string(out))
	}
	return nil
}

// Plugin ...
func (cb *ContractsBuilder) plugin(name string) error {
	dstDir := cb.root + "/plugins/"

	err := os.MkdirAll(dstDir, 0777)
	if err != nil {
		return err
	}

	origGoPath, err := ChangeGoPath(cb.root)
	if err != nil {
		return err
	}
	defer os.Setenv("GOPATH", origGoPath) // nolint: errcheck

	//contractPath := root + "/src/contract/" + name + "/main.go"

	out, err := exec.Command("go", "build", "-buildmode=plugin", "-o", dstDir+"/"+name+".so", "contract/"+name).CombinedOutput()
	if err != nil {
		return errors.Wrap(err, "can't build contract: "+string(out))
	}
	return nil
}
