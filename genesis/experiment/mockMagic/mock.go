package mockMagic

type Reference struct {
	domain string
	record string
}

type MockMagic struct {}

func (mf *MockMagic) MockGetCaller() *Reference {
	return &Reference{domain: "2", record: "1"}
}

func (mf *MockMagic) MockGetMyOwner() *Reference {
	return &Reference{domain: "1", record: "2"}
}

func (mf *MockMagic) MockSelfDestructRequest() {}

func (mf *MockMagic) MockGetSelfReference() *Reference {
	return &Reference{domain: "1", record: "1"}
}
