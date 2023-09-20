package warning

type Warning struct {
	error
}

func New(err error) Warning {
	return Warning{err}
}
