package record

// ChildRecord is a child activation record. Its used for children iterating.
type ChildRecord struct {
	ChainRecord

	Child Reference // Reference to the child's head.
}
