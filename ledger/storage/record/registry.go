package record

var registry = map[TypeID]Record{}

// Register makes provided record serializable. Should be called for each record in init().
func register(id TypeID, r Record) {
	if _, ok := registry[id]; ok {
		panic("duplicate record type")
	}

	registry[id] = r
}

// Registry returns records by type.
func Registry() map[TypeID]Record {
	res := map[TypeID]Record{}
	for id, rec := range registry {
		res[id] = rec
	}
	return res
}
