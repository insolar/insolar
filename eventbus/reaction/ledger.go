package reaction

import (
	"io"

	"github.com/insolar/insolar/core"
)

type Code struct {
	Code        []byte
	MachineType core.MachineType
}

func (e *Code) Serialize() (io.Reader, error) {
	return serialize(e, TypeCode)
}

type Class struct {
	Head  core.RecordRef
	State core.RecordRef
	Code  *core.RecordRef // Can be nil.
}

func (e *Class) Serialize() (io.Reader, error) {
	return serialize(e, TypeClass)
}

type Object struct {
	Head     core.RecordRef
	State    core.RecordRef
	Class    core.RecordRef
	Memory   []byte
	Children []core.RecordRef
}

func (e *Object) Serialize() (io.Reader, error) {
	return serialize(e, TypeObject)
}

type Delegate struct {
	Head core.RecordRef
}

func (e *Delegate) Serialize() (io.Reader, error) {
	return serialize(e, TypeDelegate)
}

type Reference struct {
	Ref core.RecordRef
}

func (e *Reference) Serialize() (io.Reader, error) {
	return serialize(e, TypeReference)
}
