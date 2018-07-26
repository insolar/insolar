package main

type HelloWorlder struct {
	Greeted int
}

func (hw *HelloWorlder) Hello() (string, error) {
	hw.Greeted++
	return "Hello world", nil
}
