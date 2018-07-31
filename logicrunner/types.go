package logicrunner

type Object struct {
	MachineType MachineType
	Reference   string
	Code        []byte
	Data        []byte
}

type API struct {
}

type Arguments []byte
type Type int
type FuncTag int
type FuncSig struct {
	Args []Type
	Ret  []Type
	Tags []FuncTag
}

type Func struct {
	Name string
	Sig  FuncSig
}

type Interface []Func

type Value struct {
	Type  Type
	Value []byte
}
