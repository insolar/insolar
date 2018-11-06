package message
func Extract(msg core.Message) core.RecordRef {
	switch m := msg.(type) {
	default:
		panic("unknow message type")
	}
}
func ExtractRole(msg core.Message) core.JetRole {
	switch _ := msg.(type) {
	default:
		panic("unknow message type")
	}
}
