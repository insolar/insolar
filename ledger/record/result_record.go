package record

type ReasonCode uint

type ResultRecord struct {
	AppDataRecord

	RequestRecord RecordReference
}

type WipeOutRecord struct {
	ResultRecord

	Replacement RecordReference
	WipedHash   RecordHash
}

type StatelessResult struct {
	ResultRecord
}

type ReadRecordResult struct {
	StatelessResult

	RecordBody []byte
}

type StatelessCallResult struct {
	StatelessResult

	resultMemory Memory
}

type StatelessExceptionResult struct {
	StatelessCallResult

	ExceptionType RecordReference
}

type ReadObjectResult struct {
	StatelessResult

	State            int
	MomoryProjection Memory
}

type SpecialResult struct {
	ResultRecord

	ReasonCode ReasonCode
}

type LockUnlockResult struct {
	SpecialResult
}

type RejectionResult struct {
	SpecialResult
}

type StatefulResult struct {
	ResultRecord
}

type ActivationRecord struct {
	StatefulResult

	GoverningDomain RecordReference
}

type ClassActivateRecord struct {
	ActivationRecord

	CodeRecord    RecordReference
	DefaultMemory Memory
}

type ObjectActivateRecord struct {
	ActivationRecord

	ClassActivateRecord RecordReference
	Memory              Memory
}

type StorageRecord struct {
	StatefulResult
}

type CodeRecord struct {
	StorageRecord

	Interfaces   []RecordReference
	TargetedCode [][]byte // []MachineBinaryCode
	SourceCode   string   // ObjectSourceCode
}

type AmendRecord struct {
	StatefulResult

	BaseRecord    RecordReference
	AmendedRecord RecordReference
}

type ClassAmendRecord struct {
	AmendRecord

	NewCode []byte // ObjectBinaryCode
}

func (r *ClassAmendRecord) MigrationCodes() []*MemoryMigrationCode {
	panic("not implemented")
}

type MemoryMigrationCode struct {
	ClassAmendRecord

	GeneratedByClassRecord RecordReference
	MigrationCodeRecord    RecordReference
}

type DeactivationRecord struct {
	AmendRecord
}

type ObjectAmendRecord struct {
	AmendRecord

	NewMemory Memory
}

type StatefulCallResult struct {
	ObjectAmendRecord

	ResultMemory Memory
}

type StatefulExceptionResult struct {
	StatefulCallResult

	ExceptionType RecordReference
}

type EnforcedObjectAmendRecord struct {
	ObjectAmendRecord
}

type ObjectAppendRecord struct {
	AmendRecord

	AppendMemory Memory
}
