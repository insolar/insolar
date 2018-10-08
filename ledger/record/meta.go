package record

// ChildRecord is a child activation record. Its used for children iterating.
type ChildRecord struct {
	PrevChild *ID

	Ref Reference // Reference to the child's head.
}

// Next returns next record.
func (r *ChildRecord) Next() *ID {
	if r == nil {
		return nil
	}

	return r.PrevChild
}
