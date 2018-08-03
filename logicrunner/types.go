package logicrunner

// Reference is a contract address
type Reference string

// Object is an inner representation of storage object for transfwering it over API
type Object struct {
	MachineType MachineType
	Reference   Reference
	Code        []byte
	Data        []byte
}

// Arguments is a dedicated type for arguments, that represented as bynary cbored blob
type Arguments []byte
