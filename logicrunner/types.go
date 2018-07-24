package logicrunner

type Object struct {
}

type API struct {
}

func (API) Call(address string, method string, args [][]byte) {
	panic("implement me")
}
